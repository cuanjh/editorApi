package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"

	// "editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitLiveRouter(Router *gin.RouterGroup) {

	LangRouter := Router.Group("live").Use(middleware.CORSMiddleware())
	{
		LangRouter.POST("shareinfo", editorapi.ShareInfo)
		LangRouter.POST("list", editorapi.LiveList)
		LangRouter.POST("edit_list_order", editorapi.EditListOrder)
		LangRouter.POST("add", editorapi.LiveAdd)
		LangRouter.POST("edit", editorapi.LiveEdit)
		LangRouter.POST("del", editorapi.LiveDel)
		LangRouter.POST("online", editorapi.LiveOnline)
		LangRouter.POST("offline", editorapi.LiveOffline)
		LangRouter.POST("course/online", editorapi.LiveCourseOnline)
		LangRouter.POST("course/offline", editorapi.LiveCourseOffline)
		LangRouter.POST("course/edit", editorapi.LiveCourseEdit)
		LangRouter.POST("chatroom/sendmsg", editorapi.ChatroomSendMsg)
		LangRouter.POST("chatroom/comments", editorapi.ChatroomComments)
		LangRouter.GET("chatroom/majia", editorapi.ChatroomMajia)
		LangRouter.GET("wxtoken", editorapi.WXToken)
		LangRouter.POST("pushurl", editorapi.LivePushUrl)

		LangRouter.POST("sub", editorapi.LiveSub)
		LangRouter.POST("usercount", editorapi.UserCount)

		LangRouter.POST("comments_upload", editorapi.CommentsUpload) //上传评论数据
		LangRouter.POST("send_live_common", editorapi.SendLiveCommon) //发布评论数据

		LangRouter.POST("gagadd", editorapi.GagAdd)
		LangRouter.POST("gagremove", editorapi.GagRemove)
		LangRouter.POST("dbtest", editorapi.DbTest)
	}
}
