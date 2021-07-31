package main

import (
	"context"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/smzdm-reptile/controller"
	"github.com/cellargalaxy/smzdm-reptile/service"
)

func main() {
	util.InitLog("smzdm.log")
	ctx := context.Background()
	ctx = util.SetLogId(ctx)
	go service.StartSearchService(ctx)
	controller.StartWebService()
}
