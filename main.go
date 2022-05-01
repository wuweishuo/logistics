package main

import (
	"context"
	"flag"
	"fmt"
	"logistics/config"
	"logistics/fetcher"
	_ "logistics/fetcher/impl"
	"logistics/model"
	"os"
	"sort"
	"strconv"

	"github.com/jinzhu/configor"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Stack().Caller().Logger()
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var countryCode string
	flag.StringVar(&countryCode, "country_code", "US", "input your destination")
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
	type fetcherResult struct {
		name string
		data []model.Logistics
		err  error
	}
	channel := make(chan fetcherResult)
	for name := range fetcher.GetRegistry() {
		go func(name string) {
			defer func() {
				err := recover()
				if err != nil {
					channel <- fetcherResult{
						name: name,
						data: nil,
						err:  errors.Errorf("panic:%+v", err),
					}
				}

			}()
			if _, ok := c.Logins[name]; !ok {
				channel <- fetcherResult{
					name: name,
					data: nil,
					err:  errors.Errorf("%s hasn't config", name),
				}
				return
			}
			data, err := fetcher.GetRegistry()[name].Fetch(ctx, c.Logins[name], countryCode, weight)
			channel <- fetcherResult{
				name: name,
				data: data,
				err:  err,
			}
		}(name)
	}
	for i := 0; i < len(fetcher.GetRegistry()); i++ {
		result := <-channel
		if result.err != nil {
			log.Err(result.err).Stack().Msgf("%s has err", result.name)
			continue
		}
		res = append(res, result.data...)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Total < res[j].Total
	})

	writer := tablewriter.NewWriter(os.Stdout)
	writer.SetHeader([]string{"来源", "url", "渠道", "重量", "总价", "单价", "运费", "燃油", "其他杂费", "备注"})
	writer.SetFooter([]string{"", "", "", "", "", "", "", "", "total", strconv.Itoa(len(res))})
	data := make([][]string, 0, len(res))
	for _, d := range res {
		data = append(data, []string{
			d.Source,
			d.URL,
			d.Method,
			fmt.Sprintf("%v", d.Weight),
			fmt.Sprintf("%v", d.Total),
			fmt.Sprintf("%v", d.Price),
			fmt.Sprintf("%v", d.Fare),
			fmt.Sprintf("%v", d.Fuel),
			fmt.Sprintf("%v", d.Other),
			"",
		})
	}
	writer.AppendBulk(data)
	writer.Render()
}
