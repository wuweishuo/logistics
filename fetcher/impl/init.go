package impl

import "logistics/fetcher"

func init() {
	fetcher.Register("mock", MockFetcher{})
	fetcher.Register("san_tong", NewSanTongFetcher())
	fetcher.Register("lian_di", NewLianDiFetcher())
}
