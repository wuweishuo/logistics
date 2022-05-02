package impl

type QHZXFetcher struct {
	ItDiDaFetcher
}

func NewQHZXFetcher() *QHZXFetcher {
	return &QHZXFetcher{
		NewItDiDaFetcher("QHZX", "http://zxi.itdida.com/itdida-flash/website/landing",
			"http://zxi.itdida.com/itdida-api/login", "http://zxi.itdida.com/itdida-api/flash/price/query"),
	}
}
