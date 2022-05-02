package impl

type LKLFetcher struct {
	ItDiDaFetcher
}

func NewLKLFetcher() *LKLFetcher {
	return &LKLFetcher{
		NewItDiDaFetcher("LKL", "http://lkl.itdida.com/itdida-flash/website/landing",
			"http://lkl.itdida.com/itdida-api/login", "http://lkl.itdida.com/itdida-api/flash/price/query"),
	}
}
