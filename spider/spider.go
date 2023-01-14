package spider

import (
	"github.com/cellargalaxy/go_common/util"
	"regexp"
)

var httpClient = util.GetHttpClientSpider()

var numRegexp *regexp.Regexp
var dateRegexp *regexp.Regexp

func init() {
	var err error

	numRegexp, err = regexp.Compile("\\d+(\\.\\d+)*")
	if err != nil {
		panic(err)
	}

	dateRegexp, err = regexp.Compile("\\d\\d:\\d\\d")
	if err != nil {
		panic(err)
	}
}
