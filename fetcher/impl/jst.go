package impl

type JSTFetcher struct {
	ItDiDaFetcher
}

func NewJSTFetcher() *JSTFetcher {
	return &JSTFetcher{
		NewItDiDaFetcherByDomian("http://just.itdida.com"),
	}
}
