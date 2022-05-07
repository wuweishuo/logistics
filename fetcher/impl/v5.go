package impl

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"logistics/fetcher"
	"logistics/model"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type V5FetcherConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Domain   string `yaml:"domain"`
}

func (v *V5FetcherConfig) Parse(node *yaml.Node) error {
	return node.Decode(v)
}

type V5FetcherFactory struct{}

func (v V5FetcherFactory) ConstructFetcher(config fetcher.FetcherConfig) (fetcher.Fetcher, error) {
	fetcherConfig, ok := config.(*V5FetcherConfig)
	if !ok {
		return nil, errors.New("config not right")
	}
	return V5Fetcher{
		config: *fetcherConfig,
		client: http.DefaultClient,
	}, nil
}

func (v V5FetcherFactory) ConstructConfig() fetcher.FetcherConfig {
	return &V5FetcherConfig{}
}

type V5Fetcher struct {
	config V5FetcherConfig
	client *http.Client
}

func (v V5Fetcher) Fetch(ctx context.Context, countryCode string, weight float64) ([]model.Logistics, error) {
	token, err := v.login(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/tms-saas-oms//bms/quoteObj/list?destNo=%s&weig=%v&packing=WPX", v.config.Domain, countryCode, weight), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	base64Token, err := v.getToken(token)
	if err != nil {
		return nil, err
	}
	req.Header.Set("locale", "zhCN")
	req.Header.Set("token", base64Token)
	req.Header.Set("version", "1.0")
	resp, err := v.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var queryResp struct {
		Message    string `json:"message"`
		ResultCode int    `json:"result_code"`
		Body       struct {
			List []struct {
				Name                string  `json:"name"`
				StandardCharge      float64 `json:"standardCharge"`      // 运费
				StandardFeeSum      float64 `json:"standardFeeSum"`      // 总费用
				OtherFeeUnitSum     float64 `json:"otherFeeUnitSum"`     // 单价
				OtherStandardFeeSum float64 `json:"otherStandardFeeSum"` // 其他费用
				Remark              string  `json:"remark"`
				Weig                float64 `json:"weig"`
			} `yaml:"list"`
			PageNum  int `yaml:"pageNum"`
			PageSize int `yaml:"pageSize"`
			Pages    int `yaml:"pages"`
			Size     int `yaml:"size"`
		} `yaml:"body"`
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&queryResp)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if queryResp.ResultCode != 0 {
		return nil, errors.New(queryResp.Message)
	}
	var res []model.Logistics
	for _, data := range queryResp.Body.List {
		res = append(res, model.Logistics{
			URL:    v.config.Domain + "/viplogin",
			Method: data.Name,
			Weight: data.Weig,
			Total:  data.StandardFeeSum,
			Price:  data.OtherFeeUnitSum,
			Fare:   data.StandardCharge,
			Other:  data.OtherStandardFeeSum,
			Remark: data.Remark,
		})
	}
	return res, nil
}

func (v V5Fetcher) getRandomStr() string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	builder := strings.Builder{}
	for i := 0; i < 8; i++ {
		builder.WriteByte(chars[rand.Intn(len(chars))])
	}
	return builder.String()
}

func (v V5Fetcher) getToken(token string) (string, error) {
	m := map[string]interface{}{
		"timestamp": time.Now().Unix() * 1000,
		"nonce":     v.getRandomStr(),
		"token":     token,
	}
	str, err := json.Marshal(m)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return v.encodeMixChar(string(str)), nil
}

func (v V5Fetcher) encodeMixChar(str string) string {
	base64Token := base64.StdEncoding.EncodeToString([]byte(str))
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(base64Token, "a", "-"),
		"c", "#"),
		"x", "^"),
		"M", "$",
	)
}

func (v V5Fetcher) login(ctx context.Context) (string, error) {
	resp, err := v.client.PostForm(v.config.Domain+"/tms-saas-oms/user/login", url.Values{
		"userNo":   []string{v.config.Username},
		"password": []string{v.encodeMixChar(v.config.Password)},
	})
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var loginResp struct {
		Message    string `json:"message"`
		ResultCode int    `json:"result_code"`
		Body       struct {
			Token string `json:"token"`
		}
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&loginResp)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if loginResp.ResultCode != 0 {
		return "", errors.New(loginResp.Message)
	}
	return loginResp.Body.Token, nil
}
