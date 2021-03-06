package router

import (
	"editorApi/controller/api"
	"editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitApiRouter(Router *gin.RouterGroup) {
	ApiRouter := Router.Group("api").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		ApiRouter.POST("createApi", api.CreateApi)   //创建Api
		ApiRouter.POST("deleteApi", api.DeleteApi)   //删除Api
		ApiRouter.POST("getApiList", api.GetApiList) //获取Api列表
		ApiRouter.POST("getApiById", api.GetApiById) //获取单条Api消息
		ApiRouter.POST("updataApi", api.UpdataApi)   //更新api
		ApiRouter.POST("getAllApis", api.GetAllApis) // 获取所有api
	}
}
