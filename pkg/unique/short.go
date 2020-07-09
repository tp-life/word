package unique

import (
	"math"
	"math/rand"
	"strings"
	"time"
)

var base = []string{
	"1", "2", "3", "4", "5", "6", "7", "8", "9", "0",
	"B", "C", "E", "F", "D", "G", "H", "J", "K", "L", "M",
	"N", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

const suffix = "A"
const codeLen = 7

// GenCode 随机生成短ID
func GenCode() string {
	re := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := len(base)
	n := math.Pow(float64(b), float64(codeLen))
	code := ShortCode(re.Int63n(int64(n)), codeLen)
	return strings.ToLower(code)
}

// ShortCode 根据ID生成短ID
func ShortCode(id int64, codeLen int) string {
	if id == 0 {
		id = time.Now().UnixNano()
	}
	if codeLen == 0 {
		codeLen = 10
	}
	var binLen = int64(len(base))
	key := len(base)
	var code = make([]string, binLen)
	j := 1
	for id/binLen > 0 {
		index := id % binLen
		key--
		code[key] = base[index]
		id /= binLen
		j++
	}
	code[key-1] = base[id%binLen]
	if j < codeLen {
		code = append(code, suffix)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := 0; i < codeLen-1-j; i++ {
			code = append(code, base[r.Intn(int(binLen))])
		}
	}
	return strings.Join(code, "")
}
