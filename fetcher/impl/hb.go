package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"logistics/config"
	"logistics/model"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

type HBFetcher struct {
	client *http.Client
}

func NewHBFetcher() *HBFetcher {
	jar, _ := cookiejar.New(nil)
	return &HBFetcher{
		client: &http.Client{
			Jar: jar,
		},
	}
}

func (h HBFetcher) Fetch(ctx context.Context, config config.LoginConfig, countryCode string, weight float64) ([]model.Logistics, error) {
	err := h.login(ctx, config)
	if err != nil {
		return nil, err
	}
	values, err := h.getParameters(ctx, countryCode, weight)
	if err != nil {
		return nil, err
	}
	resp, err := h.client.PostForm("http://gzzf.rtb56.com/usercenter/querytools/fee_trail.aspx", values)
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
	reader.Find("table tbody tr").Each(func(i int, selection *goquery.Selection) {
		if i == 0 || notFoundRecord {
			return
		}
		if i%2 != 0 {
			var td []string
			selection.Find("td").Each(func(i int, selection *goquery.Selection) {
				td = append(td, selection.Text())
			})
			if td[0] == "暂无记录" {
				notFoundRecord = true
				return
			}
			total, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimRight(td[9], "RMB")), 10)
			if err != nil {
				log.Err(errors.WithStack(err)).Msg("")
				return
			}
			price, err := strconv.ParseFloat(td[4], 10)
			if err != nil {
				log.Err(errors.WithStack(err)).Msg("")
				return
			}
			fare, err := strconv.ParseFloat(td[4], 10)
			if err != nil {
				log.Err(errors.WithStack(err)).Msg("")
				return
			}
			fuel, err := strconv.ParseFloat(td[5], 10)
			if err != nil {
				log.Err(errors.WithStack(err)).Msg("")
				return
			}
			other, err := strconv.ParseFloat(td[7], 10)
			if err != nil {
				log.Err(errors.WithStack(err)).Msg("")
				return
			}
			res = append(res, model.Logistics{
				Source: "广州中帆国际业务管理系统",
				URL:    "http://gzzf.rtb56.com/usercenter/index.aspx",
				Method: td[2],
				Weight: weight,
				Total:  total,
				Price:  price,
				Fare:   fare,
				Fuel:   fuel,
				Other:  other,
			})
		} else {
			res[len(res)-1].Remark = selection.Find("td").Text()
		}
	})
	return res, nil
}

func (h HBFetcher) getParameters(ctx context.Context, countryCode string, weight float64) (url.Values, error) {
	resp, err := h.client.Get("http://gzzf.rtb56.com/usercenter/querytools/fee_trail.aspx")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	reader, err := goquery.NewDocumentFromReader(resp.Body)
	return url.Values{
		"txtCountry":           []string{countryCode},
		"txtWeight":            []string{fmt.Sprintf("%v", weight)},
		"rbCargoType":          []string{"W"},
		"__EVENTVALIDATION":    []string{h.getValue(ctx, reader, "#__EVENTVALIDATION")},
		"__VIEWSTATEGENERATOR": []string{h.getValue(ctx, reader, "#__VIEWSTATEGENERATOR")},
		"__VIEWSTATE":          []string{h.getValue(ctx, reader, "#__VIEWSTATE")},
		"__EVENTTARGET":        []string{"btnSubmit"},
	}, nil
}

func (h HBFetcher) getValue(ctx context.Context, reader *goquery.Document, id string) string {
	attrs := reader.Find(id).Get(0).Attr
	for _, attr := range attrs {
		if attr.Key == "value" {
			return attr.Val
		}
	}
	return ""
}

func (h HBFetcher) login(ctx context.Context, config config.LoginConfig) error {
	resp, err := h.client.PostForm("http://gzzf.rtb56.com/tools/submit_ajax.ashx?action=user_login_g", url.Values{
		"txtUserName": []string{config.Username},
		"txtPassword": []string{config.Password},
	})
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var m map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&m)
	if err != nil {
		return errors.WithStack(err)
	}
	if m["status"] != float64(1) {
		return errors.New(m["msg"].(string))
	}
	return nil
}
