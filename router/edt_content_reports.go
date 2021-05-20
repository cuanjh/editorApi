package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"
	"github.com/gin-gonic/gin"
)

func InitContentReportsRouter(Router *gin.RouterGroup) {
	ContentReportsRouter := Router.Group("editor").Use(middleware.CORSMiddleware(), middleware.JWTAuth())
	{
		ContentReportsRouter.POST("content_reports/create", editorapi.ContentReportsCreate)
		ContentReportsRouter.POST("content_reports/find", editorapi.ContentReportsFind)
		ContentReportsRouter.POST("content_reports/findone", editorapi.ContentReportsFindOne)
		ContentReportsRouter.POST("content_reports/list", editorapi.ContentReportsList)
		ContentReportsRouter.POST("content_reports/update", editorapi.ContentReportsUpdate)
		ContentReportsRouter.POST("content_reports/delete", editorapi.ContentReportsDelete)
	}
}

