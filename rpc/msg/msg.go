package msg

import (
	"context"
	"fmt"
	"github.com/cellargalaxy/msg_gateway/sdk"
	"github.com/cellargalaxy/smzdm_reptile/cache"
	"time"
)

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
	return sdk.Client.SendTgMsg2ConfigChatId(ctx, text)
}
func SendWx(ctx context.Context, text string) (bool, error) {
	return sdk.Client.SendTemplateToCommonTag(ctx, text)
}

func SendWxTemplateToTag(ctx context.Context, templateId string, tagId int, url string, data map[string]interface{}) (bool, error) {
	return sdk.Client.SendWxTemplateToTag(ctx, templateId, tagId, url, data)
}
