package main

import (
	"context"
	"flag"
	"github.com/jinzhu/configor"
	"github.com/rs/zerolog/log"
	"logistics/config"
	"logistics/fetcher"
	_ "logistics/fetcher/impl"
	"logistics/model"
	"sort"
)

func main() {
	var countryCode string
	flag.StringVar(&countryCode, "country_code", "AD", "input your destination")
	var weight float64
	flag.Float64Var(&weight, "weight", 1, "input your weight")
	var configFile string
	flag.StringVar(&configFile, "config_file", "./config/config.yml", "input your config file")
	ctx := context.Background()
	var c config.Config
	err := configor.Load(&c, configFile)
	if err != nil {
		panic(err)
	}
	var res []model.Logistics
	for name, f := range fetcher.GetRegistry() {
		if _, ok := c.Logins[name]; !ok {
			log.Error().Msgf("%s hasn't config", name)
			continue
		}
		data, err := f.Fetch(ctx, c.Logins[name], countryCode, weight)
		if err != nil {
			log.Err(err).Msgf("%s has err", name)
			continue
		}
		res = append(res, data...)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Total < res[j].Total
	})
	log.Info().Msgf("out data:%v", res)
}
