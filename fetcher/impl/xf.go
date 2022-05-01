package impl

import "net/http"

type XfFetcher struct {
	BailinHuaHuiFetcher
}

func NewXfFetcher() *XfFetcher {
	return &XfFetcher{
		BailinHuaHuiFetcher{
			source:   "XF",
			url:      "http://xf.kingtrans.cn/WebPrice?action=list",
			queryUrl: "http://xf.kingtrans.cn/WebPrice?action=list",
			client:   http.DefaultClient,
		},
	}
}
