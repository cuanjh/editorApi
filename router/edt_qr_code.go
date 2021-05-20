package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"
	"github.com/gin-gonic/gin"
)

func InitQRcodeRouter(Router *gin.RouterGroup) {
	QRcodeRouterPublic := Router.Group("editor").Use(middleware.CORSMiddleware())
	{
		QRcodeRouterPublic.GET("qr_code/details", editorapi.QRcodeDetails)
	}
	QRcodeRouter := Router.Group("editor").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	{
		QRcodeRouter.POST("qr_code/list", editorapi.QRcodeList)
		QRcodeRouter.POST("qr_code/add", editorapi.QRcodeAdd)
		QRcodeRouter.POST("qr_code/delete", editorapi.QRcodeDelete)
		QRcodeRouter.POST("qr_code/update", editorapi.QRcodeUpdate)
		QRcodeRouter.POST("qr_code/image", editorapi.QRcodeImage)
	}
}
