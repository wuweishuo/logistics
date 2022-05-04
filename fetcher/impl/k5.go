package impl

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"logistics/enums"
	"logistics/fetcher"
	"logistics/model"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type K5FetcherConfig struct {
	Domain     string `yaml:"domain"`
	TotalIdx   int    `yaml:"total_idx"`
	WeightIdx  int    `yaml:"weight_idx"`
	PriceIdx   *int   `yaml:"price_idx"`
	ChannelIdx int    `yaml:"channel_idx"`
}

func (k *K5FetcherConfig) Parse(value *yaml.Node) error {
	return value.Decode(k)
}

type K5FetcherFactory struct{}

func (k K5FetcherFactory) ConstructFetcher(config fetcher.FetcherConfig) (fetcher.Fetcher, error) {
	fetcherConfig, ok := config.(*K5FetcherConfig)
	if !ok {
		return nil, errors.New("config not right")
	}
	return NewK5Fetcher(*fetcherConfig), nil
}

func (k K5FetcherFactory) ConstructConfig() fetcher.FetcherConfig {
	return &K5FetcherConfig{}
}

type K5Fetcher struct {
	config   K5FetcherConfig
	url      string
	queryUrl string
	client   *http.Client
}

func NewK5Fetcher(config K5FetcherConfig) K5Fetcher {
	return K5Fetcher{
		config:   config,
		url:      config.Domain + "/WebPrice?action=list",
		queryUrl: config.Domain + "/WebPrice?action=list",
		client:   http.DefaultClient,
	}
}

func (b K5Fetcher) Fetch(ctx context.Context, countryCode string, weight float64) ([]model.Logistics, error) {
	resp, err := b.client.PostForm(b.queryUrl, url.Values{
		"country":   []string{enums.CountryCodes[countryCode]},
		"rweight":   []string{fmt.Sprintf("%v", weight)},
		"goodstype": []string{"WPX"},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	reader, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var res []model.Logistics
	reader.Find(".ks_dl_table1").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		var td []string
		selection.Find("dd").Each(func(i int, selection *goquery.Selection) {
			td = append(td, strings.TrimSpace(selection.Text()))
		})
		if td[0] == "无记录" {
			return false
		}
		var queryWeight, total, price float64
		queryWeight, err = strconv.ParseFloat(td[b.config.WeightIdx], 64)
		if err != nil {
			err = errors.WithStack(err)
			return false
		}
		total, err = strconv.ParseFloat(td[b.config.TotalIdx], 64)
		if err != nil {
			err = errors.WithStack(err)
			return false
		}
		if b.config.PriceIdx != nil {
			price, err = strconv.ParseFloat(td[*b.config.PriceIdx], 64)
			if err != nil {
				err = errors.WithStack(err)
				return false
			}
		}

		res = append(res, model.Logistics{
			URL:    b.url,
			Method: td[b.config.ChannelIdx],
			Weight: queryWeight,
			Total:  total,
			Price:  price,
			Remark: strings.TrimSpace(selection.Next().Next().Text()),
		})
		return true
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
