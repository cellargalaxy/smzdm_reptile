package controller

import (
	"github.com/cellargalaxy/smzdm-reptile/config"
	"github.com/cellargalaxy/smzdm-reptile/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func StartWebService() {
	logrus.Info("开始web服务")
	engine := gin.Default()
	engine.GET("/", func(context *gin.Context) {
		context.Header("Content-Type", "text/html; charset=utf-8")
		context.String(200, indexHtmlString)
	})
	engine.GET("/listSearchCondition", func(context *gin.Context) {
		context.JSON(http.StatusOK, createResponse(service.ListSearchCondition()))
	})
	engine.POST("/saveSearchConditions", func(context *gin.Context) {
		searchConditionsJsonString := context.PostForm("searchConditions")
		context.JSON(http.StatusOK, createResponse(nil, service.AddSearchConditions(searchConditionsJsonString)))
	})
	engine.Run(config.ListenAddress)
	defer logrus.Info("结束web服务")
}

func createResponse(data interface{}, err error) map[string]interface{} {
	if err == nil {
		return gin.H{"code": config.SuccessCode, "message": nil, "data": data}
	} else {
		return gin.H{"code": config.FailCode, "message": err.Error(), "data": data}
	}
}

var indexHtmlString = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>张大妈爬虫</title>
    <link type="text/css" rel="stylesheet" href="//unpkg.com/bootstrap/dist/css/bootstrap.min.css"/>
    <link type="text/css" rel="stylesheet" href="//unpkg.com/bootstrap-vue@latest/dist/bootstrap-vue.min.css"/>
</head>
<body>
<div id="app">
    <b-button-group style="width: 100%">
        <b-button variant="primary" @click="saveSearchCondition">save</b-button>
        <b-button variant="info" @click="listSearchCondition">flush</b-button>
    </b-button-group>
    <b-form-textarea :rows="rows" v-model="searchConditionString" @input="flushRows"></b-form-textarea>
</div>
</body>
<script src="//polyfill.io/v3/polyfill.min.js?features=es2015%2CIntersectionObserver" crossorigin="anonymous"></script>
<script src="//unpkg.com/vue@latest/dist/vue.min.js"></script>
<script src="//unpkg.com/bootstrap-vue@latest/dist/bootstrap-vue.min.js"></script>
<script src="https://cdn.bootcss.com/jquery/3.4.1/jquery.min.js"></script>
<script>
    var app = new Vue({
        el: '#app',
        data: {
            searchConditionString: "",
            rows: 1,
        },
        methods: {
            saveSearchCondition: function () {
                if (!window.confirm("确定修改？")) {
                    return
                }
                $.ajax({
                    url: 'saveSearchConditions',
                    type: 'post',
                    data: {"searchConditions": app.searchConditionString},
                    contentType: "application/x-www-form-urlencoded",
                    dataType: "json",

                    error: ajaxErrorDeal,
                    success: function (data) {
                        if (data.status == 1) {
                            alert('修改成功')
                            app.listSearchCondition()
                        } else {
                            alert('修改失败: ' + data.massage)
                        }
                    }
                });
            },
            listSearchCondition: function () {
                $.ajax({
                    url: 'listSearchCondition',
                    type: 'get',
                    data: {},
                    contentType: "application/x-www-form-urlencoded",
                    dataType: "json",

                    error: ajaxErrorDeal,
                    success: function (data) {
                        let searchConditionString = JSON.stringify(data.data, null, 2);
                        if (searchConditionString == null || searchConditionString == "") {
                            searchConditionString = "[]"
                        }
                        app.searchConditionString = searchConditionString
                        app.rows = app.searchConditionString.split("\n").length
                        alert('刷新成功')
                    }
                });
            },
            flushRows: function (text) {
                app.rows = text.split("\n").length
            },
        },
    })

    app.listSearchCondition()

    function ajaxErrorDeal() {
        alert("网络错误!");
    }
</script>
</html>`
