package impl

import (
	"net/http"
	"net/http/cookiejar"
)

type CSYFetcher struct {
	CXFetcher
}

func NewCSYFetcher() *CSYFetcher {
	jar, _ := cookiejar.New(nil)
	return &CSYFetcher{
		CXFetcher{
			source:    "CX",
			url:       "http://csy.kingtrans.cn/old_index.jsp?retry_reason=PASSWORD",
			loginUrl:  "http://csy.kingtrans.cn/client/Logon?action=logon",
			queryUrl:  "http://csy.kingtrans.cn/nclient/CClientPrice?action=getAnalyse",
			methodIdx: 1,
			weightIdx: 3,
			totalIdx:  5,
			priceIdx:  6,
			fuelIdx:   7,
			client: &http.Client{
				Jar: jar,
			},
		},
	}
}
