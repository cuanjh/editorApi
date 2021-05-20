package repository

import (
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"editorApi/middleware"
	"editorApi/requests"
	"github.com/gin-gonic/gin"
	"time"
)

type TeacherAudit struct {
}

/**
老师审核
*/
func (m *TeacherAudit) AuditTeacher(ctx *gin.Context, params requests.TeacherAuditRequests) (inserted_id interface{}, err error) {
	collection := mgdb.OnlineClient.Database(mgdb.DbKuyu).Collection(editor.TbTeacherAudit)

	claims, _ := middleware.GetClaims(ctx)

	var data = make(map[string]interface{})
	data["audit_time"] = time.Now().Format("2006-01-02 15:04:05")
	data["created_on"] = time.Now().Format("2006-01-02 15:04:05")
	data["user_id"] = params.UserId
	data["status"] = params.Status
	data["content"] = params.Content
	data["auditor"] = claims.UUID.String()
	data["auditor_nickname"] = claims.NickName

	insertOneResult, err := collection.InsertOne(ctx, data)
	if err != nil {
		return
	}

	// 修改审核记录
	var teacher = Teacher{}
	var teacherRequests = requests.TeacherRequests{}
	teacherRequests.UserId = params.UserId
	teacherRequests.Status = params.Status
	teacher.Audit(ctx, teacherRequests)

	inserted_id = insertOneResult.InsertedID
	return
}
