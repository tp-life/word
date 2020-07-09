package holiday

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func init() {
	Register("tianapi", &TianAPI{})
}

// Interface 节假日
type Interface interface {
	Set(key, val string)
	IsHoliday(time time.Time) bool
}

// TimeMaster 时间管理
type TimeMaster struct {
	list map[string]bool
	mu   sync.RWMutex
}

var (
	drivers = make(map[string]Interface)
	list    = new(TimeMaster)
)

func (timeMaster *TimeMaster) Set(day string, holiday bool) {
	timeMaster.mu.Lock()
	defer timeMaster.mu.Unlock()
	if timeMaster.list == nil {
		timeMaster.list = make(map[string]bool)
	}
	timeMaster.list[day] = holiday
}

func (timeMaster *TimeMaster) Has(day string) bool {
	timeMaster.mu.RLock()
	defer timeMaster.mu.RUnlock()
	if timeMaster.list == nil {
		return false
	}

	var _, ok = timeMaster.list[day]
	return ok
}

func (timeMaster *TimeMaster) Get(day string) bool {
	timeMaster.mu.RLock()
	defer timeMaster.mu.RUnlock()
	if timeMaster.list != nil {
		return timeMaster.list[day]
	}
	return false
}

// Holiday 节假日
type Holiday struct {
	driver Interface
}

// Register 注册一个查询驱动
func Register(name string, driver Interface) {
	drivers[name] = driver
}

// New 选择一个查询器实例化
func New(name string) *Holiday {
	var holiday = &Holiday{driver: drivers[name]}
	return holiday
}

// Set 设置参数
func (holiday *Holiday) Set(key, val string) *Holiday {
	holiday.driver.Set(key, val)
	return holiday
}

// IsHoliday 判断是否是节假日
func (holiday *Holiday) IsHoliday(time time.Time) bool {
	var day = time.Format("2006-01-02")
	if list.Has(day) {
		return list.Get(day)
	}
	var isHoliday = holiday.driver.IsHoliday(time)
	list.Set(day, isHoliday)
	return isHoliday
}

// TianAPI 天行 API
type TianAPI struct {
	Key string
}

// Set 设置参数
func (api *TianAPI) Set(key, val string) {
	switch key {
	case "key":
		api.Key = val
	}
}

// IsHoliday 判断是否是节假日
func (api *TianAPI) IsHoliday(t time.Time) bool {
	url := fmt.Sprintf("http://api.tianapi.com/txapi/jiejiari/index?key=%s&date=%s", api.Key, t.Format("2006-01-02"))
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var data map[string]interface{}
	_ = jsoniter.Unmarshal(body, &data)
	if info, ok := data["newslist"]; ok {
		if items, ok := info.([]interface{}); ok && len(items) > 0 {
			var (
				item      = items[0].(map[string]interface{})
				isNotWork = item["isnotwork"]
			)
			switch isNotWork.(type) {
			case int, int8, int16, int32, int64, float32, float64:
				return isNotWork.(float64) > 0
			case string:
				rs, _ := strconv.ParseBool(isNotWork.(string))
				return rs
			}
		}
	}

	return false
}
