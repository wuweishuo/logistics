package impl

import (
	"context"
	"encoding/json"
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
				td = append(td, strings.TrimSpace(selection.Text()))
			})
			if td[0] == "暂无记录" {
				notFoundRecord = true
				return
			}
			total, err := h.parseFloat(ctx, strings.TrimRight(td[9], " RMB"))
			if err != nil {
				return
			}
			fare, err := h.parseFloat(ctx, td[4])
			if err != nil {
				return
			}
			fuel, err := h.parseFloat(ctx, td[5])
			if err != nil {
				return
			}
			other, err := h.parseFloat(ctx, td[7])
			if err != nil {
				return
			}
			queryWeight, err := h.parseFloat(ctx, strings.TrimRight(td[8], " KG"))
			if err != nil {
				return
			}
			res = append(res, model.Logistics{
				Source: "广州中帆国际业务管理系统",
				URL:    "http://gzzf.rtb56.com/usercenter/index.aspx",
				Method: td[2],
				Weight: queryWeight,
				Total:  total,
				Fare:   fare,
				Fuel:   fuel,
				Other:  other,
			})
		} else {
			res[len(res)-1].Remark = strings.TrimSpace(selection.Find("td").Text())
		}
	})
	return res, nil
}

func (h HBFetcher) parseFloat(ctx context.Context, str string) (float64, error) {
	if str == "-" {
		return 0, nil
	}
	float, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Err(errors.WithStack(err)).Msg("")
		return 0, errors.WithStack(err)
	}
	return float, nil
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
	if err != nil {
		return nil, errors.WithStack(err)
	}
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
