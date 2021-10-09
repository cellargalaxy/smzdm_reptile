package model

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	LogLevel   logrus.Level      `json:"log_level" yaml:"log_level"`
	Retry      int               `json:"retry" yaml:"retry"`
	Timeout    time.Duration     `json:"timeout" yaml:"timeout"`
	Sleep      time.Duration     `json:"sleep" yaml:"sleep"`
	MaxPage    int               `json:"max_page" yaml:"max_page"`
	Conditions []SearchCondition `json:"conditions" yaml:"conditions"`
}

func (this Config) String() string {
	return util.ToJsonString(this)
}

type SearchCondition struct {
	SearchId        string  `json:"search_id" yaml:"search_id"`
	SearchKey       string  `json:"search_key" yaml:"search_key"`
	TitleContain    string  `json:"title_contain" yaml:"title_contain"`
	TitleExclude    string  `json:"title_exclude" yaml:"title_exclude"`
	MinPrice        float64 `json:"min_price" yaml:"min_price"`
	MaxPrice        float64 `json:"max_price" yaml:"max_price"`
	MinZhi          int     `json:"min_zhi" yaml:"min_zhi"`
	MaxBuzhi        int     `json:"max_buzhi" yaml:"max_buzhi"`
	MerchantContain string  `json:"merchant_contain" yaml:"merchant_contain"`
	MerchantExclude string  `json:"merchant_exclude" yaml:"merchant_exclude"`
	WxTemplateId    string  `json:"wx_template_id" yaml:"wx_template_id"`
	WxTagId         int     `json:"wx_tag_id" yaml:"wx_tag_id"`
}

func (this SearchCondition) String() string {
	return util.ToJsonString(this)
}
