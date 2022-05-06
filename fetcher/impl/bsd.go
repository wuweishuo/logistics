package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"logistics/fetcher"
	"logistics/model"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gopkg.in/yaml.v3"
)

type BSDFetcherConfig struct {
	Username                string `yaml:"username"`
	Password                string `yaml:"password"`
	Tariff                  string `yaml:"tariff"`
	GoodsRules              string `yaml:"goods_rules"`
	CancelAfterVerification string `yaml:"cancel_after_verification"`
	Deliverytype            string `yaml:"deliverytype"`
	Caddrfrom               string `yaml:"caddrfrom"`
}

func (b *BSDFetcherConfig) Parse(node *yaml.Node) error {
	return node.Decode(b)
}

type BSDFetcherFactory struct{}

func (B BSDFetcherFactory) ConstructFetcher(config fetcher.FetcherConfig) (fetcher.Fetcher, error) {
	fetcherConfig, ok := config.(*BSDFetcherConfig)
	if !ok {
		return nil, errors.New("config not right")
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return NewBSDFetcher(*fetcherConfig, &http.Client{
		Jar: jar,
	}), nil
}

func (B BSDFetcherFactory) ConstructConfig() fetcher.FetcherConfig {
	return &BSDFetcherConfig{}
}

type BSDFetcher struct {
	config BSDFetcherConfig
	client *http.Client
	ids    map[string]int
}

func NewBSDFetcher(config BSDFetcherConfig, client *http.Client) *BSDFetcher {
	return &BSDFetcher{
		config: config,
		client: client,
		ids: map[string]int{
			"1U": 335,
			"2U": 336,
			"3U": 337,
			"4U": 338,
			"U1": 331,
			"U2": 330,
			"U3": 332,
			"US": 221,
			"":   412,
			"!!": 303,
			"10": 376,
			"11": 377,
			"12": 378,
			"13": 379,
			"14": 380,
			"15": 381,
			"16": 382,
			"17": 383,
			"18": 384,
			"19": 385,
			"5U": 386,
			"6U": 387,
			"7U": 388,
			"8U": 389,
			"9U": 390,
			"AD": 23,
			"AE": 219,
			"AF": 20,
			"AG": 26,
			//"AI":25,
			"AI":      349,
			"AL":      21,
			"AM":      28,
			"AN":      159,
			"AO":      24,
			"AR":      27,
			"AS":      182,
			"AT":      31,
			"AU":      30,
			"AU2":     409,
			"AW":      29,
			"AX":      350,
			"AZ":      32,
			"AZEQ":    351,
			"AZYQ":    352,
			"BA":      44,
			"BB":      36,
			"BD":      304,
			"BE":      38,
			"BF":      49,
			"BG":      48,
			"BH":      34,
			"BI":      50,
			"BJ":      40,
			"BM":      41,
			"BN":      47,
			"BO":      43,
			"BR":      46,
			"BS":      33,
			"BT":      42,
			"BUS":     357,
			"BV":      248,
			"BW":      45,
			"BY":      37,
			"BZ":      39,
			"CA":      53,
			"CAMPION": 364,
			"CC":      235,
			"CD":      359,
			//"CD":62,
			"CF": 58,
			"CG": 63,
			"CH": 203,
			//"CI":240,
			"CI":   306,
			"CK":   162,
			"CL":   60,
			"CM":   52,
			"CN":   233,
			"CO":   61,
			"CR":   64,
			"CT":   411,
			"CU":   66,
			"CV":   54,
			"CX":   234,
			"CY":   67,
			"CZ":   68,
			"DE":   91,
			"DFI":  321,
			"DJ":   70,
			"DK":   69,
			"DM":   346,
			"DO":   72,
			"DZ":   22,
			"EC":   74,
			"EE":   79,
			"EG":   75,
			"ELD":  416,
			"ER":   78,
			"ES":   194,
			"ET":   80,
			"FI":   83,
			"FJ":   82,
			"FK":   237,
			"FM":   55,
			"FO":   81,
			"FR":   84,
			"FTW1": 433,
			"GA":   87,
			"GB":   220,
			"GBB":  370,
			"GC":   348,
			"GD":   96,
			"GE":   90,
			"GF":   85,
			"GG":   249,
			"GH":   92,
			"GI":   93,
			"GL":   95,
			"GM":   88,
			"GN":   100,
			"GP":   339,
			"GQ":   77,
			"GR":   94,
			"GT":   99,
			"GU":   98,
			"GW":   101,
			//"GW":363,
			"GY":   102,
			"GYR3": 429,
			"HG":   361,
			"HK":   418,
			"HM":   239,
			"HN":   104,
			"HR":   65,
			"HT":   103,
			"HU":   105,
			"HW":   238,
			"IC":   323,
			"ID":   108,
			"IE":   111,
			"IL":   112,
			"IN":   107,
			//"IN":307,
			"IND9": 435,
			"ION":  319,
			"IQ":   110,
			"IR":   109,
			//"IR":413,
			"IS":      106,
			"IT":      113,
			"ITA":     415,
			"ITB":     414,
			"JE":      241,
			"JM":      114,
			"JO":      116,
			"JP":      115,
			"KE":      118,
			"KG":      123,
			"KH":      51,
			"KI":      119,
			"KM":      236,
			"KN":      196,
			"KP":      120,
			"KR":      121,
			"KV":      315,
			"KW":      122,
			"KY":      57,
			"KZ":      117,
			"LA":      124,
			"LAKE":    367,
			"LAS1":    430,
			"LAX9":    420,
			"LB":      126,
			"LC":      197,
			"LCC":     403,
			"LGB3":    422,
			"LGB4":    423,
			"LGB6":    424,
			"LGB8":    421,
			"LI":      130,
			"LIVIGNO": 366,
			"LK":      195,
			"LL":      393,
			"LM":      371,
			"LR":      128,
			"LS":      127,
			"LT":      131,
			"LU":      132,
			"LV":      125,
			"LY":      129,
			"MA":      152,
			"MAD":     369,
			"MC":      148,
			"MD":      147,
			"MDW2":    434,
			"ME":      150,
			"MEL":     395,
			"MG":      135,
			"MH":      142,
			"MK":      134,
			"ML":      139,
			"MLL":     408,
			"MM":      154,
			"MN":      149,
			"MO":      417,
			"MP":      141,
			"MQ":      143,
			"MR":      144,
			"MS":      151,
			"MSQ":     354,
			"MT":      140,
			"MU":      145,
			"MV":      138,
			"MW":      136,
			"MX":      146,
			"MY":      137,
			"MY1":     372,
			"MY2":     373,
			"MZ":      153,
			"NA":      155,
			"NB":      355,
			"NC":      160,
			"NE":      164,
			//"NF":166,
			//"NF":398,
			"NG":   165,
			"NGY":  394,
			"NI":   163,
			"NL":   158,
			"NO":   167,
			"NP":   157,
			"NR":   156,
			"NU":   243,
			"NZ":   161,
			"OM":   168,
			"ONT2": 425,
			"ONT6": 426,
			"ONT8": 419,
			"ONT9": 427,
			"OOL":  362,
			"PA":   170,
			"PE":   173,
			"PER":  400,
			"PF":   345,
			"PG":   171,
			"PH":   174,
			"PI":   308,
			"PK":   169,
			"PL":   175,
			"PR":   177,
			"PS":   250,
			//"PS":353,
			"PT": 176,
			//"PW":56,
			"PW":   399,
			"PY":   172,
			"QA":   178,
			"RE":   309,
			"RO":   179,
			"RS":   186,
			"RU":   180,
			"RW":   181,
			"SA":   184,
			"SB":   192,
			"SC":   187,
			"SCK4": 431,
			"SD":   199,
			"SE":   202,
			"SG":   189,
			"SI":   191,
			"SK":   190,
			"SL":   188,
			"SM":   244,
			"SMF3": 432,
			"SN":   185,
			"SNA4": 428,
			"SO":   245,
			"SP":   318,
			"SR":   200,
			"ST":   183,
			"STD":  406,
			"SV":   76,
			"SVG":  231,
			"SX":   360,
			"SY":   204,
			"SYD":  410,
			"SZ":   201,
			"TC":   215,
			"TD":   59,
			"TG":   209,
			"TH":   208,
			"TJ":   206,
			"TL":   73,
			"TM":   214,
			"TN":   212,
			"TO":   210,
			"TR":   213,
			"TT":   211,
			"TV":   216,
			"TW":   343,
			"TZ":   207,
			"UA":   218,
			"UG":   217,
			"UKL":  407,
			"USA":  392,
			"USB":  391,
			"UY":   222,
			"UZ":   310,
			"VC":   198,
			"VE":   224,
			"VG":   246,
			"VI":   232,
			"VN":   225,
			"VT":   358,
			"VTT":  368,
			"VU":   223,
			"WF":   347,
			"WL":   247,
			"WS":   341,
			"XB":   356,
			"XC":   333,
			//"XC":365,
			//"XE":89,
			"XE":  405,
			"XEE": 404,
			"XL":  242,
			"XM":  334,
			"XN":  397,
			"XS":  344,
			"XY":  401,
			"XYY": 402,
			"YE":  227,
			"YT":  375,
			"YUU": 396,
			"ZA":  193,
			"ZM":  228,
			"ZW":  229,
		},
	}
}

func (b BSDFetcher) Fetch(ctx context.Context, countryCode string, weight float64) ([]model.Logistics, error) {
	err := b.login(ctx)
	if err != nil {
		return nil, err
	}
	res, err := b.getFee(ctx, countryCode, weight)
	if err != nil {
		return nil, err

	}
	for i := range res {
		otherFee, err := b.getOtherFee(ctx, countryCode, weight, res[i].Method)
		if err != nil {
			return nil, err
		}
		res[i].Other = otherFee.InexactFloat64()
		res[i].Total = otherFee.Add(decimal.NewFromFloat(res[i].Total)).InexactFloat64()
	}
	return res, nil
}

func (b BSDFetcher) login(ctx context.Context) error {
	resp, err := b.client.PostForm("http://mis.bsdexp.com/login", url.Values{
		"Username": []string{b.config.Username},
		"Password": []string{b.config.Password},
	})
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.Request.Method != http.MethodGet {
		return errors.New("登录失败")
	}
	return nil
}

func (b BSDFetcher) getFee(ctx context.Context, countryCode string, weight float64) ([]model.Logistics, error) {
	resp, err := b.client.PostForm("http://mis.bsdexp.com/FeeSearch/GetFee", url.Values{
		"sTargetCountryID": []string{fmt.Sprintf("%d", b.ids[countryCode])},
		"sStage":           []string{fmt.Sprintf("%v", weight)},
		"sPackageType":     []string{"包裹"},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	type QueryResp struct {
		Total  int    `json:"total"`
		Errors string `json:"errors"`
		Rows   []struct {
			Algo               string `json:"Algo"`
			MinSaleAmountTotal string `json:"MinSaleAmountTotal"`
			CHCnName           string `json:"CHCnName"`
			BaseFee            string `json:"BaseFee"`
			CalWeight          string `json:"CalWeight"`
			FuelFee            string `json:"FuelFee"`
		} `json:"rows"`
	}
	var queryResp QueryResp
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&queryResp)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if queryResp.Errors != "" {
		return nil, errors.New(queryResp.Errors)
	}
	var res []model.Logistics
	for _, row := range queryResp.Rows {
		fare, err := strconv.ParseFloat(row.BaseFee, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		calcWeight, err := strconv.ParseFloat(row.CalWeight, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		fuel, err := strconv.ParseFloat(row.FuelFee, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		total, err := strconv.ParseFloat(row.MinSaleAmountTotal, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		res = append(res, model.Logistics{
			URL:    "http://mis.bsdexp.com/Account",
			Method: row.CHCnName,
			Total:  total,
			Weight: calcWeight,
			Fuel:   fuel,
			Fare:   fare,
			Remark: row.Algo,
		})
	}
	return res, nil
}

func (b BSDFetcher) getOtherFee(ctx context.Context, countryCode string, weight float64, cemskind string) (decimal.Decimal, error) {
	resp, err := b.client.PostForm("http://mis.bsdexp.com/FeeSearch/GetOtherFee", url.Values{
		"cemskind":                []string{cemskind},
		"CountryId":               []string{fmt.Sprintf("%d", b.ids[countryCode])},
		"fweight":                 []string{fmt.Sprintf("%v", weight)},
		"tariff":                  []string{b.config.Tariff},
		"goodsRules":              []string{b.config.GoodsRules},
		"cancelAfterVerification": []string{b.config.CancelAfterVerification},
		"deliverytype":            []string{b.config.Deliverytype},
		"caddrfrom":               []string{b.config.Caddrfrom},
	})
	if err != nil {
		return decimal.Zero, errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	type QueryResp struct {
		Fee float64 `json:"fee"`
	}
	var queryResp []QueryResp
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&queryResp)
	if err != nil {
		return decimal.Zero, errors.WithStack(err)
	}
	var otherFee decimal.Decimal
	for _, v := range queryResp {
		otherFee = otherFee.Add(decimal.NewFromFloat(v.Fee))
	}
	return otherFee, nil
}
