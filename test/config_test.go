package test

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/smzdm_reptile/model"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestConfig(test *testing.T) {
	var config model.Config
	config.LogLevel = logrus.InfoLevel
	config.Retry = 3
	config.Timeout = time.Second * 3
	config.Sleep = time.Second * 3
	config.MaxPage = 10
	config.Conditions = append(config.Conditions, model.SearchCondition{
		SearchId:        "SearchId",
		SearchKey:       "SearchKey",
		TitleContain:    "TitleContain",
		TitleExclude:    "TitleExclude",
		MinPrice:        123,
		MaxPrice:        234,
		MinZhi:          345,
		MaxBuzhi:        456,
		MerchantContain: "MerchantContain",
		MerchantExclude: "MerchantExclude",
		WxTemplateId:    "WxTemplateId",
		WxTagId:         567,
	})
	test.Logf("%+v", util.ToYamlString(config))
}
