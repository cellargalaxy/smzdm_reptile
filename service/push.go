package service

import (
	"github.com/cellargalaxy/smzdm-reptile/config"
	"github.com/cellargalaxy/wx-gateway/sdk"
)

var wxClient *sdk.WxClient

func init() {
	var err error
	wxClient, err = sdk.NewWxClient(config.Timeout, config.Sleep, config.Retry, config.WxPushAddress, config.WxToken)
	if err != nil {
		panic(err)
	}
}

func SendWxPush(templateId string, tagId int, url string, data map[string]interface{}) error {
	_, err := wxClient.SendTemplateToTag(templateId, tagId, url, data)
	return err
}
