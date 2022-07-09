package main

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/smzdm_reptile/model"
	"github.com/cellargalaxy/smzdm_reptile/service"
)

func init() {
	util.Init(model.DefaultServerName)
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
