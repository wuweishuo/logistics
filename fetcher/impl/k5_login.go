package impl

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"logistics/fetcher"
	"logistics/model"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type K5LoginFetcherConfig struct {
	Domain     string `yaml:"domain"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	ChannelIdx int    `yaml:"channel_idx"`
	WeightIdx  int    `yaml:"weight_idx"`
	TotalIdx   int    `yaml:"total_idx"`
	PriceIdx   int    `yaml:"price_idx"`
	FuelIdx    int    `yaml:"fuel_idx"`
}

func (k *K5LoginFetcherConfig) Parse(value *yaml.Node) error {
	return value.Decode(k)
}

type K5LoginFetcherFactory struct{}

func (k K5LoginFetcherFactory) ConstructFetcher(config fetcher.FetcherConfig) (fetcher.Fetcher, error) {
	fetcherConfig, ok := config.(*K5LoginFetcherConfig)
	if !ok {
		return nil, errors.New("config not right")
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return NewK5LoginFetcher(*fetcherConfig, &http.Client{
		Jar: jar,
	}), nil
}

func (k K5LoginFetcherFactory) ConstructConfig() fetcher.FetcherConfig {
	return &K5LoginFetcherConfig{}
}

type K5LoginFetcher struct {
	config   K5LoginFetcherConfig
	url      string
	loginUrl string
	queryUrl string
	client   *http.Client
}

func NewK5LoginFetcher(config K5LoginFetcherConfig, client *http.Client) K5LoginFetcher {
	return K5LoginFetcher{
		config:   config,
		url:      config.Domain + "/new_index.jsp?retry_reason=PASSWORD",
		loginUrl: config.Domain + "/nclient/Logon?action=logon",
		queryUrl: config.Domain + "/nclient/CClientPrice?action=getAnalyse",
		client:   client,
	}
}

func (c K5LoginFetcher) Fetch(ctx context.Context, countryCode string, weight float64) ([]model.Logistics, error) {
	err := c.login(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.PostForm(c.queryUrl, url.Values{
		"orgLookup.country": []string{countryCode},
		"weight":            []string{fmt.Sprintf("%v", weight)},
		"goodstype":         []string{"WPX"},
		"showDataType":      []string{"1"},
		"logisticflag":      []string{"1"},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.Request.Method == http.MethodGet {
		return nil, errors.New("登录失败")
	}
	reader, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var res []model.Logistics
	reader.Find(".tablelist tbody tr").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if !selection.Parent().Parent().HasClass("tablelist") {
			return true
		}
		if _, exist := selection.Attr("style"); exist {
			return true
		}
		var td []string
		selection.Find("td").Each(func(i int, selection *goquery.Selection) {
			td = append(td, strings.TrimSpace(selection.Text()))
		})
		if td[0] == "无记录" {
			return false
		}
		var total, queryWeight, price, fuel float64
		total, err = strconv.ParseFloat(td[c.config.TotalIdx], 64)
		if err != nil {
			err = errors.WithStack(err)
			return false
		}
		queryWeight, err = strconv.ParseFloat(td[c.config.WeightIdx], 64)
		if err != nil {
			err = errors.WithStack(err)
			return false
		}
		price, err = strconv.ParseFloat(td[c.config.PriceIdx], 64)
		if err != nil {
			err = errors.WithStack(err)
			return false
		}
		if td[c.config.FuelIdx] != "含燃油" {
			fuel, err = strconv.ParseFloat(td[c.config.FuelIdx], 64)
			if err != nil {
				err = errors.WithStack(err)
				return false
			}
		}
		res = append(res, model.Logistics{
			URL:    c.url,
			Method: td[c.config.ChannelIdx],
			Weight: queryWeight,
			Total:  total,
			Price:  price,
			Fuel:   fuel,
			Remark: strings.TrimSpace(selection.Next().Next().Text()),
		})
		return true
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c K5LoginFetcher) login(ctx context.Context) error {
	resp, err := c.client.PostForm(c.loginUrl, url.Values{
		"userid":   []string{c.config.Username},
		"password": []string{c.config.Password},
	})
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return nil
}
