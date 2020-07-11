package service

import (
	"encoding/json"
	"github.com/cellargalaxy/smzdm-reptile/dao"
	"github.com/cellargalaxy/smzdm-reptile/model"
	"github.com/sirupsen/logrus"
)

func AddSearchConditions(searchConditionsJsonString string) error {
	var searches []model.SearchCondition
	err := json.Unmarshal([]byte(searchConditionsJsonString), &searches)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Error("反序列化搜索条件json失败")
		return err
	}
	return dao.InsertSearchConditions(searches)
}

func ListSearchCondition() ([]model.SearchCondition, error) {
	return dao.SelectSearchConditions()
}
