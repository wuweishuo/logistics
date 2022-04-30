package impl

import (
	"github.com/pkg/errors"
	"net/http"
)

type JunYaFetcher struct {
	ZhongFeiFetcher
}

func NewJunYaFetcher() *JunYaFetcher {
	return &JunYaFetcher{
		ZhongFeiFetcher{
			source:   "广州骏亚供应链有限公司",
			url:      "http://193.112.163.140:8082/index.htm",
			signUrl:  "http://193.112.163.140:8082/signin.htm",
			queryUrl: "http://193.112.163.140:8082/priceSearchWeb.htm",
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
		},
	}
}
