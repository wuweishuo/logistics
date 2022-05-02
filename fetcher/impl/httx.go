package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"logistics/config"
	"logistics/model"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
)

type HTTXFetcher struct {
	client *http.Client
}

func NewHTTXFetcher() *HTTXFetcher {
	jar, _ := cookiejar.New(nil)
	return &HTTXFetcher{
		client: &http.Client{
			Jar: jar,
		},
	}
}

func (h HTTXFetcher) Fetch(ctx context.Context, config config.LoginConfig, countryCode string, weight float64) ([]model.Logistics, error) {
	err := h.login(ctx, config)
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
		total, err := strconv.ParseFloat(d.Totalprice, 10)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		price, err := strconv.ParseFloat(d.Price, 10)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		freight, err := strconv.ParseFloat(d.Freight, 10)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		tax, err := strconv.ParseFloat(d.Tax, 10)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		other, err := strconv.ParseFloat(d.Otherfee, 10)
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

func (h HTTXFetcher) login(ctx context.Context, loginConfig config.LoginConfig) error {
	resp, err := h.client.PostForm("https://vip.httx56.com/index.php/default/login", url.Values{
		"username": []string{loginConfig.Username},
		"password": []string{loginConfig.Password},
	})
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return nil
}
