package unique

import (
	"github.com/sony/sonyflake"
	"strconv"
	"time"
)

var flake = sonyflake.NewSonyflake(sonyflake.Settings{
	StartTime: time.Date(2020, 4, 13, 0, 0, 0, 0, time.Local),
})

// ID 获取一个新的唯一id
func ID() uint64 {
	id, err := flake.NextID()
	if err != nil {
		return 0
	}
	return id
}

// IDStr 获取较短的一个ID值
func IDStr() string {
	id, err := flake.NextID()
	if err != nil {
		return ""
	}

	return strconv.FormatUint(id, 10)
}
