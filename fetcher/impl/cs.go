package impl

type CSFetcher struct {
	ItDiDaFetcher
}

func NewCSFetcher() *CSFetcher {
	return &CSFetcher{
		NewItDiDaFetcher("CS", "http://csgj.itdida.com/itdida-flash/website/landing",
			"http://csgj.itdida.com/itdida-api/login", "http://csgj.itdida.com/itdida-api/flash/price/query"),
	}
}
