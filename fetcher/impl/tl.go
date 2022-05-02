package impl

import (
	"net/http"
	"net/http/cookiejar"
)

type TLFetcher struct {
	ZhongFeiFetcher
}

func NewTLFetcher() *TLFetcher {
	jar, _ := cookiejar.New(nil)
	return &TLFetcher{
		ZhongFeiFetcher{
			source:   "天蓝",
			url:      "http://139.159.213.246:8082/login.htm",
			signUrl:  "http://139.159.213.246:8082/signin.htm",
			queryUrl: "http://139.159.213.246:8082/priceSearchQuery.htm",
			client: &http.Client{
				Jar: jar,
			},
		},
	}
}
