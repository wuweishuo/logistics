package model

type Logistics struct {
	Source string  // 来源
	Method string  // 方式
	Weight float64 // 重量
	Total  float64 // 总价
	Price  float64 // 单价
	Fare   float64 // 运费
	Fuel   float64 // 燃油
	Other  float64 // 其他杂费
	Remark string  // 备注
}
