package test

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/smzdm_reptile/config"
	"github.com/cellargalaxy/smzdm_reptile/spider"
	"testing"
)

func TestListGoods(test *testing.T) {
	config.Init()

	ctx := util.GenCtx()
	object, err := spider.ListGoods(ctx, "鼠标")
	if err != nil {
		panic(err)
	}
	test.Logf("object, %+v", len(object))
	for i := range object {
		test.Logf("object, %+v", util.ToJsonString(object[i]))
	}
}
