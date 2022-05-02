package impl

type LMTFetcher struct {
	IOMSFetcher
}

func NewLMTFetcher() *LMTFetcher {
	return &LMTFetcher{
		NewIOMSFetcher("LMT"),
	}
}
