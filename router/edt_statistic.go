package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"
	"github.com/gin-gonic/gin"
)

func InitStatisticRouter(Router *gin.RouterGroup) {

	StatisticRouter := Router.Group("editor").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	{
		StatisticRouter.POST("statistic/unlock_chapter", editorapi.StatisticUnlockChapter)
		StatisticRouter.POST("statistic/unlock_part", editorapi.StatisticUnlockPart)
	}
}
