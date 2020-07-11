package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cellargalaxy/smzdm-reptile/config"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
)

func sendWxPush(templateId string, tagId string, url string, data interface{}) error {
	jsonString, err := requestWxPush(templateId, tagId, url, data)
	if err == nil {
		return analysisWxPush(jsonString)
	}
	for i := 0; i < config.Retry; i++ {
		jsonString, err := requestWxPush(templateId, tagId, url, data)
		if err == nil {
			return analysisWxPush(jsonString)
		}
	}
	return fmt.Errorf("微信推送失败重试过多")
}

func analysisWxPush(jsonString string) error {
	var result struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Error("微信推送响应反序列化异常")
		return err
	}
	if result.Code != config.SuccessCode {
		logrus.Error("微信推送失败")
		return fmt.Errorf("微信推送失败: %s", result.Message)
	}
	return nil
}

func requestWxPush(templateId string, tagId string, url string, data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Error("序列化微信推送数据失败")
		return "", err
	}
	form := map[string]interface{}{"templateId": templateId, "tagId": tagId, "url": url, "data": string(bytes)}
	logrus.WithFields(logrus.Fields{"form": form}).Info("微信推送请求表单")

	request := gorequest.New()
	response, body, errs := request.Post(config.WxPushAddress).
		Type("multipart").
		Send(form).
		Timeout(config.Timeout).
		End()
	logrus.WithFields(logrus.Fields{"errs": errs}).Info("发送微信推送请求")
	if errs != nil && len(errs) > 0 {
		logrus.Error("发送微信推送请求异常")
		return "", errors.New("发送微信推送请求异常")
	}
	logrus.WithFields(logrus.Fields{"StatusCode": response.StatusCode, "body": body}).Info("发送微信推送请求")
	if response.StatusCode != 200 {
		logrus.Error("发送微信推送请求响应码异常")
		return "", errors.New("发送微信推送请求响应码异常")
	}
	return body, nil
}
