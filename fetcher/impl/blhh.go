package impl

import (
	"context"
	"fmt"
	"logistics/config"
	"logistics/enums"
	"logistics/model"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type BailinHuaHuiFetcher struct {
	source   string
	url      string
	queryUrl string
	hasPrice bool
	client   *http.Client
}

func NewBailinHuaHuiFetcher() *BailinHuaHuiFetcher {
	return &BailinHuaHuiFetcher{
		hasPrice: true,
		source:   "柏林华惠",
		url:      "http://blhh.kingtrans.cn/WebPrice?action=list",
		queryUrl: "http://blhh.kingtrans.cn/WebPrice?action=list",
		client:   http.DefaultClient,
	}
}

func (b BailinHuaHuiFetcher) Fetch(ctx context.Context, config config.LoginConfig, countryCode string, weight float64) ([]model.Logistics, error) {
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
	var notFoundRecord bool
	reader.Find(".ks_dl_table1").Each(func(i int, selection *goquery.Selection) {
		if notFoundRecord {
			return
		}
		var td []string
		selection.Find("dd").Each(func(i int, selection *goquery.Selection) {
			td = append(td, strings.TrimSpace(selection.Text()))
		})
		if td[0] == "无记录" {
			notFoundRecord = true
			return
		}
		weight, err := strconv.ParseFloat(td[1], 64)
		if err != nil {
			log.Err(errors.WithStack(err)).Msgf("%s has err", b.source)
			return
		}
		var total, price float64
		if b.hasPrice {
			total, err = strconv.ParseFloat(td[3], 64)
			if err != nil {
				log.Err(errors.WithStack(err)).Msgf("%s has err", b.source)
				return
			}
			price, err = strconv.ParseFloat(td[2], 64)
			if err != nil {
				log.Err(errors.WithStack(err)).Msgf("%s has err", b.source)
				return
			}
		} else {
			total, err = strconv.ParseFloat(td[2], 64)
			if err != nil {
				log.Err(errors.WithStack(err)).Msgf("%s has err", b.source)
				return
			}
		}

		res = append(res, model.Logistics{
			Source: b.source,
			URL:    b.url,
			Method: strings.TrimSpace(td[0]),
			Weight: weight,
			Total:  total,
			Price:  price,
			Remark: strings.TrimSpace(selection.Next().Next().Text()),
		})
	})
	return res, nil
}
