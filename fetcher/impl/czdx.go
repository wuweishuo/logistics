package impl

type CZDXFetcher struct {
	IOMSFetcher
}

func NewCZDXFetcher() *CZDXFetcher {
	return &CZDXFetcher{
		NewIOMSFetcher("CZDX"),
	}
}
