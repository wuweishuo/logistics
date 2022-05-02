package impl

type JLTYFetcher struct {
	IOMSFetcher
}

func NewJLTYFetcher() *JLTYFetcher {
	return &JLTYFetcher{
		NewIOMSFetcher("JLTY"),
	}
}
