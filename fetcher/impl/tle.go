package impl

type TLEFetcher struct {
	IOMSFetcher
}

func NewTLEFetcher() *TLEFetcher {
	return &TLEFetcher{
		NewIOMSFetcher("TLE"),
	}
}
