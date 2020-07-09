package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ErrorHandler 错误中转
func ErrorHandler(ctx *gin.Context, errCode int, errMsg string, hint ...interface{}) {
	if errCode != Success {
		httpStatus := http.StatusOK
		NewResponse(errCode, hint, errMsg).End(ctx, httpStatus)
		ctx.Abort()
	}
}
