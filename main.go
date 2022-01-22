package main

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/smzdm_reptile/config"
	"github.com/cellargalaxy/smzdm_reptile/service"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(config.Config.LogLevel)
	util.InitDefaultLog(util.GetServerNameWithPanic())
}

/**
export server_name=smzdm_reptile
export server_center_address=http://127.0.0.1:7557
export server_center_secret=secret_secret

server_name=smzdm_reptile;server_center_address=http://127.0.0.1:7557;server_center_secret=secret_secret
*/
func main() {
	service.StartSearchService()
}
