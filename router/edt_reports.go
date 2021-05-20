package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"
	"github.com/gin-gonic/gin"
)

func InitReportsRouter(Router *gin.RouterGroup) {
	ReportsRouter := Router.Group("editor").Use(middleware.CORSMiddleware(), middleware.JWTAuth())
	{
		ReportsRouter.POST("reports/create", editorapi.ReportsCreate)
		ReportsRouter.POST("reports/find", editorapi.ReportsFind)
		ReportsRouter.POST("reports/findone", editorapi.ReportsFindOne)
		ReportsRouter.POST("reports/list", editorapi.ReportsList)
		ReportsRouter.POST("reports/update", editorapi.ReportsUpdate)
		ReportsRouter.POST("reports/delete", editorapi.ReportsDelete)
	}
}
