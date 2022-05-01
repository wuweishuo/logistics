package impl

import (
	"net/http"
	"net/http/cookiejar"
)

type SZXYFetcher struct {
	CXFetcher
}

func NewSZXYFetcher() *SZXYFetcher {
	jar, _ := cookiejar.New(nil)
	return &SZXYFetcher{
		CXFetcher{
			source:    "SZXY",
			url:       "http://xyex.kingtrans.cn/new_index.jsp?retry_reason=PASSWORD",
			loginUrl:  "http://xyex.kingtrans.cn/nclient/Logon?action=logon",
			queryUrl:  "http://xyex.kingtrans.cn/nclient/CClientPrice?action=getAnalyse",
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
