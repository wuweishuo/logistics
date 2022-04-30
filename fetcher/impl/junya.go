package impl

import (
	"net/http"
	"net/http/cookiejar"
)

type JunYaFetcher struct {
	ZhongFeiFetcher
}

func NewJunYaFetcher() *JunYaFetcher {
	jar, _ := cookiejar.New(nil)
	return &JunYaFetcher{
		ZhongFeiFetcher{
			source:   "广州骏亚供应链有限公司",
			url:      "http://193.112.163.140:8082/index.htm",
			signUrl:  "http://193.112.163.140:8082/signin.htm",
			queryUrl: "http://193.112.163.140:8082/priceSearchWeb.htm",
			client: &http.Client{
				Jar: jar,
			},
		},
	}
}
