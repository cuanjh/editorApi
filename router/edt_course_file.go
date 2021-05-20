package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"
	"github.com/gin-gonic/gin"
)

func InitCourseFileRouter(Router *gin.RouterGroup) {

	LangRouter := Router.Group("editor").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	{
		LangRouter.POST("course_files/delete_file", editorapi.DeleteFile)
		LangRouter.POST("course_files/file_list", editorapi.FileList)
		LangRouter.POST("course_files/create_transcode", editorapi.CreateTranscode)
		LangRouter.POST("course_files/describe_transcode", editorapi.DescribeTranscode)
	}
}
