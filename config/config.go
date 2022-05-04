package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"logistics/fetcher"
)

type Config map[string]map[string]fetcher.FetcherConfig

func (c *Config) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.New("config not right")
	}
	m := make(map[string]map[string]fetcher.FetcherConfig, len(node.Content))
	for i := 0; i < len(node.Content); i += 2 {
		nk := node.Content[i]
		if nk.ShortTag() != "!!str" {
			return errors.New("config not right")
		}
		factory, ok := fetcher.GetFetcherFactoryRegistry()[nk.Value]
		if !ok {
			continue
		}
		nv := node.Content[i+1]
		if nv.Kind != yaml.MappingNode {
			return errors.New("config not right")
		}
		m[nk.Value] = make(map[string]fetcher.FetcherConfig, len(node.Content))
		for j := 0; j < len(nv.Content); j += 2 {
			if nv.Content[j].ShortTag() != "!!str" {
				return errors.New("config not right")
			}
			if nv.Content[j+1].Kind != yaml.MappingNode {
				return errors.New("config not right")
			}
			config := factory.ConstructConfig()
			err := config.Parse(nv.Content[j+1])
			if err != nil {
				return err
			}
			m[nk.Value][nv.Content[j].Value] = config
		}
	}
	*c = m
	return nil
}
