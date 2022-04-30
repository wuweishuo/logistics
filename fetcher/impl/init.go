package impl

import "logistics/fetcher"

func init() {
	fetcher.Register("mock", MockFetcher{})
	fetcher.Register("santong", NewSanTongFetcher())
	fetcher.Register("liandi", NewLianDiFetcher())
	fetcher.Register("zhongfei", NewZhongFeiFetcher())
	fetcher.Register("junya", NewJunYaFetcher())
	fetcher.Register("weisuyi", NewWeiSuYiFetcher())
	fetcher.Register("hb", NewHBFetcher())
}
