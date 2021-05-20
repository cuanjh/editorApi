package router

import (
	"editorApi/controller/api"
	"editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		UserRouter.POST("remove", api.UserRemove)                 // 删除用户
		UserRouter.POST("changePassword", api.ChangePassword)     // 修改密码
		UserRouter.POST("uploadHeaderImg", api.UploadHeaderImg)   //上传头像
		UserRouter.POST("getUserList", api.GetUserList)           // 分页获取用户列表
		UserRouter.POST("setUserAuthority", api.SetUserAuthority) //设置用户权限
	}
}
