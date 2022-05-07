package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"logistics/fetcher"
	"logistics/model"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type IOMSFetcherConfig struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	CompanyNo string `yaml:"company_no"`
}

func (i *IOMSFetcherConfig) Parse(value *yaml.Node) error {
	return value.Decode(i)
}

type IOMSFetcherFactory struct{}

func (I IOMSFetcherFactory) ConstructFetcher(config fetcher.FetcherConfig) (fetcher.Fetcher, error) {
	fetcherConfig, ok := config.(*IOMSFetcherConfig)
	if !ok {
		return nil, errors.New("config not right")
	}
	return NewIOMSFetcher(*fetcherConfig), nil
}

func (I IOMSFetcherFactory) ConstructConfig() fetcher.FetcherConfig {
	return &IOMSFetcherConfig{}
}

type IOMSFetcher struct {
	config    IOMSFetcherConfig
	url       string
	companyNo string
	client    *http.Client
}

func NewIOMSFetcher(config IOMSFetcherConfig) *IOMSFetcher {
	return &IOMSFetcher{
		config:    config,
		url:       fmt.Sprintf("https://www.i-oms.cn/#/tmslogin?companyNo=%s", config.CompanyNo),
		companyNo: config.CompanyNo,
		client:    http.DefaultClient,
	}
}

type IOMSResp struct {
	ResultCode int    `json:"result_code"`
	Message    string `json:"message"`
	Body       string `json:"body"`
}

type IOMSQueryData struct {
	Datas []struct {
		TransTypeName   string  `json:"transTypeName"`
		UnitPrice       float64 `json:"unitPrice"`
		Weight          float64 `json:"weight"`
		TotalCharge     float64 `json:"totalCharge"`
		TransportCharge float64 `json:"transportCharge"`
		FuelCharge      float64 `json:"fuelCharge"`
		ExtraCharge     float64 `json:"extraCharge"`
	} `json:"datas"`
}

func (l IOMSFetcher) Fetch(ctx context.Context, countryCode string, weight float64) ([]model.Logistics, error) {
	token, err := l.getToken(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.client.PostForm("https://www.i-oms.cn/oms-web/findQuotes/1.0", url.Values{
		"head": []string{fmt.Sprintf("{\"appid\":\"LeonPC\",\"device_id\":\"Leon\",\"command\":\"findQuotes\",\"version\":\"1.0\",\"token\":\"%s\",\"sign\":\"\",\"encrypt_type\":0}", token)},
		"body": []string{fmt.Sprintf("{\"transCategory\":\"\",\"goodsType\":\"\",\"packageType\":\"WPX\",\"goodsTypeName\":\"\",\"dest\":\"%s\",\"payment\":\"\",\"partnerCompanyNo\":\"%s\",\"clientNo\":\"\",\"postCode\":\"\",\"businessTypes\":\"\",\"declareMethod\":\"其他\",\"weight\":\"%v\",\"long\":0,\"width\":0,\"height\":0,\"clientName\":\"\",\"vol\":\"0\",\"pageId\":1,\"pageIndex\":1,\"pageSize\":100,\"optType\":\"priceQuery_findQuotes\",\"payType\":\"PP\",\"destType\":0}", countryCode, l.companyNo, weight)},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var iOMSResp IOMSResp
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&iOMSResp)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if iOMSResp.ResultCode != 0 {
		return nil, errors.New(iOMSResp.Message)
	}
	var datas IOMSQueryData
	err = json.Unmarshal([]byte(iOMSResp.Body), &datas)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var res []model.Logistics
	for _, data := range datas.Datas {
		res = append(res, model.Logistics{
			URL:    l.url,
			Method: data.TransTypeName,
			Weight: data.Weight,
			Total:  data.TotalCharge,
			Price:  data.UnitPrice,
			Fare:   data.TransportCharge,
			Fuel:   data.FuelCharge,
			Other:  data.ExtraCharge,
		})
	}
	return res, nil
}

func (l IOMSFetcher) getToken(ctx context.Context) (string, error) {
	m := url.Values{
		"head": []string{"{\"appid\":\"LeonPC\",\"device_id\":\"Leon\",\"command\":\"tmsLogin\",\"version\":\"1.0\",\"token\":null,\"sign\":\"\",\"encrypt_type\":0}"},
		"body": []string{fmt.Sprintf("{\"userNo\":\"%s\",\"password\":\"%s\",\"companyNo\":\"%s\",\"domainName\":\"\"}", l.config.Username, l.config.Password, l.companyNo)},
	}
	resp, err := l.client.PostForm("https://www.i-oms.cn/user-center/tmsLogin/1.0", m)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var loginMsg IOMSResp
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&loginMsg)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if loginMsg.ResultCode != 0 {
		return "", errors.New(loginMsg.Message)
	}
	var tokenMsg map[string]interface{}
	err = json.Unmarshal([]byte(loginMsg.Body), &tokenMsg)
	if err != nil {
		return "", errors.WithStack(err)
	}
	token := tokenMsg["token"].(string)
	return token, nil
}
