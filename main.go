package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"logistics/config"
	"logistics/enums"
	"logistics/fetcher"
	_ "logistics/fetcher/impl"
	"logistics/model"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.Logger = zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Stack().Caller().Logger().Level(zerolog.Disabled)
	if os.Getenv("LOG_DEBUG") == "debug" {
		log.Logger = log.Logger.Level(zerolog.InfoLevel)
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	var cmd string
	flag.StringVar(&cmd, "c", "query", "input your cmd")
	var countryName string
	flag.StringVar(&countryName, "countryName", "", "input your destination")
	var countryCode string
	flag.StringVar(&countryCode, "countryCode", "US", "input your destination")
	var weight float64
	flag.Float64Var(&weight, "weight", 1, "input your weight")
	var configFile string
	flag.StringVar(&configFile, "configFile", "config.yml", "input your config file")
	var sourceStr string
	flag.StringVar(&sourceStr, "sources", "", "input your source")
	flag.Parse()
	switch cmd {
	case "listCountry":
		listCountry(countryName)
	case "query":
		var sources []string
		if sourceStr != "" {
			sources = strings.Split(sourceStr, ",")
		}
		query(countryCode, weight, configFile, sources)
	default:
		log.Info().Msg("-c option [listCountry, query]")
	}
}

func listCountry(countryName string) {
	writer := tablewriter.NewWriter(os.Stdout)
	writer.SetHeader([]string{"国家/地区名称", "国家/地区编码"})
	var data [][]string
	for name, code := range enums.Countries {
		if countryName != "" && !strings.Contains(name, countryName) {
			continue
		}
		data = append(data, []string{
			name,
			code,
		})
	}
	writer.AppendBulk(data)
	writer.Render()
}

func query(countryCode string, weight float64, configFile string, sources []string) {
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()

	var c config.Config
	bs, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bs, &c)
	if err != nil {
		panic(err)
	}

	var res []model.Logistics
	type fetcherResult struct {
		name      string
		data      []model.Logistics
		err       error
		startTime time.Time
	}
	channel := make(chan fetcherResult)
	registry := fetcher.GetFetcherFactoryRegistry()
	type task struct {
		f    fetcher.Fetcher
		name string
	}
	var count int
	names := make(map[string]bool)
	for k := range c {
		factory, ok := registry[k]
		if !ok {
			continue
		}
		taskChannel := make(chan task, len(c[k]))
		go func(ch chan task) {
			for true {
				select {
				case t, ok := <-ch:
					if !ok {
						return
					}
					func(t task) {
						startTime := time.Now()
						defer func() {
							err := recover()
							if err != nil {
								channel <- fetcherResult{
									name:      t.name,
									data:      nil,
									err:       errors.Errorf("panic:%+v", err),
									startTime: startTime,
								}
							}
						}()
						data, err := t.f.Fetch(ctx, countryCode, weight)
						channel <- fetcherResult{
							name:      t.name,
							data:      data,
							err:       err,
							startTime: startTime,
						}
					}(t)
				case <-ctx.Done():
					return
				}
			}
		}(taskChannel)
		for name, fetcherConfig := range c[k] {
			match := len(sources) == 0
			if len(sources) != 0 {
				for _, source := range sources {
					if source == name {
						match = true
					}
				}
			}
			if match {
				f, err := factory.ConstructFetcher(fetcherConfig)
				if err != nil {
					log.Err(err).Stack().Msgf("%s construct err", name)
					continue
				}
				taskChannel <- task{
					name: name,
					f:    f,
				}
				count++
				names[name] = true
			}
		}
		close(taskChannel)
	}
	var errorName []string
	for i := 0; i < count; i++ {
		select {
		case result := <-channel:
			delete(names, result.name)
			if result.err != nil {
				errorName = append(errorName, result.name)
				log.Err(result.err).Stack().Msgf("num:%d name:%s has err, cost time:%v s", i, result.name, time.Since(result.startTime).Seconds())
				continue
			}
			for _, data := range result.data {
				data.Source = result.name
				res = append(res, data)
			}
			log.Info().Msgf("num:%d name:%s success,cost time:%v s", i, result.name, time.Since(result.startTime).Seconds())
		case <-ctx.Done():
			log.Info().Msgf("num:%d timeout", i)
			break
		}
	}
	close(channel)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Total < res[j].Total
	})
	var timeout []string
	for name := range names {
		timeout = append(timeout, name)
	}

	writer := tablewriter.NewWriter(os.Stdout)
	writer.SetHeader([]string{"来源", "url", "渠道", "重量", "总价", "单价", "运费", "燃油", "其他杂费"})
	writer.SetFooter([]string{
		"",
		"timeout", strings.Join(timeout, ","),
		"errors", strings.Join(errorName, ","),
		"cost_time", fmt.Sprintf("%vs", time.Since(startTime).Seconds()),
		"total", strconv.Itoa(len(res)),
	})
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
		})
	}
	writer.AppendBulk(data)
	writer.Render()
}

func getConfig() interface{} {
	var fieldDef []reflect.StructField
	for name, fetcherFactory := range fetcher.GetFetcherFactoryRegistry() {
		fieldDef = append(fieldDef, reflect.StructField{
			Name: strings.ToUpper(name),
			Type: reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(fetcherFactory.ConstructConfig())),
			Tag:  reflect.StructTag(fmt.Sprintf(`yaml:%s`, name)),
		})
	}
	typ := reflect.StructOf(fieldDef)
	return reflect.New(typ).Interface()
}
