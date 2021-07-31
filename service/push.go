package service

import (
	"context"
	"github.com/cellargalaxy/msg-gateway/sdk"
	"github.com/cellargalaxy/smzdm-reptile/config"
)

var msgClient *sdk.MsgClient

func init() {
	var err error
	msgClient, err = sdk.NewMsgClient(config.Timeout, config.Sleep, config.Retry, config.WxPushAddress, config.WxToken)
	if err != nil {
		panic(err)
	}
}

func SendWxPush(ctx context.Context, templateId string, tagId int, url string, data map[string]interface{}) error {
	_, err := msgClient.SendWxTemplateToTag(ctx, templateId, tagId, url, data)
	return err
}
