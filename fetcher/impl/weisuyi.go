package impl

import (
	"github.com/pkg/errors"
	"net/http"
)

type WeiSuYiFetcher struct {
	ZhongFeiFetcher
}

func NewWeiSuYiFetcher() *WeiSuYiFetcher {
	return &WeiSuYiFetcher{
		ZhongFeiFetcher{
			source:   "威速易供应链管理",
			url:      "http://www.gdwse.com:8082/index.htm",
			signUrl:  "http://www.gdwse.com:8082/signin.htm",
			queryUrl: "http://www.gdwse.com:8082/priceSearchQuery.htm",
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
