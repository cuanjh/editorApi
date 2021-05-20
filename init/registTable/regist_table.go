package registTable

import (
	"editorApi/model/dbModel"
	"editorApi/model/sysModel"
	"github.com/jinzhu/gorm"
)

//注册数据库表专用
func RegistTable(db *gorm.DB) {
	db.AutoMigrate(sysModel.SysUser{},
		sysModel.SysAuthority{},
		sysModel.SysMenu{},
		sysModel.SysApi{},
		sysModel.SysBaseMenu{},
		sysModel.JwtBlacklist{},
		dbModel.ExaFileUploadAndDownload{},
		sysModel.SysWorkflow{},
		sysModel.SysWorkflowStepInfo{},
	)
}
