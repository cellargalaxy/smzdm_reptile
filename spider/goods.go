package spider

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/smzdm_reptile/config"
	"github.com/cellargalaxy/smzdm_reptile/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ListGoods(ctx context.Context, searchKey string) ([]model.Goods, error) {
	var goodses []model.Goods
	for page := 1; page <= config.Config.MaxPage; page = page + 1 {
		util.SleepWare(ctx, config.Config.Sleep)
		html, err := requestListGoods(ctx, searchKey, page)
		if err != nil {
			continue
		}
		gs, err := analysisListGoods(ctx, html)
		if err != nil {
			continue
		}
		if len(gs) == 0 {
			break
		}
		goodses = append(goodses, gs...)
	}
	return goodses, nil
}

// 商品列表页面
func analysisListGoods(ctx context.Context, html string) ([]model.Goods, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("商品列表页面，html解析异常")
		return nil, err
	}

	toyearString := time.Now().Format("2006-")
	todayString := time.Now().Format("2006-01-02 ")

	var goodses []model.Goods
	doc.Find(".feed-row-wide").Each(func(i int, goodsSelection *goquery.Selection) {
		var goods model.Goods

		//title,url
		titleSelection := goodsSelection.Find(".feed-nowrap")
		if titleSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("商品列表页面，title为空")
			return
		}
		titleSelection = titleSelection.First()
		if titleSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("商品列表页面，title为空")
			return
		}
		title := strings.TrimSpace(titleSelection.Text())
		if title == "" {
			logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("商品列表页面，title为空")
			return
		}
		goods.Title = title
		url, exists := titleSelection.Attr("href")
		if !exists {
			logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("商品列表页面，url为空")
			return
		}
		url = strings.TrimSpace(url)
		if url == "" {
			logrus.WithContext(ctx).WithFields(logrus.Fields{}).Warn("商品列表页面，url为空")
			return
		}
		goods.Url = url

		//price
		priceSelection := goodsSelection.Find(".z-highlight")
		if priceSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，price为空")
			return
		}
		priceSelection = priceSelection.First()
		if priceSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，price为空")
			return
		}
		priceString := strings.TrimSpace(priceSelection.Text())
		if !numRegexp.MatchString(priceString) {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，price为空")
			return
		}
		priceString = numRegexp.FindString(priceString)
		price, err := strconv.ParseFloat(priceString, 32)
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Warn("商品列表页面，price反序列化异常")
			return
		}
		goods.Price = price

		//zhi
		zhiSelection := goodsSelection.Find(".price-btn-up")
		if zhiSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，zhi为空")
			return
		}
		zhiSelection = zhiSelection.First()
		if zhiSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，zhi为空")
			return
		}
		zhiSelection = zhiSelection.Find(".unvoted-wrap")
		if zhiSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，zhi为空")
			return
		}
		zhiSelection = zhiSelection.First()
		if zhiSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，zhi为空")
			return
		}
		zhiSelection = zhiSelection.Find("span")
		if zhiSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，zhi为空")
			return
		}
		zhiSelection = zhiSelection.First()
		if zhiSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，zhi为空")
			return
		}
		zhiString := strings.TrimSpace(zhiSelection.Text())
		if !numRegexp.MatchString(zhiString) {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，zhi为空")
			return
		}
		zhiString = numRegexp.FindString(zhiString)
		zhi, err := strconv.Atoi(zhiString)
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Warn("商品列表页面，zhi反序列化异常")
			return
		}
		goods.Zhi = zhi

		//buzhi
		buzhiSelection := goodsSelection.Find(".price-btn-down")
		if buzhiSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，buzhi为空")
			return
		}
		buzhiSelection = buzhiSelection.First()
		if buzhiSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，buzhi为空")
			return
		}
		buzhiSelection = buzhiSelection.Find(".unvoted-wrap")
		if buzhiSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，buzhi为空")
			return
		}
		buzhiSelection = buzhiSelection.First()
		if buzhiSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，buzhi为空")
			return
		}
		buzhiSelection = buzhiSelection.Find("span")
		if buzhiSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，buzhi为空")
			return
		}
		buzhiSelection = buzhiSelection.First()
		if buzhiSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，buzhi为空")
			return
		}
		buzhiString := strings.TrimSpace(buzhiSelection.Text())
		if !numRegexp.MatchString(buzhiString) {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，buzhi为空")
			return
		}
		buzhiString = numRegexp.FindString(buzhiString)
		buzhi, err := strconv.Atoi(buzhiString)
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Warn("商品列表页面，buzhi反序列化异常")
			return
		}
		goods.Buzhi = buzhi

		//date,merchant
		merchantSelection := goodsSelection.Find(".feed-block-extras")
		if merchantSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，merchant为空")
			return
		}
		merchantSelection = merchantSelection.First()
		if merchantSelection == nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，merchant为空")
			return
		}
		merchantString := strings.TrimSpace(merchantSelection.Text())
		ss := strings.Split(merchantString, " ")
		if len(ss) < 2 {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，merchant非法")
			return
		}
		dateString := strings.Join(ss[:len(ss)-1], " ")
		dateString = strings.TrimSpace(dateString)
		if !dateRegexp.MatchString(dateString) {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，date非法")
			return
		}
		date, err := time.Parse("2006-01-02 15:04", toyearString+dateString)
		if err != nil {
			date, err = time.Parse("2006-01-02 15:04", todayString+dateString)
		}
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Debug("商品列表页面，date非法")
			return
		}
		if time.Hour*24 < time.Now().Sub(date) {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Debug("商品列表页面，商品过期")
			return
		}
		goods.Date = date
		merchant := strings.TrimSpace(ss[len(ss)-1])
		if merchant == "" {
			logrus.WithContext(ctx).WithFields(logrus.Fields{"url": url}).Warn("商品列表页面，merchant为空")
			return
		}
		goods.Merchant = merchant

		logrus.WithContext(ctx).WithFields(logrus.Fields{"goods": goods}).Info("创建商品对象")
		goodses = append(goodses, goods)
	})
	return goodses, nil
}

// 商品列表页面
func requestListGoods(ctx context.Context, searchKey string, page int) (string, error) {
	response, err := httpClient.R().SetContext(ctx).
		SetCookies(listCookie(ctx)).
		SetQueryParam("c", "home").
		SetQueryParam("s", searchKey).
		SetQueryParam("v", "a").
		SetQueryParam("mx_v", "b").
		SetQueryParam("p", fmt.Sprintf("%d", page)).
		Get("https://search.smzdm.com/")

	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("商品列表页面，请求异常")
		return "", fmt.Errorf("商品列表页面，请求异常")
	}
	if response == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("商品列表页面，响应为空")
		return "", fmt.Errorf("商品列表页面，响应为空")
	}
	setCookie(ctx, response.Cookies())
	statusCode := response.StatusCode()
	body := response.String()
	logrus.WithContext(ctx).WithFields(logrus.Fields{"statusCode": statusCode, "len(body)": len(body)}).Info("商品列表页面，响应")
	if statusCode != http.StatusOK {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"StatusCode": statusCode}).Error("商品列表页面，响应码失败")
		return "", fmt.Errorf("商品列表页面，响应码失败: %+v", statusCode)
	}
	return body, nil
}

var cookieTime time.Time
var cookies []*http.Cookie

func listCookie(ctx context.Context) []*http.Cookie {
	if len(cookies) > 0 && time.Now().Sub(cookieTime).Minutes() < 10 {
		return cookies
	}

	response, err := httpClient.R().SetContext(ctx).
		SetCookies(cookies).
		Get("https://www.smzdm.com/")

	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("首页，请求异常")
		return cookies
	}
	if response == nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"err": err}).Error("首页，响应为空")
		return cookies
	}
	setCookie(ctx, response.Cookies())

	return cookies
}
func setCookie(ctx context.Context, list []*http.Cookie) {
	if len(list) == 0 {
		return
	}
	cookies = list
	cookieTime = time.Now()
}
