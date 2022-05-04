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

type ZhongFeiFetcher struct {
	source   string
	url      string
	signUrl  string
	queryUrl string
	client   *http.Client
}

func NewZhongFeiFetcher() *ZhongFeiFetcher {
	jar, _ := cookiejar.New(nil)
	return &ZhongFeiFetcher{
		source:   "中飞国际",
		url:      "http://193.112.219.243:8082/index.htm",
		signUrl:  "http://193.112.219.243:8082/signin.htm",
		queryUrl: "http://193.112.219.243:8082/priceSearchQuery.htm",
		client: &http.Client{
			Jar: jar,
		},
	}
}

func (z ZhongFeiFetcher) Fetch(ctx context.Context, config config.LoginConfig, countryCode string, weight float64) ([]model.Logistics, error) {
	err := z.login(ctx, config)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", z.queryUrl, strings.NewReader(url.Values{
		"country":   []string{countryCode},
		"cargoType": []string{"P"},
		"weight":    []string{fmt.Sprintf("%v", weight)},
	}.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := z.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	//POST http://193.112.219.243:8082/priceSearchQuery.htm
	//GET http://111.230.211.49:8082/login.htm
	if resp.Request.Method == http.MethodGet {
		return nil, errors.New("登录失败")
	}
	reader, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var res []model.Logistics
	reader.Find("table tbody tr").Each(func(i int, selection *goquery.Selection) {
		if i%2 == 0 {
			var td []string
			selection.Find("td").Each(func(i int, selection *goquery.Selection) {
				td = append(td, strings.TrimSpace(selection.Text()))
			})
			weight, err := strconv.ParseFloat(td[2], 64)
			if err != nil {
				log.Err(errors.WithStack(err)).Msg("")
				return
			}
			total, err := strconv.ParseFloat(td[3], 64)
			if err != nil {
				log.Err(errors.WithStack(err)).Msg("")
				return
			}
			price, err := strconv.ParseFloat(td[4], 64)
			if err != nil {
				log.Err(errors.WithStack(err)).Msg("")
				return
			}
			fare, err := strconv.ParseFloat(td[5], 64)
			if err != nil {
				log.Err(errors.WithStack(err)).Msg("")
				return
			}
			fuel, err := strconv.ParseFloat(td[6], 64)
			if err != nil {
				log.Err(errors.WithStack(err)).Msg("")
				return
			}
			other, err := strconv.ParseFloat(td[7], 64)
			if err != nil {
				log.Err(errors.WithStack(err)).Msg("")
				return
			}
			res = append(res, model.Logistics{
				Source: z.source,
				URL:    z.url,
				Method: td[0],
				Weight: weight,
				Total:  total,
				Price:  price,
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

func (z ZhongFeiFetcher) login(ctx context.Context, config config.LoginConfig) error {
	resp, err := z.client.PostForm(z.signUrl, url.Values{
		"username": []string{config.Username},
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
