package impl

type JRFetcher struct {
	IOMSFetcher
}

func NewJRFetcher() *JRFetcher {
	return &JRFetcher{
		NewIOMSFetcher("JRE"),
	}
}
