package impl

import (
	"net/http"
	"net/http/cookiejar"
)

type THFetcher struct {
	CXFetcher
}

func NewTHFetcher() *THFetcher {
	jar, _ := cookiejar.New(nil)
	return &THFetcher{
		CXFetcher{
			source:    "TH",
			url:       "http://sky.kingtrans.net/old_index.jsp?retry_reason=PASSWORD",
			loginUrl:  "http://sky.kingtrans.net/nclient/Logon?action=logon",
			queryUrl:  "http://sky.kingtrans.net/nclient/CClientPrice?action=getAnalyse",
			methodIdx: 0,
			weightIdx: 2,
			totalIdx:  4,
			priceIdx:  5,
			fuelIdx:   6,
			client: &http.Client{
				Jar: jar,
			},
		},
	}
}
