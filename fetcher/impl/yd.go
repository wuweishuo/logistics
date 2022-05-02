package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"logistics/config"
	"logistics/model"
	"net/http"
	"net/url"
)

type YDFetcher struct {
	ItDiDaFetcher
}

func NewYDFetcher() *YDFetcher {
	return &YDFetcher{
		NewItDiDaFetcher("YD", "http://1st.itdida.com/itdida-flash/website/landing",
			"http://1st.itdida.com/itdida-api/login", "http://1st.itdida.com/itdida-api/flash/price/query"),
	}
}

type ItDiDaFetcher struct {
	source   string
	url      string
	loginUrl string
	queryUrl string
	client   *http.Client
}

func NewItDiDaFetcher(source string, url string, loginUrl string, queryUrl string) ItDiDaFetcher {
	return ItDiDaFetcher{
		source:   source,
		url:      url,
		loginUrl: loginUrl,
		queryUrl: queryUrl,
		client:   http.DefaultClient,
	}
}

func (i ItDiDaFetcher) Fetch(ctx context.Context, config config.LoginConfig, countryCode string, weight float64) ([]model.Logistics, error) {
	//token, err := i.login(ctx, config)
	//if err != nil {
	//	return nil, err
	//}
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
		//"websocketToken": token,
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
			SummaryConversion        float64 `json:"summaryConversion"`        // 总价
			FreightValue             float64 `json:"freightValue"`             // 运费
			FreightFuelCostsValue    float64 `json:"freightFuelCostsValue"`    // 燃油费
			TotalSurchargeConversion float64 `json:"totalSurchargeConversion"` // 附加费
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
			Source: i.source,
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

func (i ItDiDaFetcher) login(ctx context.Context, loginConfig config.LoginConfig) (token string, err error) {
	resp, err := i.client.PostForm(i.loginUrl, url.Values{
		"username": []string{loginConfig.Username},
		"password": []string{loginConfig.Password},
	})
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	type LoginResp struct {
		Success    bool   `json:"success"`
		StatusCode int    `json:"statusCode"`
		Data       string `json:"data"`
	}
	var loginResp LoginResp
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&loginResp)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if !loginResp.Success {
		return "", errors.New(loginResp.Data)
	}
	return loginResp.Data, nil
}
