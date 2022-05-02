package impl

type ZFEFetcher struct {
	IOMSFetcher
}

func NewZFEFetcher() *ZFEFetcher {
	return &ZFEFetcher{
		NewIOMSFetcher("ZFE"),
	}
}
