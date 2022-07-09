package service

import (
	"context"
	"fmt"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/smzdm_reptile/cache"
	"github.com/cellargalaxy/smzdm_reptile/config"
	"github.com/cellargalaxy/smzdm_reptile/model"
	"github.com/cellargalaxy/smzdm_reptile/rpc/msg"
	"github.com/cellargalaxy/smzdm_reptile/spider"
	"github.com/sirupsen/logrus"
	"regexp"
	"time"
)

func StartSearchService() {
	for {
		ctx := util.GenCtx()
		for _, searchCondition := range config.Config.Conditions {
			err := searchAndSend(ctx, searchCondition)
			if err != nil {
				logrus.WithContext(ctx).WithFields(logrus.Fields{"searchCondition": searchCondition, "err": err}).Error("搜索并发送商品失败")
				msg.SendErr(ctx, "搜索并发送商品", err)
			}
			time.Sleep(util.WareDuration(config.Config.Sleep))
		}
		time.Sleep(util.WareDuration(config.Config.Sleep))
	}
}

func searchAndSend(ctx context.Context, searchCondition model.SearchCondition) error {
	goodses, err := searchGoods(ctx, searchCondition)
	if err != nil {
		return err
	}

	for _, goods := range goodses {
		if cache.YetSendGoods(ctx, searchCondition, goods) {
			continue
		}
		data := map[string]interface{}{
			"id":       searchCondition.SearchId,
			"title":    goods.Title,
			"price":    fmt.Sprint(goods.Price),
			"zhi":      fmt.Sprint(goods.Zhi),
			"buzhi":    fmt.Sprint(goods.Buzhi),
			"merchant": goods.Merchant,
		}
		msg.SendWxTemplateToTag(ctx, searchCondition.WxTemplateId, searchCondition.WxTagId, goods.Url, data)
	}
	return nil
}

//搜索商品
func searchGoods(ctx context.Context, searchCondition model.SearchCondition) ([]model.Goods, error) {
	logrus.WithContext(ctx).WithFields(logrus.Fields{"searchCondition": searchCondition}).Info("搜索商品")
	goodses, err := spider.ListGoods(ctx, searchCondition.SearchKey)
	if err != nil {
		return nil, err
	}
	goodses, err = filterMeetGoods(ctx, goodses, searchCondition)
	return goodses, nil
}

func filterMeetGoods(ctx context.Context, goodses []model.Goods, searchCondition model.SearchCondition) ([]model.Goods, error) {
	var titleContainRegular *regexp.Regexp
	var titleExcludeRegular *regexp.Regexp
	var merchantContainRegular *regexp.Regexp
	var merchantExcludeRegular *regexp.Regexp
	var err error

	if searchCondition.TitleContain != "" {
		titleContainRegular, err = regexp.Compile(searchCondition.TitleContain)
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"TitleContain": searchCondition.TitleContain, "err": err}).Error("创建标题包含正则对象异常")
			return nil, err
		}
	}
	if searchCondition.TitleExclude != "" {
		titleExcludeRegular, err = regexp.Compile(searchCondition.TitleExclude)
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"TitleExclude": searchCondition.TitleExclude, "err": err}).Error("创建标题排除正则对象异常")
			return nil, err
		}
	}
	if searchCondition.MerchantContain != "" {
		merchantContainRegular, err = regexp.Compile(searchCondition.MerchantContain)
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"MerchantContain": searchCondition.MerchantContain, "err": err}).Error("创建商家包含正则对象异常")
			return nil, err
		}
	}
	if searchCondition.MerchantExclude != "" {
		merchantExcludeRegular, err = regexp.Compile(searchCondition.MerchantExclude)
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"MerchantExclude": searchCondition.MerchantExclude, "err": err}).Error("创建商家排除正则对象异常")
			return nil, err
		}
	}

	gs := make([]model.Goods, 0, len(goodses))
	for i := range goodses {
		goods := goodses[i]

		if goods.Price < searchCondition.MinPrice || goods.Price > searchCondition.MaxPrice {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"MinPrice": searchCondition.MinPrice, "Price": goods.Price, "MaxPrice": searchCondition.MaxPrice, "url": goods.Url}).Info("商品【价格】不在范围内")
			continue
		}
		if goods.Zhi < searchCondition.MinZhi {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"Zhi": goods.Zhi, "MinZhi": searchCondition.MinZhi, "url": goods.Url}).Info("商品【值】不在范围内")
			continue
		}
		if goods.Buzhi > searchCondition.MaxBuzhi {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"Buzhi": goods.Buzhi, "MaxBuzhi": searchCondition.MaxBuzhi, "url": goods.Url}).Info("商品【不值】不在范围内")
			continue
		}
		if titleContainRegular != nil && !titleContainRegular.MatchString(goods.Title) {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"Title": goods.Title, "TitleContain": searchCondition.TitleContain, "url": goods.Url}).Info("【标题】被包含正则过滤")
			continue
		}
		if titleExcludeRegular != nil && titleExcludeRegular.MatchString(goods.Title) {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"Title": goods.Title, "TitleExclude": searchCondition.TitleExclude, "url": goods.Url}).Info("【标题】被排除正则过滤")
			continue
		}
		if merchantContainRegular != nil && !merchantContainRegular.MatchString(goods.Merchant) {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"Merchant": goods.Merchant, "MerchantContain": searchCondition.MerchantContain, "url": goods.Url}).Info("【商家】被包含正则过滤")
			continue
		}
		if merchantExcludeRegular != nil && merchantExcludeRegular.MatchString(goods.Merchant) {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"Merchant": goods.Merchant, "MerchantExclude": searchCondition.MerchantExclude, "url": goods.Url}).Info("【商家】被排除正则过滤")
			continue
		}

		gs = append(gs, goods)
	}

	return gs, nil
}
