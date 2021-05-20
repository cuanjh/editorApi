package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitOnelineJobRouter(Router *gin.RouterGroup) {

	jobRouter := Router.Group("editor").Use(middleware.CORSMiddleware())
	{
		jobRouter.POST("online/job", editorapi.OnlineJob)
		jobRouter.POST("online/courseInfo", editorapi.OnlineCourseInfoJob)
		jobRouter.POST("online/list", editorapi.OnlineList)
		jobRouter.POST("online/del", editorapi.OnlineJobDel)
	}
}
