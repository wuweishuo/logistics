package impl

import (
	"net/http"
	"net/http/cookiejar"
)

type ZWJYFetcher struct {
	CXFetcher
}

func NewZWJYFetcher() *ZWJYFetcher {
	jar, _ := cookiejar.New(nil)
	return &ZWJYFetcher{
		CXFetcher{
			source:    "ZWJY",
			url:       "http://zwjy.kingtrans.cn/new_index.jsp",
			loginUrl:  "http://zwjy.kingtrans.net/CUserLogon?action=logon",
			queryUrl:  "http://zwjy.kingtrans.net/nclient/CClientPrice?action=getAnalyse",
			methodIdx: 0,
			weightIdx: 2,
			totalIdx:  3,
			priceIdx:  4,
			fuelIdx:   5,
			client: &http.Client{
				Jar: jar,
			},
		},
	}
}
