package msg

import (
	"context"
	"fmt"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/msg_gateway/sdk"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

var localCache *util.LocalCache

func init() {
	ctx := util.GenCtx()
	var err error
	localCache, err = util.NewDefaultLocalCache(ctx)
	if err != nil {
		panic(err)
	}
}

func Sends(ctx context.Context, name string, duration time.Duration, texts ...string) {
	text := name
	if len(texts) > 0 {
		text = fmt.Sprintf("%+v\n%+v", text, strings.Join(texts, "\n"))
	}
	Send(ctx, name, text, duration)
}

func SendErr(ctx context.Context, name string, err interface{}) {
	if err == nil {
		return
	}
	text := fmt.Sprintf("%+v\nerr: %+v", name, err)
	Send(ctx, name, text, time.Minute*10)
}

func Send(ctx context.Context, key, text string, duration time.Duration) {
	ctx = util.CopyCtx(ctx)
	if key != "" && !localCache.TryLock(ctx, key, duration) {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"text": text}).Warn("发送消息，限频")
		return
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"text": text}).Info("发送消息")
	go sendTg(ctx, text)
	go sendWx(ctx, text)
}

func sendTg(ctx context.Context, text string) error {
	return sdk.Client.SendTgMsg2ConfigChatId(ctx, text)
}
func sendWx(ctx context.Context, text string) error {
	return sdk.Client.SendTemplateToCommonTag(ctx, text)
}

func SendWxTemplateToTag(ctx context.Context, templateId string, tagId int, url string, data map[string]interface{}) error {
	return sdk.Client.SendWxTemplateToTag(ctx, templateId, tagId, url, data)
}
