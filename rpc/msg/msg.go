package msg

import (
	"context"
	"fmt"
	"github.com/cellargalaxy/msg_gateway/sdk"
	"github.com/cellargalaxy/smzdm_reptile/cache"
	"time"
)

var client *sdk.MsgClient

func init() {
	var err error
	client, err = sdk.NewDefaultMsgClient()
	if err != nil {
		panic(err)
	}
}

func SendErr(ctx context.Context, name string, err error) {
	if cache.RateLimit("SendErr-"+name, time.Hour) {
		return
	}
	Send(ctx, fmt.Sprintf("%+v\nerr: %+v", name, err))
}

func Send(ctx context.Context, text string) {
	go SendTg(ctx, text)
	go SendWx(ctx, text)
}

func SendTg(ctx context.Context, text string) (bool, error) {
	return client.SendTgMsg2ConfigChatId(ctx, text)
}
func SendWx(ctx context.Context, text string) (bool, error) {
	return client.SendTemplateToCommonTag(ctx, text)
}

func SendWxTemplateToTag(ctx context.Context, templateId string, tagId int, url string, data map[string]interface{}) (bool, error) {
	return client.SendWxTemplateToTag(ctx, templateId, tagId, url, data)
}
