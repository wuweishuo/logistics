package impl

import "net/http"

type BkFetcher struct {
	BailinHuaHuiFetcher
}

func NewBkFetcher() *BkFetcher {
	return &BkFetcher{
		BailinHuaHuiFetcher{
			source:   "百科国际物流",
			url:      "http://bk.kingtrans.cn/WebPrice?action=list",
			queryUrl: "http://bk.kingtrans.cn/WebPrice?action=list",
			client:   http.DefaultClient,
		},
	}
}
