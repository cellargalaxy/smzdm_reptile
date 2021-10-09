package cache

import (
	"context"
	"fmt"
	"github.com/cellargalaxy/smzdm_reptile/model"
	"time"
)

func YetSendGoods(ctx context.Context, searchCondition model.SearchCondition, goods model.Goods) bool {
	key := fmt.Sprintf("YetSendGoods-%+v-%+v", searchCondition.WxTagId, goods.Url)
	_, ok := get(key)
	if ok {
		return ok
	}
	set(key, key, time.Hour*30)
	return false
}
