package config

import (
	"context"
	"github.com/cellargalaxy/go_common/util"
	sc_model "github.com/cellargalaxy/server_center/model"
	"github.com/cellargalaxy/server_center/sdk"
	"github.com/cellargalaxy/smzdm_reptile/model"
	"github.com/sirupsen/logrus"
	"time"
)

var Config = model.Config{}

func init() {
	ctx := util.CreateLogCtx()
	client, err := sdk.NewDefaultServerCenterClient(ctx, &ServerCenterHandler{})
	if err != nil {
		panic(err)
	}
	client.StartConfWithInitConf(ctx)
}

func checkAndResetConfig(ctx context.Context, config model.Config) (model.Config, error) {
	if config.Timeout < 0 {
		config.Timeout = 3 * time.Second
	}
	if config.Sleep < 0 {
		config.Sleep = 3 * time.Second
	}
	return config, nil
}

type ServerCenterHandler struct {
}

func (this *ServerCenterHandler) GetServerName(ctx context.Context) string {
	return sdk.GetEnvServerName(ctx, model.DefaultServerName)
}
func (this *ServerCenterHandler) GetAddress(ctx context.Context) string {
	return sdk.GetEnvServerCenterAddress(ctx)
}
func (this *ServerCenterHandler) GetSecret(ctx context.Context) string {
	return sdk.GetEnvServerCenterSecret(ctx)
}
func (this *ServerCenterHandler) GetInterval(ctx context.Context) time.Duration {
	return 5 * time.Minute
}
func (this *ServerCenterHandler) ParseConf(ctx context.Context, object sc_model.ServerConfModel) error {
	var config model.Config
	err := util.UnmarshalYamlString(object.ConfText, &config)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("反序列化配置异常")
		return err
	}
	config, err = checkAndResetConfig(ctx, config)
	if err != nil {
		return err
	}
	Config = config
	logrus.WithContext(ctx).WithFields(logrus.Fields{"Config": Config}).Info("加载配置")
	return nil
}
func (this *ServerCenterHandler) GetDefaultConf(ctx context.Context) string {
	var config model.Config
	config, _ = checkAndResetConfig(ctx, config)
	return util.ToYamlString(config)
}
func (this *ServerCenterHandler) GetLocalFilePath(ctx context.Context) string {
	return ""
}
