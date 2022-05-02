package impl

type TLNFetcher struct {
	IOMSFetcher
}

func NewTLNFetcher() *TLNFetcher {
	return &TLNFetcher{
		NewIOMSFetcher("TLN"),
	}
}
