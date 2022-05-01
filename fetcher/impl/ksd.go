package impl

import "net/http"

type KsdFetcher struct {
	BailinHuaHuiFetcher
}

func NewKsdFetcher() *KsdFetcher {
	return &KsdFetcher{
		BailinHuaHuiFetcher{
			source:   "KSD",
			url:      "http://ksd.ksds.com.cn/WebPrice?action=list",
			queryUrl: "http://ksd.ksds.com.cn/WebPrice?action=list",
			client:   http.DefaultClient,
		},
	}
}
