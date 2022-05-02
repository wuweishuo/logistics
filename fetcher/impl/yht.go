package impl

type YHTFetcher struct {
	IOMSFetcher
}

func NewYHTFetcher() *YHTFetcher {
	return &YHTFetcher{
		NewIOMSFetcher("YHT"),
	}
}
