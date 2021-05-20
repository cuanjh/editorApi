package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"
	"github.com/gin-gonic/gin"
)

//ImportCart
func InitCardRouter(Router *gin.RouterGroup) {
	CardRouter := Router.Group("editor").Use(middleware.CORSMiddleware())
	{
		CardRouter.POST("card/import", editorapi.ImportCard)
		CardRouter.POST("card/delete_sentence", editorapi.DeleteSentence)
		CardRouter.POST("card/read_files", editorapi.ReadFiles)
		CardRouter.POST("card/read_soundinfos", editorapi.ReadSoundInfos)
		CardRouter.POST("card/read_sentence", editorapi.ReadSentence)
		CardRouter.POST("card/read_words", editorapi.ReadWords)
		CardRouter.POST("card/add_tag", editorapi.AddTag)
		CardRouter.POST("card/add_sound", editorapi.AddSound)
		CardRouter.POST("card/aeneas_job", editorapi.AeneasJob)
	}
}
