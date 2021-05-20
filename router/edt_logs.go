package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"

	// "editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitLogsRouter(Router *gin.RouterGroup) {

	logsRouter := Router.Group("editor").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	{
		logsRouter.POST("logs/list", editorapi.EditorLogs)
	}
}
