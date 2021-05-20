package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"
	"github.com/gin-gonic/gin"
)

func InitOperateLogRouter(Router *gin.RouterGroup) {

	LangRouter := Router.Group("editor").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	{
		LangRouter.POST("operate_log/list", editorapi.OperateLogList)
		LangRouter.POST("operate_log/details", editorapi.OperateLogDetails)
		LangRouter.POST("operate_log/rollback", editorapi.Rollback)
	}
}