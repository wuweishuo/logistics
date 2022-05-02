package impl

import (
	"net/http"
	"net/http/cookiejar"
)

type WJFetcher struct {
	ZhongFeiFetcher
}

func NewWJFetcher() *WJFetcher {
	jar, _ := cookiejar.New(nil)
	return &WJFetcher{
		ZhongFeiFetcher{
			source:   "WJ",
			url:      "http://111.230.211.49:8082/signin.htm",
			signUrl:  "http://111.230.211.49:8082/signin.htm",
			queryUrl: "http://111.230.211.49:8082/priceSearchQuery.htm",
			client: &http.Client{
				Jar: jar,
			},
		},
	}
}
