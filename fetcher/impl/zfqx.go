package impl

import (
	"net/http"
	"net/http/cookiejar"
)

type ZFQXFetcher struct {
	CXFetcher
}

func NewZFQXFetcher() *ZFQXFetcher {
	jar, _ := cookiejar.New(nil)
	return &ZFQXFetcher{
		CXFetcher{
			source:    "ZFQX",
			url:       "http://zfqx.kingtrans.cn/new_index.jsp",
			loginUrl:  "http://zfqx.kingtrans.net/nclient/Logon?action=logon",
			queryUrl:  "http://zfqx.kingtrans.net/nclient/CClientPrice?action=getAnalyse",
			methodIdx: 0,
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
