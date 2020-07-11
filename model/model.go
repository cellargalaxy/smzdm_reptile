package model

import "time"

type SearchCondition struct {
	SearchId        string  `json:"searchId"`
	SearchKey       string  `json:"searchKey"`
	TitleContain    string  `json:"titleContain"`
	TitleExclude    string  `json:"titleExclude"`
	MinPrice        float32 `json:"minPrice"`
	MaxPrice        float32 `json:"maxPrice"`
	MinZhi          int     `json:"minZhi"`
	MaxBuzhi        int     `json:"maxBuzhi"`
	MerchantContain string  `json:"merchantContain"`
	MerchantExclude string  `json:"merchantExclude"`
	WxTemplateId    string  `json:"wxTemplateId"`
	WxTagId         string  `json:"wxTagId"`
}

type Goods struct {
	Title    string
	Price    float32
	Zhi      int
	Buzhi    int
	Merchant string
	Url      string
	Date     time.Time
}
