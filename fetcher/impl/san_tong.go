package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"logistics/config"
	"logistics/model"
	"net/http"
	"net/url"
	"strings"
)

// SanTongFetcher [三通订单系统](http://119.23.34.110:8088/)
type SanTongFetcher struct {
	client *http.Client
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

func NewSanTongFetcher() *SanTongFetcher {
	return &SanTongFetcher{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) > 10 {
					return errors.New("stopped after 10 redirects")
				}
				for _, cookie := range req.Response.Cookies() {
					req.AddCookie(cookie)
				}
				return nil
			},
		},
	}
}

func (s SanTongFetcher) Fetch(ctx context.Context, config config.LoginConfig, countryCode string, weight float64) ([]model.Logistics, error) {
	cookies, err := s.getCookies(ctx, config)
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
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
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
	var data SanTongResp
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	res := make([]model.Logistics, 0, data.Total)
	for _, d := range data.Data {
		res = append(res, model.Logistics{
			Source: "三通订单系统",
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

func (s SanTongFetcher) getCookies(ctx context.Context, config config.LoginConfig) ([]*http.Cookie, error) {
	resp, err := s.client.PostForm("http://119.23.34.110:8088/default/index/login", url.Values{
		"userName": []string{config.Username},
		"userPass": []string{config.Password},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return resp.Request.Cookies(), nil
}
