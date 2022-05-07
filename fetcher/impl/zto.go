package impl

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gopkg.in/yaml.v3"
	"logistics/fetcher"
	"logistics/model"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

type ZTOFetcherConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	PubKey   string `yaml:"pub_key"`
}

func (z *ZTOFetcherConfig) Parse(node *yaml.Node) error {
	return node.Decode(z)
}

type ZTOFetcherFactory struct{}

func (z ZTOFetcherFactory) ConstructFetcher(config fetcher.FetcherConfig) (fetcher.Fetcher, error) {
	fetcherConfig, ok := config.(*ZTOFetcherConfig)
	if !ok {
		return nil, errors.New("config not right")
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return ZTOFetcher{
		config: *fetcherConfig,
		client: &http.Client{
			Jar: jar,
		},
	}, nil
}

func (z ZTOFetcherFactory) ConstructConfig() fetcher.FetcherConfig {
	return &ZTOFetcherConfig{}
}

type ZTOFetcher struct {
	config ZTOFetcherConfig
	client *http.Client
}

func (z ZTOFetcher) Fetch(ctx context.Context, countryCode string, weight float64) ([]model.Logistics, error) {
	err := z.login(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := z.client.PostForm("https://ioes-client.ztoglobal.com/exp-client/admin/base/quoteItem/queryQuoteItem", url.Values{
		"customerChargeWeight": []string{fmt.Sprintf("%v", weight)},
		"packTypeCode":         []string{"WPX"},
		"destCode":             []string{countryCode},
		"pageSize":             []string{"100"},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var queryResp struct {
		Total int `yaml:"total"`
		List  []struct {
			ProductChannelName   string  `json:"productChannelName"`
			CustomerChargeWeight float64 `json:"customerChargeWeight"`
			Shipping             string  `json:"shipping"`
			FeeTotal             string  `json:"feeTotal"`
			WarehouseFee         string  `json:"warehouseFee"`
			PeakSeasonFee        string  `json:"peakSeasonFee"`
		} `yaml:"list"`
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&queryResp)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var res []model.Logistics
	for _, data := range queryResp.List {
		var total, fare, other float64
		total, err = strconv.ParseFloat(strings.TrimLeft(data.FeeTotal, "CNY"), 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		fare, err = strconv.ParseFloat(data.Shipping, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		warehouseFee, peakSeasonFee := decimal.Zero, decimal.Zero
		if data.WarehouseFee != "-" {
			warehouseFee, err = decimal.NewFromString(data.WarehouseFee)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}
		if data.PeakSeasonFee != "-" {
			peakSeasonFee, err = decimal.NewFromString(data.PeakSeasonFee)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}
		other = warehouseFee.Add(peakSeasonFee).InexactFloat64()
		res = append(res, model.Logistics{
			URL:    "https://ioes-client.ztoglobal.com/exp-client/admin/main",
			Method: data.ProductChannelName,
			Weight: data.CustomerChargeWeight,
			Total:  total,
			Fare:   fare,
			Other:  other,
		})
	}
	return res, nil
}

func (z ZTOFetcher) login(ctx context.Context) error {
	redirectUrl, err := z.getCode(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	resp, err := z.client.Get(redirectUrl)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if !strings.Contains(resp.Request.URL.Path, "/exp-client/admin/main") {
		return errors.New("登录失败")
	}
	return nil
}

func (z ZTOFetcher) getCode(ctx context.Context) (string, error) {
	authorization, err := z.getAuthorization(ctx)
	if err != nil {
		return "", err
	}
	m := map[string]interface{}{
		"response_type":      "code",
		"authorization_type": "passwd",
		"redirect_url":       "https://ioes-client.ztoglobal.com/exp-client/admin/iamLogin/callback",
		"app_id":             "ztudIbxQ7UUH6--KmPCq8CYw",
		"view":               "web",
		"authorization":      authorization,
	}
	bs, err := json.Marshal(m)
	if err != nil {
		return "", errors.WithStack(err)
	}
	req, err := http.NewRequest(http.MethodPost, "https://iam-int.zto.com/api/oauth2/authorize?app_id=ztudIbxQ7UUH6--KmPCq8CYw", bytes.NewReader(bs))
	if err != nil {
		return "", errors.WithStack(err)
	}
	resp, err := z.client.Do(req)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var authResp struct {
		ErrorCode    string `json:"error_code"`
		ErrorMessage string `json:"error_message"`
		Data         struct {
			RedirectUrl string `json:"redirect_url"`
		} `json:"data"`
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&authResp)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if authResp.ErrorCode != "" {
		return "", errors.New(authResp.ErrorMessage)
	}
	return authResp.Data.RedirectUrl, nil
}

func (z ZTOFetcher) getAuthorization(ctx context.Context) (string, error) {
	var pubKey = []byte(z.config.PubKey)
	block, _ := pem.Decode(pubKey) //将密钥解析成公钥实例
	if block == nil || block.Type != "PUBLIC KEY" {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	if err != nil {
		return "", errors.WithStack(err)
	}
	pub := pubInterface.(*rsa.PublicKey)
	data, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(fmt.Sprintf("%s %s", z.config.Username, z.config.Password)))
	if err != nil {
		return "", errors.WithStack(err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
