package fetcher

import (
	"context"
	"gopkg.in/yaml.v3"
	"logistics/model"
)

type FetcherConfig interface {
	Parse(node *yaml.Node) error
}

type FetcherFactory interface {
	ConstructFetcher(config FetcherConfig) (Fetcher, error)
	ConstructConfig() FetcherConfig
}

type Fetcher interface {
	Fetch(ctx context.Context, countryCode string, weight float64) ([]model.Logistics, error)
}

type FetcherFactoryRegistry map[string]FetcherFactory

var fetcherFactoryRegistry = make(FetcherFactoryRegistry)

func GetFetcherFactoryRegistry() FetcherFactoryRegistry {
	return fetcherFactoryRegistry
}

func RegisterFetcherFactory(name string, factory FetcherFactory) {
	fetcherFactoryRegistry[name] = factory
}
