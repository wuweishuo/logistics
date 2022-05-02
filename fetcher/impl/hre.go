package impl

type HRFetcher struct {
	IOMSFetcher
}

func NewHRFetcher() *HRFetcher {
	return &HRFetcher{
		NewIOMSFetcher("HRE"),
	}
}
