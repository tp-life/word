// Package app 提供全局公用依赖性极低操作
//  不要尝试写入复杂操作逻辑到这里, 可能会引起令人头疼的循环调用问题
//  其他包可以调用app, 但app不要调用其他包, 需要调用的在其他包中调用 Register 将服务注册
package app

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"word/docs"
	"word/pkg/i18n"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

// GinMode 这个变量由 Makefile 指定
var GinMode string

func init() {
	if !filepath.IsAbs(os.Args[0]) {
		os.Args[0], _ = filepath.Abs(os.Args[0])
	}
	gin.SetMode(GinMode)
	startLogger()
}

var (
	// 翻译
	translate = i18n.NewBundle(language.Chinese).LoadFiles(fmt.Sprintf("%s/locales", Root()), "yaml", yaml.Unmarshal)
)

// Response HTTP返回数据结构体, 可使用这个, 也可以自定义
type Response struct {
	Code    int         `json:"code"`    // 状态码,这个状态码是与前端和APP约定的状态码,非HTTP状态码
	Data    interface{} `json:"data"`    // 返回数据
	Message string      `json:"message"` // 自定义返回的消息内容
	msgData interface{} `json:"-"`       // 消息解析使用的数据
}

// SetMsgData 消息解析使用的数据
func (rsp *Response) MsgData(data interface{}) *Response {
	rsp.msgData = data
	return rsp
}

// End 在调用了这个方法之后,还是需要 return 的
func (rsp *Response) End(c *gin.Context, httpStatus ...int) {
	status := http.StatusOK
	if len(httpStatus) > 0 {
		status = httpStatus[0]
	}

	if rsp.Message != "" {
		rsp.Message = translate.NewPrinter(i18n.GetAcceptLanguages(c)...).Translate(rsp.Message, rsp.msgData)
	}
	c.JSON(status, rsp)
}

// Object 直接获得本对象
func (rsp *Response) Object(ctx *gin.Context) *Response {
	if rsp.Message != "" {
		rsp.Message = translate.NewPrinter(i18n.GetAcceptLanguages(ctx)...).Translate(rsp.Message, rsp.msgData)
	}
	return rsp
}

// NewResponse 接口返回统一使用这个
//  code 服务端与客户端和web端约定的自定义状态码
//  data 具体的返回数据
//  message 可不传,自定义消息
func NewResponse(code int, data interface{}, message ...string) *Response {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	return &Response{Code: code, Data: data, Message: msg}
}

func OK(c *gin.Context, data interface{}, message ...string) {
	NewResponse(Success, data, message...).End(c)
}

func OriOK(c *gin.Context, data interface{}, message ...string) {
	status := http.StatusOK
	NewResponse(Success, data, message...).Object(c)
	c.PureJSON(status, data)
}

func F(c *gin.Context, code int, message ...string) {
	NewResponse(code, nil, message...).End(c)
}

// Translator 翻译, 通过分析 AcceptLanguage 来获取用户接受的语言
func Translator(ctx *gin.Context) *i18n.Printer {
	return translate.NewPrinter(i18n.GetAcceptLanguages(ctx)...)
}

// Root 根目录
//  返回程序运行时的运行目录
func Root() string {
	binaryRootPath, _ := filepath.Abs(os.Args[0])
	if gin.Mode() == gin.ReleaseMode {
		return filepath.Dir(binaryRootPath)
	}
	if filepath.Base(filepath.Dir(binaryRootPath)) != docs.GetMod() {
		binaryRootPath, _ = filepath.Abs("./")
	} else {
		return filepath.Dir(binaryRootPath)
	}

	//TODO 下面增加几行用于方便本地测试
	exit := pathExists(filepath.Join(binaryRootPath, "configs"))
	if !exit {
		binaryRootPath, _ = filepath.Abs("./configs")
		for filepath.Base(binaryRootPath) != docs.GetMod() {
			binaryRootPath = filepath.Dir(binaryRootPath)
		}
	}
	return binaryRootPath
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

// Name 程序名
//  返回程序名称
func Name() string {
	stat, _ := os.Stat(os.Args[0])
	return stat.Name()
}

// Md5 md5 hash
func Md5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// GetUserAgent 得到用户ua
func GetUserAgent(c *gin.Context) string {
	ua := c.Request.Header.Get("User-Agent")
	if ok, _ := regexp.MatchString("Android|android", ua); ok {
		return "android"
	} else if ok, _ := regexp.MatchString("iPad|iPhone", ua); ok {
		return "ios"
	}
	return "pc"
}
