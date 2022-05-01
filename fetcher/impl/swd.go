package impl

import "net/http"

type SwdFetcher struct {
	BailinHuaHuiFetcher
}

func NewSwdFetcher() *SwdFetcher {
	return &SwdFetcher{
		BailinHuaHuiFetcher{
			source:   "SWD",
			url:      "http://gzwsd.kingtrans.cn/WebPrice?action=list",
			queryUrl: "http://gzwsd.kingtrans.cn/WebPrice?action=list",
			client:   http.DefaultClient,
		},
	}
}
