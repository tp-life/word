package idcard

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"word/pkg/app"
)

var (
	// 实名认证接口
	requestURI = "https://checkid.market.alicloudapi.com/IDCard"
	// ErrQueryFailed 查询出错
	ErrQueryFailed = errors.New("query auth info error")
	cities         = map[int]string{
		11: "北京",
		12: "天津",
		13: "河北",
		14: "山西",
		15: "内蒙古",
		21: "辽宁",
		22: "吉林",
		23: "黑龙江 ",
		31: "上海",
		32: "江苏",
		33: "浙江",
		34: "安徽",
		35: "福建",
		36: "江西",
		37: "山东",
		41: "河南",
		42: "湖北 ",
		43: "湖南",
		44: "广东",
		45: "广西",
		46: "海南",
		50: "重庆",
		51: "四川",
		52: "贵州",
		53: "云南",
		54: "西藏 ",
		61: "陕西",
		62: "甘肃",
		63: "青海",
		64: "宁夏",
		65: "新疆",
		71: "台湾",
		81: "香港",
		82: "澳门",
		91: "国外",
	}
	reg, _ = regexp.Compile("^\\d{6}(18|19|20)?\\d{2}(0[1-9]|1[012])(0[1-9]|[12]\\d|3[01])\\d{3}(\\d|[xX])$")
)

// ValidIDCard 验证身份证是否符合规则
func ValidIDCard(idCard string) bool {
	if idCard == "" || !reg.Match([]byte(idCard)) {
		return false
	}

	var address, _ = strconv.Atoi(idCard[:2])
	if _, ok := cities[address]; !ok {
		return false
	}

	var chars = strings.Split(idCard, "")
	var factor = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	var parity = []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}
	var sum = 0
	var ai = 0
	var wi = 0
	for i := 0; i < 17; i++ {
		ai, _ = strconv.Atoi(chars[i])
		wi = factor[i]
		sum += ai * wi
	}
	if parity[sum%11] != strings.ToUpper(chars[17]) {
		return false
	}
	return true
}

var (
	conf config
)

type config struct {
	AppCode string `mapstructure:"app_code" toml:"app_code" env:"IDCARD_CODE"`
}

// Auth 验证用户身份信息
func Auth(username, idCard string) (bool, error) {
	var (
		uri, _ = url.Parse(requestURI)
		query  = uri.Query()
	)

	query.Set("idCard", idCard)
	query.Set("name", username)
	uri.RawQuery = query.Encode()
	request, _ := http.NewRequest(http.MethodGet, uri.String(), nil)
	request.Header.Set("Authorization", fmt.Sprintf("APPCODE %s", app.InitConfig("idauth", &conf)))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return false, ErrQueryFailed
	}

	data, _ := ioutil.ReadAll(response.Body)
	return jsoniter.Get(data, "status").ToString() == "01", nil
}
