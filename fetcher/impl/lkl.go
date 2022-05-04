package impl

type LKLFetcher struct {
	ItDiDaFetcher
}

func NewLKLFetcher() *LKLFetcher {
	return &LKLFetcher{
		NewItDiDaFetcherByDomian("http://lkl.itdida.com"),
	}
}
