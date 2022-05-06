package impl

import (
	"context"
	"fmt"
	"logistics/fetcher"
	"logistics/model"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type HLFetcherConfig struct {
	Username string  `yaml:"username"`
	Password string  `yaml:"password"`
	Domain   string  `yaml:"domain"`
	QueryUrl *string `yaml:"query_url"`
}

func (h *HLFetcherConfig) Parse(value *yaml.Node) error {
	return value.Decode(h)
}

type HLFetcherFactory struct{}

func (H HLFetcherFactory) ConstructFetcher(config fetcher.FetcherConfig) (fetcher.Fetcher, error) {
	fetcherConfig, ok := config.(*HLFetcherConfig)
	if !ok {
		return nil, errors.New("config not right")
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return NewHLFetcher(&http.Client{
		Jar: jar,
	}, *fetcherConfig), nil
}

func (H HLFetcherFactory) ConstructConfig() fetcher.FetcherConfig {
	return &HLFetcherConfig{}
}

type HLFetcher struct {
	config   HLFetcherConfig
	url      string
	signUrl  string
	queryUrl string
	client   *http.Client
}

func NewHLFetcher(client *http.Client, config HLFetcherConfig) *HLFetcher {
	queryUrl := config.Domain + "/priceSearchQuery.htm"
	if config.QueryUrl != nil {
		queryUrl = *config.QueryUrl
	}
	return &HLFetcher{
		config:   config,
		url:      config.Domain + "/index.htm",
		signUrl:  config.Domain + "/signin.htm",
		queryUrl: queryUrl,
		client:   client,
	}
}

func (z HLFetcher) Fetch(ctx context.Context, countryCode string, weight float64) ([]model.Logistics, error) {
	err := z.login(ctx)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", z.queryUrl, strings.NewReader(url.Values{
		"country":   []string{countryCode},
		"cargoType": []string{"P"},
		"weight":    []string{fmt.Sprintf("%v", weight)},
	}.Encode()))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := z.client.Do(req)
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
	reader.Find("table tbody tr").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if i%2 == 0 {
			var td []string
			selection.Find("td").Each(func(i int, selection *goquery.Selection) {
				td = append(td, strings.TrimSpace(selection.Text()))
			})
			var queryWeight, total, price, fare, fuel, other float64
			queryWeight, err = strconv.ParseFloat(td[2], 64)
			if err != nil {
				err = errors.WithStack(err)
				return false
			}
			total, err = strconv.ParseFloat(td[3], 64)
			if err != nil {
				err = errors.WithStack(err)
				return false
			}
			price, err = strconv.ParseFloat(td[4], 64)
			if err != nil {
				err = errors.WithStack(err)
				return false
			}
			fare, err = strconv.ParseFloat(td[5], 64)
			if err != nil {
				err = errors.WithStack(err)
				return false
			}
			fuel, err = strconv.ParseFloat(td[6], 64)
			if err != nil {
				err = errors.WithStack(err)
				return false
			}
			other, err = strconv.ParseFloat(td[7], 64)
			if err != nil {
				err = errors.WithStack(err)
				return false
			}
			res = append(res, model.Logistics{
				URL:    z.url,
				Method: td[0],
				Weight: queryWeight,
				Total:  total,
				Price:  price,
				Fare:   fare,
				Fuel:   fuel,
				Other:  other,
			})
		} else {
			res[len(res)-1].Remark = strings.TrimSpace(selection.Find("td").Text())
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (z HLFetcher) login(ctx context.Context) error {
	resp, err := z.client.PostForm(z.signUrl, url.Values{
		"username": []string{z.config.Username},
		"password": []string{z.config.Password},
	})
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return nil
}
