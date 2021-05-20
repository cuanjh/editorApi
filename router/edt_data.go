package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"
	"github.com/gin-gonic/gin"
)

func InitDataRouter(Router *gin.RouterGroup) {
	DataRouter := Router.Group("editor").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	{
		DataRouter.POST("data/export", editorapi.ExportData)
		DataRouter.POST("data/course", editorapi.CourseData)
	}
}
