package utils

import (
	"errors"
	"fmt"
	"github.com/uniplaces/carbon"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

//驼峰组合单词转换为下划线分隔
func HumpToUnderline(str string) string {
	strLen := len(str)
	res := make([]string, 0)
	index := 0
	seq := 0
	for i := 1; i < strLen; i++ {
		if str[i] >= 'A' && str[i] <= 'Z' {
			if seq == i-1 {
				seq = i
			} else {
				res = append(res, str[index:i])
				index = i
				seq = i
			}
		}
		if i == strLen-1 {
			res = append(res, str[index:])
		}
	}
	return strings.ToLower(strings.Join(res, "_"))
}

func DateRange(t int64) (start, end int64) {
	current := time.Unix(t, 0)
	start = current.Add(-time.Duration(current.Hour()) * time.Hour).Add(-time.Duration(current.Minute()) * time.Minute).Add(-time.Duration(current.Second()) * time.Second).Unix()
	end = start + 24*3600
	return
}

func IsDuplicateError(err string, dbName string) bool {
	switch dbName {
	case "mysql":
		if strings.Index(err, "Error 1062:") != -1 {
			return true
		}
	case "mongo":
		if strings.Index(err, "E11000") != -1 {
			return true
		}
	}
	return false
}

func IsNotExistError(err string, dbName string) bool {
	switch dbName {
	case "mysql":
		if strings.Index(err, "record not found") != -1 {
			return true
		}
	case "mongo":
		if strings.Index(err, "no documents in result") != -1 {
			return true
		}
	case "elastic":
		if strings.Index(err, "Error 404 (Not Found)") != -1 {
			return true
		}
	}
	return false
}

// HttpMethod 将http method转化为对应数字
func HttpMethod(method string) int8 {
	switch method {
	case "POST":
		return 1
	case "DELETE":
		return 2
	case "PUT":
		return 4
	case "GET":
		return 8
	default:
		return 0
	}
}

var added int64 = 0

// Random 生成随机数字串 主要用于4位数生成短信验证码
func Random(n int) string {
	const alphanum = "0123456789"
	var bytes = make([]byte, n)
	v := atomic.AddInt64(&added, 1)
	rand.Seed(time.Now().UnixNano() + v)
	for i := range bytes {
		bytes[i] = alphanum[rand.Intn(len(alphanum))] // 返回一个取值范围在[0,n)的伪随机int值
	}
	return string(bytes)
}

// RandomAlphaNum 生成随机字符串
func RandomAlphaNum(length int) string {
	var (
		r     = rand.New(rand.NewSource(time.Now().Unix()))
		words = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		str   []byte
	)

	for i := 0; i < length; i++ {
		num := r.Intn(61)
		str = append(str, words[num])
	}

	return string(str)
}

// ToString 转换为string
func ToString(v interface{}) (vv string, err error) {
	switch v.(type) {
	case uint8, int8, int, uint, int16, int32, int64, uint64:
		vv = fmt.Sprintf("%d", v)
	case string:
		vv = v.(string)
	case float64:
		temp := v.(float64)
		vv = strconv.FormatFloat(temp, 'f', 10, 64)
	case float32:
		temp := v.(float32)
		vv = strconv.FormatFloat(float64(temp), 'f', 10, 64)
	default:
		err = errors.New("type is errors")
	}
	if vv == "" {
		err = errors.New("fail: value is empty")
	} else {
		vv = strings.TrimSpace(vv)
	}
	return
}

// ConversionMap 将结构体转换为 map
func ConversionMap(d interface{}) map[string]interface{} {
	ref := reflect.TypeOf(d)
	refValue := reflect.ValueOf(d)
	var resData = make(map[string]interface{})
	if refValue.Kind() == reflect.Struct {
		for i := 0; i < refValue.NumField(); i++ {
			filed := ref.Field(i)
			tag := filed.Tag.Get("json")
			if tag == "-" || tag == "" {
				continue
			}
			resData[tag] = refValue.Field(i).Interface()
		}
	}
	return resData
}

// 根据start,end 时间戳计算 聚合bucket 参数
func BucketForTime(start, end int64) (timeList []int64, err error) {
	if start > end {
		return nil, errors.New("param error: start > end")
	}
	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)
	start = startTime.Add(-time.Duration(startTime.Hour()) * time.Hour).Add(-time.Duration(startTime.Minute()) * time.Minute).Add(-time.Duration(startTime.Second()) * time.Second).Unix()
	end = endTime.Add(-time.Duration(endTime.Hour()) * time.Hour).Add(-time.Duration(endTime.Minute()) * time.Minute).Add(-time.Duration(endTime.Second()) * time.Second).Unix()
	end = end + 86400
	timeList = make([]int64, 0)
	for i := int64(0); ; i++ {
		if start+i*86400 > end {
			break
		}
		timeList = append(timeList, start+i*86400)
	}
	/*if len(timeList) > 30 {
		err = errors.New("时间区间过大，请控制在30天内!")
		return
	}*/
	return
}

// 环比计算
func MoMRatio(now, pre float32) float32 {
	ratio := float32(100)
	if pre == 0 {
		if now == 0 {
			return 0
		} else {
			return 100
		}
	} else {
		value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (now-pre)*100/pre), 10)
		ratio = float32(value)
	}
	return ratio
}

// 比值计算
func Ratio(now, pre float32) float32 {
	ratio := float32(100)
	if pre == 0 {
		if now == 0 {
			return 0
		} else {
			return 100
		}
	} else {
		value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (now)*100/pre), 10)
		ratio = float32(value)
	}
	return ratio
}

// FormatFloat 格式化浮点类型，转为需要的小数位数  d 保留位数
func FormatFloat(n float64, d int) float64 {
	fmtStr := "%." + fmt.Sprintf("%d", d) + "f"
	str := fmt.Sprintf(fmtStr, n)
	newF, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return n
	}
	return newF
}

// FormatFloatRatio 获取乘以100后的比例，转为需要的小数位数  d 保留位数
func FormatFloatRatio(n float64, d int) float64 {
	return FormatFloat(n, d+2) * 100
}

// 通过时间戳获取时间对象
func NewCarbonByTime(s int64) (*carbon.Carbon, error) {
	return carbon.CreateFromTimestamp(s, "Asia/Shanghai")
}

// LinkRatio 环比
func LinkRatio(total, preTotal float64) float64 {
	ratio := FormatFloatRatio(1, 2)
	diff := total - preTotal
	if diff == 0 {
		return 0
	}
	if preTotal > 0 {
		ratio = FormatFloatRatio(diff/preTotal, 2)
	}

	return ratio
}
