package spider

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/go-resty/resty/v2"
	"regexp"
)

var httpClient *resty.Client

var numRegexp *regexp.Regexp
var dateRegexp *regexp.Regexp

func init() {
	var err error
	httpClient = util.GetHttpClientRetry()

	numRegexp, err = regexp.Compile("\\d+(\\.\\d+)*")
	if err != nil {
		panic(err)
	}

	dateRegexp, err = regexp.Compile("\\d\\d:\\d\\d")
	if err != nil {
		panic(err)
	}
}
