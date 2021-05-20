package initRouter

import (
	_ "editorApi/docs"
	"editorApi/middleware"
	"editorApi/router"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.uber.org/zap"
	"net/http"
	"time"
)

//初始化总路由
func InitRouter() *gin.Engine {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowMethods("HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
	config.AddAllowHeaders("Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

	var Router = gin.Default()
	Router.Use(cors.New(config)) //解决跨域问题
	Router.Use(middleware.RequestIdMiddleware())
	//Router.Use(middleware.LoadTls())  // 打开就能玩https了
	Router.Use(middleware.Ginzap(zap.L(), time.RFC3339, true))
	Router.Use(middleware.RecoveryWithZap(zap.L(), true))
	Router.Use(middleware.ErrHandler()) // 全局panic错误
	Router.Use(middleware.Logger())     // 如果不需要日志 请关闭这里
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	Router.StaticFS("/editor/data", http.Dir("./data")) //设置静态文件目录

	ApiGroup := Router.Group("") // 方便统一添加路由组前缀 多服务器上线使用
	//Router.Use(middleware.Logger())
	router.InitUserRouter(ApiGroup)                  // 注册用户路由
	router.InitBaseRouter(ApiGroup)                  // 注册基础功能路由 不做鉴权
	router.InitMenuRouter(ApiGroup)                  // 注册menu路由
	router.InitAuthorityRouter(ApiGroup)             // 注册角色路由
	router.InitApiRouter(ApiGroup)                   // 注册功能api路由
	router.InitFileUploadAndDownloadRouter(ApiGroup) // 文件上传下载功能路由
	router.InitWorkflowRouter(ApiGroup)              // 工作流相关路由
	router.InitCasbinRouter(ApiGroup)                // 权限相关路由
	router.InitJwtRouter(ApiGroup)                   // jwt相关路由
	contentGroup := Router.Group("")                 // 方便统一添加路由组前缀 多服务器上线使用
	router.InitLangRouter(contentGroup)              //编辑器语言相关路由
	router.InitCourseRouter(contentGroup)            //编辑器语言相关路由
	router.InitContentVersionRouter(contentGroup)    //编辑器语言相关路由
	router.InitContentRouter(contentGroup)           //编辑器语言相关路由
	router.InitContentTypeRouter(contentGroup)       //编辑器语言相关路由
	router.InitModelRouter(contentGroup)             //编辑器语言相关路由
	router.InitTagRouter(contentGroup)               //编辑器语言相关路由
	router.InitLogsRouter(contentGroup)              //编辑器语言相关路由
	router.InitOnelineJobRouter(contentGroup)        //编辑器语言相关路由
	router.InitCatalogRouter(contentGroup)           //编辑器语言相关路由
	router.InitImageRouter(contentGroup)             //编辑器语言相关路由
	router.InitEditorAuthorityRouter(contentGroup)   //编辑器语言相关路由
	router.InitLiveRouter(contentGroup)              //编辑器语言相关路由
	router.InitDiscoverRouter(contentGroup)          //编辑器语言相关路由
	router.InitSearcWordRouter(contentGroup)         //编辑器语言相关路由
	router.InitClassWeixinRouter(contentGroup)       //编辑器语言相关路由
	router.InitOperateLogRouter(contentGroup)        //操作日志相关路由
	router.InitTeacherRouter(contentGroup)           //教师相关路由
	router.InitCourseFileRouter(contentGroup)        //教师相关路由
	router.InitQRcodeRouter(contentGroup)            //二维码相关路由
	router.InitDiceRouter(contentGroup)              //词典相关路由
	router.InitStatisticRouter(contentGroup)         //统计相关
	router.InitCardRouter(contentGroup)              //多语卡相关
	router.InitSentenceRouter(contentGroup)          //句子相关
	router.InitReportsRouter(contentGroup)           //审核相关
	router.InitActorsRouter(contentGroup)            //审核相关
	router.InitContentReportsRouter(contentGroup)    //审核相关
	router.InitDataRouter(contentGroup)              //审核相关
	return Router
}
