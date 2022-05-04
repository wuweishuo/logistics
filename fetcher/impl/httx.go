package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"logistics/fetcher"
	"logistics/model"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type HTTXFetcherConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (h *HTTXFetcherConfig) Parse(value *yaml.Node) error {
	return value.Decode(h)
}

type HTTXFetcherFactory struct{}

func (H HTTXFetcherFactory) ConstructFetcher(config fetcher.FetcherConfig) (fetcher.Fetcher, error) {
	fetcherConfig, ok := config.(*HTTXFetcherConfig)
	if !ok {
		return nil, errors.New("config not right")
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return NewHTTXFetcher(*fetcherConfig, &http.Client{
		Jar: jar,
	}), nil
}

func (H HTTXFetcherFactory) ConstructConfig() fetcher.FetcherConfig {
	return &HTTXFetcherConfig{}
}

type HTTXFetcher struct {
	config HTTXFetcherConfig
	client *http.Client
}

func NewHTTXFetcher(config HTTXFetcherConfig, client *http.Client) *HTTXFetcher {
	return &HTTXFetcher{
		config: config,
		client: client,
	}
}

func (h HTTXFetcher) Fetch(ctx context.Context, countryCode string, weight float64) ([]model.Logistics, error) {
	err := h.login(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := h.client.PostForm("https://vip.httx56.com/index.php?controller=finance&action=ajaxprice", url.Values{
		"re_country": []string{countryCode},
		"hwlx":       []string{"02"},
		"weight":     []string{fmt.Sprintf("%v", weight)},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	type QueryResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    []struct {
			Service    string  `json:"service"`
			Weight     float64 `json:"weight"`
			Totalprice string  `json:"totalprice"`
			Price      string  `json:"price"`
			Freight    string  `json:"freight"`
			Tax        string  `json:"tax"`
			Otherfee   string  `json:"otherfee"`
			Bz         string  `json:"bz"`
		} `json:"data"`
	}
	var queryResp QueryResp
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&queryResp)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !queryResp.Success {
		return nil, errors.New(queryResp.Message)
	}
	var res []model.Logistics
	for _, d := range queryResp.Data {
		total, err := strconv.ParseFloat(d.Totalprice, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		price, err := strconv.ParseFloat(d.Price, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		freight, err := strconv.ParseFloat(d.Freight, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		tax, err := strconv.ParseFloat(d.Tax, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		other, err := strconv.ParseFloat(d.Otherfee, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		res = append(res, model.Logistics{
			Source: "汇通天下",
			URL:    "https://vip.httx56.com/index.php/default/login",
			Method: d.Service,
			Weight: d.Weight,
			Total:  total,
			Price:  price,
			Fare:   freight,
			Fuel:   tax,
			Other:  other,
			Remark: d.Bz,
		})
	}
	return res, nil
}

func (h HTTXFetcher) login(ctx context.Context) error {
	resp, err := h.client.PostForm("https://vip.httx56.com/index.php/default/login", url.Values{
		"username": []string{h.config.Username},
		"password": []string{h.config.Password},
	})
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return nil
}
