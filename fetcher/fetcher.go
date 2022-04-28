package fetcher

import (
	"context"
	"logistics/config"
	"logistics/model"
)

type Fetcher interface {
	Fetch(ctx context.Context, config config.LoginConfig, countryCode string, weight float64) ([]model.Logistics, error)
}

var registry = make(Registry)

func GetRegistry() Registry {
	return registry
}

func Register(name string, fetcher Fetcher) {
	registry[name] = fetcher
}

type Registry map[string]Fetcher
