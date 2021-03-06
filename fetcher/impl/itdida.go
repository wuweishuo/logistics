package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"logistics/fetcher"
	"logistics/model"
	"net/http"

	"github.com/pkg/errors"
)

type ItDiDaFetcherConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Domain   string `yaml:"domain"`
}

func (i *ItDiDaFetcherConfig) Parse(value *yaml.Node) error {
	return value.Decode(i)
}

type ItDiDaFetcherFactory struct{}

func (i ItDiDaFetcherFactory) ConstructFetcher(config fetcher.FetcherConfig) (fetcher.Fetcher, error) {
	c, ok := config.(*ItDiDaFetcherConfig)
	if !ok {
		return nil, errors.New("config not right")
	}
	return NewItDiDaFetcher(*c), nil
}

func (i ItDiDaFetcherFactory) ConstructConfig() fetcher.FetcherConfig {
	return &ItDiDaFetcherConfig{}
}

type ItDiDaFetcher struct {
	config   ItDiDaFetcherConfig
	url      string
	queryUrl string
	client   *http.Client
}

func NewItDiDaFetcher(config ItDiDaFetcherConfig) ItDiDaFetcher {
	return ItDiDaFetcher{
		config:   config,
		url:      config.Domain + "/itdida-flash/website/landing",
		queryUrl: config.Domain + "/itdida-api/flash/price/query",
		client:   http.DefaultClient,
	}
}

func (i ItDiDaFetcher) Fetch(ctx context.Context, countryCode string, weight float64) ([]model.Logistics, error) {
	m := map[string]interface{}{
		"searchType":    2,
		"priceZoneType": 0,
		"weight":        weight,
		"packageType":   1,
		"pieceCount":    1,
		"wayTypeList":   []int{2, 1, 0},
		"containerType": 1,
		"unitModelList": []int{},
		"countryCode":   countryCode,
	}
	bb, err := json.Marshal(m)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req, err := http.NewRequest("POST", i.queryUrl, bytes.NewReader(bb))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := i.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	type QueryResp struct {
		Success    bool `json:"success"`
		StatusCode int  `json:"statusCode"`
		Data       []struct {
			ChannelName              string  `json:"channelName"`
			ChargeableWeight         float64 `json:"chargeableWeight"`
			SummaryConversion        float64 `json:"summaryConversion"`        // ??????
			FreightValue             float64 `json:"freightValue"`             // ??????
			FreightFuelCostsValue    float64 `json:"freightFuelCostsValue"`    // ?????????
			TotalSurchargeConversion float64 `json:"totalSurchargeConversion"` // ?????????
			Description              string  `json:"description"`
		} `json:"data"`
	}
	var queryResp QueryResp
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&queryResp)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !queryResp.Success {
		return nil, errors.New("error")
	}
	var res []model.Logistics
	for _, data := range queryResp.Data {
		res = append(res, model.Logistics{
			URL:    i.url,
			Method: data.ChannelName,
			Total:  data.SummaryConversion,
			Weight: data.ChargeableWeight,
			Fuel:   data.FreightFuelCostsValue,
			Fare:   data.FreightValue,
			Other:  data.TotalSurchargeConversion,
			Remark: data.Description,
		})
	}
	return res, nil
}
