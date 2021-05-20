package editorapi

import (
	"editorApi/controller/servers"

	"github.com/gin-gonic/gin"
)

var (
	COURSE_ASSETS_DOMAIN     string = "https://course-assets1.talkmate.com"
	COURSE_ASSETS_BUCKET     string = "assets"
	COURSE_UPLOADFILE_BUCKET string = "uploadfiles"
	COURSE_UPLOADFILE_DOMAIN string = "https://uploadfile1.talkmate.com"
	EDITOR_DB                string = "editor"
)

// @Tags EditorInfoAPI（公共信息接口）
// @Summary 获取上传课程资料的七牛Token  , bucket:assets
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/info/token [get]
func QiniuToken(c *gin.Context) {
	servers.ReportFormat(c, true, "课程资料上传Token", gin.H{
		"token": servers.UploadToken(COURSE_ASSETS_BUCKET),
	})
}

// @Tags EditorInfoAPI（公共信息接口）
// @Summary 获取上传课程资料的七牛Token , bucket:uploadfile
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/info/token/uploadfile [get]
func QiniuUploadFileToken(c *gin.Context) {
	servers.ReportFormat(c, true, "课程资料上传Token", gin.H{
		"token": servers.UploadToken(COURSE_UPLOADFILE_BUCKET),
	})
}

// @Tags EditorInfoAPI（公共信息接口）
// @Summary 获取配置信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/info/config [get]
func ConfigInfo(c *gin.Context) {
	servers.ReportFormat(c, true, "配置信息", gin.H{
		"langInfos": []map[string]string{
			map[string]string{
				"langKey": "zh-CN",
				"name":    "中文",
			},
			map[string]string{
				"langKey": "en",
				"name":    "英文",
			},
		},
		"assetsDomain":     COURSE_ASSETS_DOMAIN,
		"uploadfileDomain": COURSE_UPLOADFILE_DOMAIN,
	})
}
