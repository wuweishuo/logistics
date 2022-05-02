package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"logistics/config"
	"logistics/model"
	"net/http"
	"net/url"
)

// IOMSFetcher 联递国际物流，https://www.i-oms.cn/#/tmslogin?companyNo=lde
type LDEFetcher struct {
	IOMSFetcher
}

type IOMSFetcher struct {
	source    string
	url       string
	companyNo string
	client    *http.Client
}

func NewLDEFetcher() *LDEFetcher {
	return &LDEFetcher{
		NewIOMSFetcher("LDE"),
	}
}

func NewIOMSFetcher(companyNo string) IOMSFetcher {
	return IOMSFetcher{
		source:    companyNo,
		url:       fmt.Sprintf("https://www.i-oms.cn/#/tmslogin?companyNo=%s", companyNo),
		companyNo: companyNo,
		client:    http.DefaultClient,
	}
}

type LianDiResp struct {
	ResultCode int    `json:"result_code"`
	Message    string `json:"message"`
	Body       string `json:"body"`
}

type LianDiQueryData struct {
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

func (l IOMSFetcher) Fetch(ctx context.Context, config config.LoginConfig, countryCode string, weight float64) ([]model.Logistics, error) {
	token, err := l.getToken(ctx, config)
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
	var lianDiResp LianDiResp
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&lianDiResp)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if lianDiResp.ResultCode != 0 {
		return nil, errors.New(lianDiResp.Message)
	}
	var datas LianDiQueryData
	err = json.Unmarshal([]byte(lianDiResp.Body), &datas)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var res []model.Logistics
	for _, data := range datas.Datas {
		res = append(res, model.Logistics{
			Source: l.source,
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

func (l IOMSFetcher) getToken(ctx context.Context, config config.LoginConfig) (string, error) {
	m := url.Values{
		"head": []string{"{\"appid\":\"LeonPC\",\"device_id\":\"Leon\",\"command\":\"tmsLogin\",\"version\":\"1.0\",\"token\":null,\"sign\":\"\",\"encrypt_type\":0}"},
		"body": []string{fmt.Sprintf("{\"userNo\":\"%s\",\"password\":\"%s\",\"companyNo\":\"%s\",\"domainName\":\"\"}", config.Username, config.Password, l.companyNo)},
	}
	resp, err := l.client.PostForm("https://www.i-oms.cn/user-center/tmsLogin/1.0", m)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var loginMsg LianDiResp
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
