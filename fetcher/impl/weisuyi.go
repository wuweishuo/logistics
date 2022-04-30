package impl

import (
	"net/http"
	"net/http/cookiejar"
)

type WeiSuYiFetcher struct {
	ZhongFeiFetcher
}

func NewWeiSuYiFetcher() *WeiSuYiFetcher {
	jar, _ := cookiejar.New(nil)
	return &WeiSuYiFetcher{
		ZhongFeiFetcher{
			source:   "威速易供应链管理",
			url:      "http://www.gdwse.com:8082/index.htm",
			signUrl:  "http://www.gdwse.com:8082/signin.htm",
			queryUrl: "http://www.gdwse.com:8082/priceSearchQuery.htm",
			client: &http.Client{
				Jar: jar,
			},
		},
	}
}
