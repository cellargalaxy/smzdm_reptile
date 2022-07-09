package spider

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/smzdm_reptile/config"
	"github.com/go-resty/resty/v2"
	"regexp"
	"time"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36"

var httpClient *resty.Client

var numRegexp *regexp.Regexp
var dateRegexp *regexp.Regexp

func init() {
	var err error
	httpClient = util.CreateHttpClient(config.Config.Timeout, config.Config.Sleep, time.Minute*5, config.Config.Retry, map[string]string{"User-Agent": userAgent}, true)

	numRegexp, err = regexp.Compile("\\d+(\\.\\d+)*")
	if err != nil {
		panic(err)
	}

	dateRegexp, err = regexp.Compile("\\d\\d:\\d\\d")
	if err != nil {
		panic(err)
	}
}
