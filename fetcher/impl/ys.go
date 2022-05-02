package impl

import (
	"net/http"
	"net/http/cookiejar"
)

type YSFetcher struct {
	ZhongFeiFetcher
}

func NewYSFetcher() *YSFetcher {
	jar, _ := cookiejar.New(nil)
	return &YSFetcher{
		ZhongFeiFetcher{
			source:   "YS",
			url:      "http://159.75.217.165:8082/login.htm",
			signUrl:  "http://159.75.217.165:8082/signin.htm",
			queryUrl: "http://159.75.217.165:8082/priceSearchQuery.htm",
			client: &http.Client{
				Jar: jar,
			},
		},
	}
}
