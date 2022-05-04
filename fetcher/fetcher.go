package fetcher

import (
	"context"
	"logistics/config"
	"logistics/model"
)

type FetcherFactory interface {
	ConstructFetcher(config interface{}) (Fetcher, error)
	ConstructConfig() interface{}
}

type Fetcher interface {
	Fetch(ctx context.Context, config config.LoginConfig, countryCode string, weight float64) ([]model.Logistics, error)
}

type FetcherFactoryRegistry map[string]FetcherFactory

var fetcherFactoryRegistry = make(FetcherFactoryRegistry)

func GetFetcherFactoryRegistry() FetcherFactoryRegistry {
	return fetcherFactoryRegistry
}

func RegisterFetcherFactory(name string, factory FetcherFactory) {
	fetcherFactoryRegistry[name] = factory
}

var registry = make(Registry)

func GetRegistry() Registry {
	return registry
}

func Register(name string, fetcher Fetcher) {
	registry[name] = fetcher
}

type Registry map[string]Fetcher
