package service

import (
	"context"
	"encoding/json"
	"github.com/cellargalaxy/smzdm-reptile/dao"
	"github.com/cellargalaxy/smzdm-reptile/model"
	"github.com/sirupsen/logrus"
)

func AddSearchConditions(ctx context.Context, searchConditionsJsonString string) error {
	var searches []model.SearchCondition
	err := json.Unmarshal([]byte(searchConditionsJsonString), &searches)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("反序列化搜索条件json失败")
		return err
	}
	return dao.InsertSearchConditions(ctx, searches)
}

func ListSearchCondition(ctx context.Context) ([]model.SearchCondition, error) {
	return dao.SelectSearchConditions(ctx)
}
