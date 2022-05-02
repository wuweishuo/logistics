package impl

type JSTFetcher struct {
	ItDiDaFetcher
}

func NewJSTFetcher() *JSTFetcher {
	return &JSTFetcher{
		NewItDiDaFetcher("JST", "http://just.itdida.com/itdida-flash/website/landing",
			"http://just.itdida.com/itdida-api/login", "http://just.itdida.com/itdida-api/flash/price/query"),
	}
}
