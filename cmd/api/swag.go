//+build doc

package api

import (
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func init() {
	swagHandler = ginSwagger.WrapHandler(swaggerFiles.Handler)
}
