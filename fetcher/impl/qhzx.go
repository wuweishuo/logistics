package impl

type QHZXFetcher struct {
	ItDiDaFetcher
}

func NewQHZXFetcher() *QHZXFetcher {
	return &QHZXFetcher{
		NewItDiDaFetcherByDomian("http://zxi.itdida.com"),
	}
}
