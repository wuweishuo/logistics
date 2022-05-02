package impl

import (
	"net/http"
	"net/http/cookiejar"
)

type SWLFetcher struct {
	ZhongFeiFetcher
}

func NewSWLFetcher() *SWLFetcher {
	jar, _ := cookiejar.New(nil)
	return &SWLFetcher{
		ZhongFeiFetcher{
			source:   "SWL",
			url:      "http://193.112.33.159:8082/login.htm",
			signUrl:  "http://193.112.33.159:8082/signin.htm",
			queryUrl: "http://193.112.33.159:8082/priceSearchQuery.htm",
			client: &http.Client{
				Jar: jar,
			},
		},
	}
}
