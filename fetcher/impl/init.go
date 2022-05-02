package impl

import "logistics/fetcher"

func init() {
	fetcher.Register("mock", MockFetcher{})
	fetcher.Register("santong", NewSanTongFetcher())
	fetcher.Register("hb", NewHBFetcher())

	// k5 不需登陆
	fetcher.Register("blhh", NewBailinHuaHuiFetcher())
	fetcher.Register("bk", NewBkFetcher())
	fetcher.Register("ksd", NewKsdFetcher())
	fetcher.Register("swd", NewSwdFetcher())
	fetcher.Register("xf", NewXfFetcher())
	// k5
	fetcher.Register("cx", NewCXFetcher())
	fetcher.Register("th", NewTHFetcher())
	fetcher.Register("szxy", NewSZXYFetcher())
	fetcher.Register("zwjy", NewZWJYFetcher())
	fetcher.Register("zfqx", NewZFQXFetcher())

	// 华磊科技
	fetcher.Register("zhongfei", NewZhongFeiFetcher())
	fetcher.Register("junya", NewJunYaFetcher())
	fetcher.Register("weisuyi", NewWeiSuYiFetcher())
	fetcher.Register("ts", NewTSFetcher())
	fetcher.Register("rh", NewRHFetcher())
	fetcher.Register("swl", NewSWLFetcher())
	fetcher.Register("tl", NewTLFetcher())
	fetcher.Register("ys", NewYSFetcher())

	// i-oms
	fetcher.Register("lde", NewLDEFetcher())
	fetcher.Register("czdx", NewCZDXFetcher())
	fetcher.Register("hre", NewHRFetcher())
	fetcher.Register("jlty", NewJLTYFetcher())
	fetcher.Register("jre", NewJRFetcher())
	fetcher.Register("lmt", NewLMTFetcher())
	fetcher.Register("tle", NewTLEFetcher())
	fetcher.Register("yht", NewYHTFetcher())
	fetcher.Register("zfe", NewZFEFetcher())
}
