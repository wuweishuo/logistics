package impl

type CSFetcher struct {
	ItDiDaFetcher
}

func NewCSFetcher() *CSFetcher {
	return &CSFetcher{
		NewItDiDaFetcherByDomian("http://csgj.itdida.com"),
	}
}
