package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"

	// "editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitCourseRouter(Router *gin.RouterGroup) {

	LangRouter := Router.Group("editor").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	{
		LangRouter.POST("course/list", editorapi.CourseList)
		LangRouter.POST("course/add", editorapi.CourseAdd)
		LangRouter.POST("course/edit", editorapi.CourseEdit)
		LangRouter.POST("course/del", editorapi.CourseDel)
		LangRouter.GET("course/types", editorapi.CourseTypes)
		LangRouter.POST("course/detail", editorapi.CourseDetail)
		LangRouter.POST("course/detail/edit", editorapi.CourseDetailEdit)
		LangRouter.POST("course/unlock", editorapi.ContentUnlock)
	}
}
