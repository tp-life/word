package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ErrorHandler 错误中转
func ErrorHandler(ctx *gin.Context, errCode int, errMsg string, hint ...interface{}) {
	if errCode != Success {
		httpStatus := http.StatusOK
		NewResponse(errCode, hint, errMsg).End(ctx, httpStatus)
		ctx.Abort()
	}
}

//limit 操作
func Limit(ctx *gin.Context) int {
	limit := ctx.DefaultQuery("size", "10")
	num, _ := strconv.Atoi(limit)
	return num
}

//offset 操作
func Offset(ctx *gin.Context) int {
	offset := ctx.DefaultQuery("pages", "1")
	page, _ := strconv.Atoi(offset)
	if page < 1 {
		page = 1
	}
	return (page - 1) * Limit(ctx)
}
