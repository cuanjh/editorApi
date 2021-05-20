package router


import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"
	"github.com/gin-gonic/gin"
)

func InitTeacherRouter(Router *gin.RouterGroup) {

	LangRouter := Router.Group("editor").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	{
		LangRouter.POST("teacher/list", editorapi.TeacherList)
		LangRouter.POST("teacher/details", editorapi.TeacherDetails)
		LangRouter.POST("teacher/audit", editorapi.Audit)
	}
}