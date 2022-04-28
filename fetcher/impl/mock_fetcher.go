package impl

import (
	"context"
	"logistics/config"
	"logistics/fetcher"
	"logistics/model"
)

func init() {
	fetcher.Register("mock", MockFetcher{})
}

type MockFetcher struct{}

func (m MockFetcher) Fetch(ctx context.Context, config config.LoginConfig, countryCode string, weight float64) ([]model.Logistics, error) {
	return []model.Logistics{
		{
			Source: "mock",
			Weight: 1,
			Total:  1,
			Remark: "remark",
		},
	}, nil
}
