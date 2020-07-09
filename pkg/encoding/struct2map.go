package encoding

import (
	"word/pkg/common/log"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Struct2Map struct==>map
func Struct2Map(v interface{}) map[string]interface{} {
	by, err := json.Marshal(v)
	if err != nil {
		log.Error(err)
		return nil
	}

	var m map[string]interface{}
	err = json.Unmarshal(by, &m)
	if err != nil {
		log.Error(err)
		return m
	}
	return m
}
