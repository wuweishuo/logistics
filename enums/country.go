package enums

var Countries = map[string]string{
	"澳大利亚": "AU",
}

var CountryCodes map[string]string

func init() {
	CountryCodes = make(map[string]string, len(Countries))
	for k, v := range Countries {
		CountryCodes[v] = k
	}
}
