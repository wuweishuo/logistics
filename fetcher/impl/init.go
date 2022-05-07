package impl

import "logistics/fetcher"

func init() {
	fetcher.RegisterFetcherFactory("itdida", ItDiDaFetcherFactory{})
	fetcher.RegisterFetcherFactory("i-oms", IOMSFetcherFactory{})
	fetcher.RegisterFetcherFactory("hl", HLFetcherFactory{})
	fetcher.RegisterFetcherFactory("k5", K5FetcherFactory{})
	fetcher.RegisterFetcherFactory("k5-login", K5LoginFetcherFactory{})
	fetcher.RegisterFetcherFactory("san-tong", SanTongFetcherFactory{})
	fetcher.RegisterFetcherFactory("hb", HBFetcherFactory{})
	fetcher.RegisterFetcherFactory("httx", HTTXFetcherFactory{})
	fetcher.RegisterFetcherFactory("bsd", BSDFetcherFactory{})
	fetcher.RegisterFetcherFactory("v5", V5FetcherFactory{})
}
