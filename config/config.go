package config

import (
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

const (
	SuccessCode = 1
	FailCode    = 2

	ConfigFilePath          = "resources/config.ini"
	SearchConditionFilePath = "resources/searchCondition.json"

	defaultRetry         = 2
	defaultMaxPage       = 10
	defaultTimeoutSecond = 5
	defaultTimeout       = defaultTimeoutSecond * time.Second
	defaultSleepSecond   = 2
	defaultSleep         = defaultSleepSecond * time.Second
	defaultListenAddress = ":8080"
	defaultWxPushAddress = ""
)

var Retry = defaultRetry
var MaxPage = defaultMaxPage
var Timeout = defaultTimeout
var Sleep = defaultSleep
var ListenAddress = defaultListenAddress
var WxPushAddress = defaultWxPushAddress

func init() {
	logrus.Info("加载配置开始")

	configFile, err := ini.Load(ConfigFilePath)
	if err == nil {
		Retry = configFile.Section("").Key("retry").MustInt(defaultRetry)
		MaxPage = configFile.Section("").Key("maxPage").MustInt(defaultMaxPage)
		Timeout = time.Duration(configFile.Section("").Key("timeout").MustInt(defaultTimeoutSecond)) * time.Second
		Sleep = time.Duration(configFile.Section("").Key("sleep").MustInt(defaultSleepSecond)) * time.Second
		ListenAddress = configFile.Section("").Key("listenAddress").MustString(defaultListenAddress)
		WxPushAddress = configFile.Section("").Key("wxPushAddress").MustString(defaultWxPushAddress)
	} else {
		logrus.WithFields(logrus.Fields{"err": err}).Error("加载配置文件失败")
	}

	retryString := os.Getenv("RETRY")
	logrus.WithFields(logrus.Fields{"retryString": retryString}).Info("环境变量读取配置Retry")
	retry, err := strconv.Atoi(retryString)
	if err == nil {
		Retry = retry
	}

	maxPageString := os.Getenv("MAX_PAGE")
	logrus.WithFields(logrus.Fields{"maxPageString": maxPageString}).Info("环境变量读取配置MaxPage")
	maxPage, err := strconv.Atoi(maxPageString)
	if err == nil && maxPage > 0 {
		MaxPage = maxPage
	}

	timeoutString := os.Getenv("TIMEOUT")
	logrus.WithFields(logrus.Fields{"timeoutString": timeoutString}).Info("环境变量读取配置Timeout")
	timeout, err := strconv.Atoi(timeoutString)
	if err == nil && timeout > 0 {
		Timeout = time.Duration(timeout) * time.Second
	}

	sleepString := os.Getenv("SLEEP")
	logrus.WithFields(logrus.Fields{"sleepString": sleepString}).Info("环境变量读取配置Sleep")
	sleep, err := strconv.Atoi(sleepString)
	if err == nil && sleep > 0 {
		Sleep = time.Duration(sleep) * time.Second
	}

	listenAddress := os.Getenv("LISTEN_ADDRESS")
	logrus.WithFields(logrus.Fields{"listenAddress": listenAddress}).Info("环境变量读取配置ListenAddress")
	if listenAddress != "" {
		ListenAddress = listenAddress
	}

	wxPushAddress := os.Getenv("WX_PUSH_ADDRESS")
	logrus.WithFields(logrus.Fields{"wxPushAddress": wxPushAddress}).Info("环境变量读取配置WxPushAddress")
	if wxPushAddress != "" {
		WxPushAddress = wxPushAddress
	}

	logrus.WithFields(logrus.Fields{"Retry": Retry}).Info("配置Retry")
	logrus.WithFields(logrus.Fields{"MaxPage": MaxPage}).Info("配置MaxPage")
	logrus.WithFields(logrus.Fields{"Timeout": Timeout}).Info("配置Timeout")
	logrus.WithFields(logrus.Fields{"Sleep": Sleep}).Info("配置Sleep")
	logrus.WithFields(logrus.Fields{"ListenAddress": ListenAddress}).Info("配置ListenAddress")
	logrus.WithFields(logrus.Fields{"WxPushAddress": WxPushAddress}).Info("配置WxPushAddress")
}
