package impl

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"logistics/config"
	"logistics/model"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ZhongFeiFetcher struct {
	client *http.Client
}

func NewZhongFeiFetcher() *ZhongFeiFetcher {
	return &ZhongFeiFetcher{
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

func (z ZhongFeiFetcher) Fetch(ctx context.Context, config config.LoginConfig, countryCode string, weight float64) ([]model.Logistics, error) {
	cookies, err := z.getCookies(ctx, config)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "http://193.112.219.243:8082/priceSearchQuery.htm", strings.NewReader(url.Values{
		"country":   []string{countryCode},
		"cargoType": []string{"P"},
		"weight":    []string{fmt.Sprintf("%v", weight)},
	}.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := z.client.Do(req)
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
	reader.Find("table tbody tr").Each(func(i int, selection *goquery.Selection) {
		if i%2 == 0 {
			var td []string
			selection.Find("td").Each(func(i int, selection *goquery.Selection) {
				td = append(td, selection.Text())
			})
			weight, err := strconv.ParseFloat(td[2], 10)
			if err != nil {
				log.Err(errors.WithStack(err)).Msg("")
				return
			}
			total, err := strconv.ParseFloat(td[3], 10)
			if err != nil {
				log.Err(errors.WithStack(err)).Msg("")
				return
			}
			price, err := strconv.ParseFloat(td[4], 10)
			if err != nil {
				log.Err(errors.WithStack(err)).Msg("")
				return
			}
			fare, err := strconv.ParseFloat(td[5], 10)
			if err != nil {
				log.Err(errors.WithStack(err)).Msg("")
				return
			}
			fuel, err := strconv.ParseFloat(td[6], 10)
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
				Source: "中飞国际",
				URL:    "http://193.112.219.243:8082/priceSearchQuery.htm",
				Method: td[0],
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

func (z ZhongFeiFetcher) getCookies(ctx context.Context, config config.LoginConfig) ([]*http.Cookie, error) {
	resp, err := z.client.PostForm("http://193.112.219.243:8082/signin.htm", url.Values{
		"username": []string{config.Username},
		"password": []string{config.Password},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return resp.Request.Cookies(), nil
}
