package impl

import (
	"net/http"
	"net/http/cookiejar"
)

type TSFetcher struct {
	ZhongFeiFetcher
}

func NewTSFetcher() *TSFetcher {
	jar, _ := cookiejar.New(nil)
	return &TSFetcher{
		ZhongFeiFetcher{
			source:   "TS",
			url:      "http://www.gztsexp.com:8082/login.htm",
			signUrl:  "http://www.gztsexp.com:8082/signin.htm",
			queryUrl: "http://www.gztsexp.com:8082/priceSearchQuery.htm",
			client: &http.Client{
				Jar: jar,
			},
		},
	}
}
