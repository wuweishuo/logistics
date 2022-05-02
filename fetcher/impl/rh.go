package impl

import (
	"net/http"
	"net/http/cookiejar"
)

type RHFetcher struct {
	ZhongFeiFetcher
}

func NewRHFetcher() *RHFetcher {
	jar, _ := cookiejar.New(nil)
	return &RHFetcher{
		ZhongFeiFetcher{
			source:   "RH",
			url:      "http://www.rh168.com:8082/signin.htm",
			signUrl:  "http://www.rh168.com:8082/signin.htm",
			queryUrl: "http://www.rh168.com:8082/priceSearchQuery.htm",
			client: &http.Client{
				Jar: jar,
			},
		},
	}
}
