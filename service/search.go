package service

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cellargalaxy/smzdm-reptile/config"
	"github.com/cellargalaxy/smzdm-reptile/dao"
	"github.com/cellargalaxy/smzdm-reptile/model"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var sentGoodsMap = make(map[string]model.Goods)

func StartSearchService() {
	logrus.Info("开始搜索服务")
	for {
		searchConditions, err := dao.SelectSearchConditions()
		if err != nil {
			time.Sleep(config.Sleep)
			continue
		}
		for _, searchCondition := range searchConditions {
			if !searchAndSend(searchCondition) {
				logrus.WithFields(logrus.Fields{"searchCondition": searchCondition}).Error("搜索或者发送商品失败")
			}
			time.Sleep(config.Sleep)
		}
		time.Sleep(config.Sleep)
	}
	logrus.Info("结束搜索服务")
}

func searchAndSend(searchCondition model.SearchCondition) bool {
	goodses, err := searchGoods(searchCondition)
	if err != nil {
		return false
	}

	for _, goods := range goodses {
		if isSentGoods(searchCondition, goods) {
			continue
		}
		addSentGoods(searchCondition, goods)
		data := map[string]interface{}{
			"id":       searchCondition.SearchId,
			"title":    goods.Title,
			"price":    fmt.Sprint(goods.Price),
			"zhi":      fmt.Sprint(goods.Zhi),
			"buzhi":    fmt.Sprint(goods.Buzhi),
			"merchant": goods.Merchant,
		}
		SendWxPush(searchCondition.WxTemplateId, searchCondition.WxTagId, goods.Url, data)
	}
	return true
}

func isSentGoods(searchCondition model.SearchCondition, newGoods model.Goods) bool {
	today := time.Now()
	var expiredKeys []string
	for key, goods := range sentGoodsMap {
		if goods.Date.Day()+2 < today.Day() {
			expiredKeys = append(expiredKeys, key)
		}
	}
	for _, expiredKey := range expiredKeys {
		delete(sentGoodsMap, expiredKey)
	}

	newKey := fmt.Sprintf("%+v:%+v", searchCondition.WxTagId, newGoods.Url)
	for key := range sentGoodsMap {
		if newKey == key {
			return true
		}
	}
	return false
}

func addSentGoods(searchCondition model.SearchCondition, newGoods model.Goods) {
	newKey := fmt.Sprintf("%+v:%+v", searchCondition.WxTagId, newGoods.Url)
	sentGoodsMap[newKey] = newGoods
}

//搜索商品
func searchGoods(searchCondition model.SearchCondition) ([]model.Goods, error) {
	logrus.WithFields(logrus.Fields{"searchCondition": searchCondition}).Info("搜索商品")

	titleContainRegular, err := regexp.Compile(searchCondition.TitleContain)
	if err != nil {
		logrus.WithFields(logrus.Fields{"TitleContain": searchCondition.TitleContain, "err": err}).Error("创建标题包含正则对象失败")
		return nil, err
	}
	if searchCondition.TitleContain == "" {
		logrus.Info("标题包含正则为空，取消创建标题包含正则对象")
		titleContainRegular = nil
	}

	titleExcludeRegular, err := regexp.Compile(searchCondition.TitleExclude)
	if err != nil {
		logrus.WithFields(logrus.Fields{"TitleExclude": searchCondition.TitleExclude, "err": err}).Error("创建标题排除正则对象失败")
		return nil, err
	}
	if searchCondition.TitleExclude == "" {
		logrus.Info("标题排除正则为空，取消创建标题排除正则对象")
		titleExcludeRegular = nil
	}

	merchantContainRegular, err := regexp.Compile(searchCondition.MerchantContain)
	if err != nil {
		logrus.WithFields(logrus.Fields{"MerchantContain": searchCondition.MerchantContain, "err": err}).Error("创建商家包含正则对象失败")
		return nil, err
	}
	if searchCondition.MerchantContain == "" {
		logrus.Info("商家包含正则为空，取消创建商家包含正则对象")
		merchantContainRegular = nil
	}

	merchantExcludeRegular, err := regexp.Compile(searchCondition.MerchantExclude)
	if err != nil {
		logrus.WithFields(logrus.Fields{"MerchantExclude": searchCondition.MerchantExclude, "err": err}).Error("创建商家排除正则对象失败")
		return nil, err
	}
	if searchCondition.MerchantExclude == "" {
		logrus.Info("商家排除正则为空，取消创建商家排除正则对象")
		merchantExcludeRegular = nil
	}

	var goodses []model.Goods
	for page := 1; page <= config.MaxPage; page = page + 1 {
		html, err := requestListGoods(searchCondition.SearchKey, page)
		if err != nil {
			logrus.WithFields(logrus.Fields{"err": err}).Error("请求商品列表页面失败")
			continue
		}
		gs, err := analysisListGoods(html)
		if err != nil {
			logrus.WithFields(logrus.Fields{"err": err}).Error("分析商品列表页面")
			continue
		}
		if len(gs) == 0 {
			logrus.WithFields(logrus.Fields{"page": page}).Info("遍历完成全部商品列表页面")
			break
		}

		for _, goods := range gs {
			if isMeetCondition(goods, searchCondition, titleContainRegular, titleExcludeRegular, merchantContainRegular, merchantExcludeRegular) {
				goodses = append(goodses, goods)
			}
		}

		time.Sleep(config.Sleep)
	}
	return goodses, nil
}

func isMeetCondition(goods model.Goods, searchCondition model.SearchCondition,
	titleContainRegular *regexp.Regexp, titleExcludeRegular *regexp.Regexp,
	merchantContainRegular *regexp.Regexp, merchantExcludeRegular *regexp.Regexp) bool {
	if goods.Price < searchCondition.MinPrice || goods.Price > searchCondition.MaxPrice {
		logrus.WithFields(logrus.Fields{"MinPrice": searchCondition.MinPrice, "Price": goods.Price, "MaxPrice": searchCondition.MaxPrice, "url": goods.Url}).Info("商品【价格】不在范围内")
		return false
	}
	if goods.Zhi < searchCondition.MinZhi {
		logrus.WithFields(logrus.Fields{"Zhi": goods.Zhi, "MinZhi": searchCondition.MinZhi, "url": goods.Url}).Info("商品【值】不在范围内")
		return false
	}
	if goods.Buzhi > searchCondition.MaxBuzhi {
		logrus.WithFields(logrus.Fields{"Buzhi": goods.Buzhi, "MaxBuzhi": searchCondition.MaxBuzhi, "url": goods.Url}).Info("商品【不值】不在范围内")
		return false
	}
	if titleContainRegular != nil && !titleContainRegular.MatchString(goods.Title) {
		logrus.WithFields(logrus.Fields{"Title": goods.Title, "TitleContain": searchCondition.TitleContain, "url": goods.Url}).Info("【标题】被包含正则过滤")
		return false
	}
	if titleExcludeRegular != nil && titleExcludeRegular.MatchString(goods.Title) {
		logrus.WithFields(logrus.Fields{"Title": goods.Title, "TitleExclude": searchCondition.TitleExclude, "url": goods.Url}).Info("【标题】被排除正则过滤")
		return false
	}
	if merchantContainRegular != nil && !merchantContainRegular.MatchString(goods.Merchant) {
		logrus.WithFields(logrus.Fields{"Merchant": goods.Merchant, "MerchantContain": searchCondition.MerchantContain, "url": goods.Url}).Info("【商家】被包含正则过滤")
		return false
	}
	if merchantExcludeRegular != nil && merchantExcludeRegular.MatchString(goods.Merchant) {
		logrus.WithFields(logrus.Fields{"Merchant": goods.Merchant, "MerchantExclude": searchCondition.MerchantExclude, "url": goods.Url}).Info("【商家】被排除正则过滤")
		return false
	}
	return true
}

//分析商品列表页面
func analysisListGoods(html string) ([]model.Goods, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Error("商品列表页面，html解析失败")
		return nil, err
	}

	numRegexp, err := regexp.Compile("\\d+(\\.\\d+)*")
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Error("创建数字正则对象失败")
		return nil, err
	}

	dateRegexp, err := regexp.Compile("\\d\\d:\\d\\d")
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Error("创建日期正则对象失败")
		return nil, err
	}

	todayString := time.Now().Format("2006-01-02 ")

	var goodses []model.Goods
	doc.Find(".feed-row-wide").Each(func(i int, goodsSelection *goquery.Selection) {
		//title
		titlesSelection := goodsSelection.Find(".feed-nowrap")
		titleSelection := titlesSelection.First()
		if titleSelection == nil || titleSelection.Text() == "" {
			//logrus.Warn("商品【标题】为空")
			return
		}

		//price
		pricesSelection := goodsSelection.Find(".z-highlight")
		priceSelection := pricesSelection.First()
		if priceSelection == nil || priceSelection.Text() == "" {
			//logrus.Warn("商品【价格】为空")
			return
		}

		//zhi
		zhisSelection1 := goodsSelection.Find(".price-btn-up")
		zhiSelection1 := zhisSelection1.First()
		if zhiSelection1 == nil || zhiSelection1.Text() == "" {
			//logrus.Warn("商品【值1】为空")
			return
		}
		zhisSelection2 := zhiSelection1.Find(".unvoted-wrap")
		zhiSelection2 := zhisSelection2.First()
		if zhiSelection2 == nil || zhiSelection2.Text() == "" {
			//logrus.Warn("商品【值2】为空")
			return
		}
		zhisSelection3 := zhiSelection2.Find("span")
		zhiSelection3 := zhisSelection3.First()
		if zhiSelection3 == nil || zhiSelection3.Text() == "" {
			//logrus.Warn("商品【值3】为空")
			return
		}

		//buzhi
		buzhisSelection1 := goodsSelection.Find(".price-btn-down")
		buzhiSelection1 := buzhisSelection1.First()
		if buzhiSelection1 == nil || buzhiSelection1.Text() == "" {
			//logrus.Warn("商品【不值1】为空")
			return
		}
		buzhisSelection2 := buzhiSelection1.Find(".unvoted-wrap")
		buzhiSelection2 := buzhisSelection2.First()
		if buzhiSelection2 == nil || buzhiSelection2.Text() == "" {
			//logrus.Warn("商品【不值2】为空")
			return
		}
		buzhisSelection3 := buzhiSelection2.Find("span")
		buzhiSelection3 := buzhisSelection3.First()
		if buzhiSelection3 == nil || buzhiSelection3.Text() == "" {
			//logrus.Warn("商品【不值3】为空")
			return
		}

		//date,merchant
		merchantsSelection := goodsSelection.Find(".feed-block-extras")
		merchantSelection := merchantsSelection.First()
		if merchantSelection == nil || merchantSelection.Text() == "" {
			//logrus.Warn("商品【日期与商家】为空")
			return
		}

		//title
		title := strings.TrimSpace(titleSelection.Text())

		url, exists := titleSelection.Attr("href")
		if !exists {
			//logrus.WithFields(logrus.Fields{"title": title}).Warn("商品【链接】为空")
			return
		}

		//price
		priceString := strings.TrimSpace(priceSelection.Text())
		if !numRegexp.MatchString(priceString) {
			//logrus.WithFields(logrus.Fields{"title": title, "priceString": priceString}).Warn("商品【价格】非法")
			return
		}
		priceString = numRegexp.FindString(priceString)
		price, err := strconv.ParseFloat(priceString, 32)
		if err != nil {
			//logrus.WithFields(logrus.Fields{"title": title, "priceString": priceString, "err": err}).Warn("商品【价格】格式化为数字失败")
			return
		}

		//zhi
		zhiString := strings.TrimSpace(zhiSelection3.Text())
		if !numRegexp.MatchString(zhiString) {
			//logrus.WithFields(logrus.Fields{"title": title, "zhiString": zhiString}).Warn("商品【值】非法")
			return
		}
		zhiString = numRegexp.FindString(zhiString)
		zhi, err := strconv.ParseInt(zhiString, 10, 32)
		if err != nil {
			//logrus.WithFields(logrus.Fields{"title": title, "zhiString": zhiString, "err": err}).Warn("商品【值】格式化为数字失败")
			return
		}

		//buzhi
		buzhiString := strings.TrimSpace(buzhiSelection3.Text())
		if !numRegexp.MatchString(buzhiString) {
			//logrus.WithFields(logrus.Fields{"title": title, "buzhiString": buzhiString}).Warn("商品【不值】非法")
			return
		}
		buzhiString = numRegexp.FindString(buzhiString)
		buzhi, err := strconv.ParseInt(buzhiString, 10, 32)
		if err != nil {
			//logrus.WithFields(logrus.Fields{"title": title, "buzhiString": buzhiString, "err": err}).Warn("商品【不值】格式化为数字失败")
			return
		}

		//date,merchant
		merchantString := strings.TrimSpace(merchantSelection.Text())
		if !strings.Contains(merchantString, " ") {
			//logrus.WithFields(logrus.Fields{"title": title, "merchantString": merchantString}).Warn("商品【日期与商家】非法")
			return
		}
		ss := strings.Split(merchantString, " ")
		dateString := ss[0]
		if len(ss) > 2 {
			for i := 1; i < len(ss)-1; i++ {
				dateString = dateString + " " + ss[i]
			}
		}
		dateString = strings.TrimSpace(dateString)
		merchant := strings.TrimSpace(ss[len(ss)-1])
		if !dateRegexp.MatchString(dateString) {
			//logrus.WithFields(logrus.Fields{"title": title, "dateString": dateString}).Warn("商品【日期】非法")
			return
		}
		date, err := time.Parse("2006-01-02 15:04", todayString+dateString)
		if err != nil {
			//logrus.WithFields(logrus.Fields{"title": title, "todayString+dateString": todayString + dateString, "err": err}).Warn("商品【日期】格式化为日期失败")
			return
		}

		goods := model.Goods{
			Title:    title,
			Url:      url,
			Price:    float32(price),
			Zhi:      int(zhi),
			Buzhi:    int(buzhi),
			Merchant: merchant,
			Date:     date,
		}
		logrus.WithFields(logrus.Fields{"goods": goods}).Info("创建商品对象")
		goodses = append(goodses, goods)
	})
	return goodses, nil
}

//请求商品列表页面
func requestListGoods(searchKey string, page int) (string, error) {
	request := gorequest.New()
	response, body, errs := request.Get("https://search.smzdm.com").
		Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36").
		Param("c", "home").
		Param("s", searchKey).
		Param("v", "b").
		Param("p", fmt.Sprintf("%d", page)).
		Timeout(config.Timeout).
		End()
	logrus.WithFields(logrus.Fields{"errs": errs}).Info("获取商品列表页面请求")
	if errs != nil && len(errs) > 0 {
		logrus.Error("获取商品列表页面请求异常")
		return "", fmt.Errorf("获取商品列表页面请求异常")
	}
	logrus.WithFields(logrus.Fields{"StatusCode": response.StatusCode, "body": len(body)}).Info("获取商品列表页面请求")
	if response.StatusCode != 200 {
		logrus.Error("获取商品列表页面请求响应码异常")
		return "", fmt.Errorf("获取商品列表页面请求响应码异常")
	}
	return body, nil
}
