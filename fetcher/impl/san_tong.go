package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"logistics/fetcher"
	"logistics/model"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type SanTongFetcherConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (s *SanTongFetcherConfig) Parse(value *yaml.Node) error {
	return value.Decode(s)
}

type SanTongFetcherFactory struct{}

func (s SanTongFetcherFactory) ConstructFetcher(config fetcher.FetcherConfig) (fetcher.Fetcher, error) {
	fetcherConfig, ok := config.(*SanTongFetcherConfig)
	if !ok {
		return nil, errors.New("config not right")
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return NewSanTongFetcher(*fetcherConfig, &http.Client{
		Jar: jar,
	}), nil
}

func (s SanTongFetcherFactory) ConstructConfig() fetcher.FetcherConfig {
	return &SanTongFetcherConfig{}
}

// SanTongFetcher [三通订单系统](http://119.23.34.110:8088/)
type SanTongFetcher struct {
	config SanTongFetcherConfig
	client *http.Client
}

func NewSanTongFetcher(config SanTongFetcherConfig, client *http.Client) *SanTongFetcher {
	return &SanTongFetcher{
		config: config,
		client: client,
	}
}

func (s SanTongFetcher) Fetch(ctx context.Context, countryCode string, weight float64) ([]model.Logistics, error) {
	err := s.login(ctx)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "http://119.23.34.110:8088/order/fee-trail/list/page/1/pageSize/20", strings.NewReader(url.Values{
		"country_code": []string{countryCode},
		"product_type": []string{"W"},
		"weight":       []string{fmt.Sprintf("%v", weight)},
	}.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	type SanTongResp struct {
		CountryCode string
		Message     string
		Total       int
		Data        []struct {
			ChargeWeight  float64
			ServiceCnName string
			TotalFee      float64
			FreightFee    float64
			FuelFee       float64
			OtherFee      float64
			Remark        string
		}
	}
	var data SanTongResp
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	res := make([]model.Logistics, 0, data.Total)
	for _, d := range data.Data {
		res = append(res, model.Logistics{
			URL:    "http://119.23.34.110:8088",
			Method: d.ServiceCnName,
			Weight: d.ChargeWeight,
			Total:  d.TotalFee,
			Fare:   d.FreightFee,
			Fuel:   d.FuelFee,
			Other:  d.OtherFee,
			Remark: d.Remark,
		})
	}
	return res, nil
}

func (s SanTongFetcher) login(ctx context.Context) error {
	resp, err := s.client.PostForm("http://119.23.34.110:8088/default/index/login", url.Values{
		"userName": []string{s.config.Username},
		"userPass": []string{s.config.Password},
	})
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return nil
}
