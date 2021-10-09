package main

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/smzdm_reptile/config"
	"github.com/cellargalaxy/smzdm_reptile/service"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(config.Config.LogLevel)
	util.InitLog(util.GetServerNameWithPanic())
}

func main() {
	service.StartSearchService()
}
