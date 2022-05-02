package impl

import "logistics/fetcher"

func init() {
	//fetcher.Register("mock", MockFetcher{})
	fetcher.Register("三通", NewSanTongFetcher())
	fetcher.Register("环邦", NewHBFetcher())
	fetcher.Register("汇通天下", NewHTTXFetcher())

	// k5 不需登陆
	fetcher.Register("柏林华惠", NewBailinHuaHuiFetcher())
	fetcher.Register("百科", NewBkFetcher())
	fetcher.Register("凯时达", NewKsdFetcher())
	fetcher.Register("森威达", NewSwdFetcher())
	fetcher.Register("喜丰", NewXfFetcher())
	// k5
	fetcher.Register("宸轩", NewCXFetcher())
	fetcher.Register("天豪", NewTHFetcher())
	fetcher.Register("深圳迅一", NewSZXYFetcher())
	fetcher.Register("中外急运", NewZWJYFetcher())
	fetcher.Register("逐风前行", NewZFQXFetcher())
	//fetcher.Register("传送易", NewCSYFetcher())

	// 华磊科技
	fetcher.Register("中飞", NewZhongFeiFetcher())
	fetcher.Register("骏亚", NewJunYaFetcher())
	fetcher.Register("腾顺", NewTSFetcher())
	fetcher.Register("融航", NewRHFetcher())
	fetcher.Register("顺万里", NewSWLFetcher())
	fetcher.Register("天蓝", NewTLFetcher())
	fetcher.Register("云顺", NewYSFetcher())
	fetcher.Register("威捷", NewWJFetcher())

	// i-oms
	fetcher.Register("联递", NewLDEFetcher())
	fetcher.Register("创智", NewCZDXFetcher())
	fetcher.Register("华仁", NewHRFetcher())
	fetcher.Register("金龙", NewJLTYFetcher())
	fetcher.Register("巨人", NewJRFetcher())
	fetcher.Register("蓝马特", NewLMTFetcher())
	fetcher.Register("天霖", NewTLNFetcher())
	fetcher.Register("永恒泰", NewYHTFetcher())
	fetcher.Register("深圳中飞", NewZFEFetcher())

	// ItDiDa
	fetcher.Register("一代", NewYDFetcher())
	fetcher.Register("超顺", NewCSFetcher())
	fetcher.Register("加时特", NewJSTFetcher())
	fetcher.Register("乐凯龙", NewLKLFetcher())
	fetcher.Register("众鑫", NewQHZXFetcher())
}
