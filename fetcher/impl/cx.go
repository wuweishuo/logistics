package impl

import (
	"context"
	"fmt"
	"logistics/config"
	"logistics/model"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type CXFetcher struct {
	source    string
	url       string
	loginUrl  string
	queryUrl  string
	methodIdx int
	weightIdx int
	totalIdx  int
	priceIdx  int
	fuelIdx   int
	client    *http.Client
}

func NewCXFetcher() *CXFetcher {
	jar, _ := cookiejar.New(nil)
	return &CXFetcher{
		source:    "CX",
		url:       "http://cx.kingtrans.cn/old_index.jsp?retry_reason=PASSWORD",
		loginUrl:  "http://cx.kingtrans.cn/client/Logon?action=logon",
		queryUrl:  "http://cx.kingtrans.cn/nclient/CClientPrice?action=getAnalyse",
		methodIdx: 1,
		weightIdx: 3,
		totalIdx:  5,
		priceIdx:  6,
		fuelIdx:   7,
		client: &http.Client{
			Jar: jar,
		},
	}
}

func (c CXFetcher) Fetch(ctx context.Context, config config.LoginConfig, countryCode string, weight float64) ([]model.Logistics, error) {
	err := c.login(ctx, config)
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
	// POST http://cx.kingtrans.cn/nclient/CClientPrice
	// GET http://csy.kingtrans.cn/client.jsp
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
		total, err := strconv.ParseFloat(td[c.totalIdx], 64)
		if err != nil {
			log.Err(errors.WithStack(err)).Msgf("%s has err", c.source)
			return true
		}
		weight, err := strconv.ParseFloat(td[c.weightIdx], 64)
		if err != nil {
			log.Err(errors.WithStack(err)).Msgf("%s has err", c.source)
			return true
		}
		price, err := strconv.ParseFloat(td[c.priceIdx], 64)
		if err != nil {
			log.Err(errors.WithStack(err)).Msgf("%s has err", c.source)
			return true
		}
		var fuel float64
		if td[c.fuelIdx] != "含燃油" {
			fuel, err = strconv.ParseFloat(td[c.fuelIdx], 64)
			if err != nil {
				log.Err(errors.WithStack(err)).Msgf("%s has err", c.source)
				return true
			}
		}
		res = append(res, model.Logistics{
			Source: c.source,
			URL:    c.url,
			Method: td[c.methodIdx],
			Weight: weight,
			Total:  total,
			Price:  price,
			Fuel:   fuel,
			Remark: strings.TrimSpace(selection.Next().Next().Text()),
		})
		return true
	})
	return res, nil
}

func (c CXFetcher) login(ctx context.Context, config config.LoginConfig) error {
	resp, err := c.client.PostForm(c.loginUrl, url.Values{
		"userid":   []string{config.Username},
		"password": []string{config.Password},
	})
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return nil
}
