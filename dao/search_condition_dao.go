package dao

import (
	"encoding/json"
	"github.com/cellargalaxy/smzdm-reptile/config"
	"github.com/cellargalaxy/smzdm-reptile/model"
	"github.com/cellargalaxy/smzdm-reptile/utils"
	"github.com/sirupsen/logrus"
)

var searchConditions []model.SearchCondition

func InsertSearchConditions(searches []model.SearchCondition) error {
	bytes, err := json.Marshal(searches)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Error("序列化搜索条件失败")
		return err
	}
	err = utils.WriteFileWithBytesOrCreateIfNotExist(config.SearchConditionFilePath, bytes)
	if err == nil {
		searchConditions = searches
	}
	return err
}

func SelectSearchConditions() ([]model.SearchCondition, error) {
	if searchConditions != nil && len(searchConditions) > 0 {
		return searchConditions, nil
	}
	jsonString, err := utils.ReadFileOrCreateIfNotExist(config.SearchConditionFilePath, "[]")
	if err != nil {
		return nil, err
	}
	var searches []model.SearchCondition
	err = json.Unmarshal([]byte(jsonString), &searches)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Error("反序列化搜索条件失败")
		return nil, err
	}
	searchConditions = searches
	return searchConditions, nil
}
