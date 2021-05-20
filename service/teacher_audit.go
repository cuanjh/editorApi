package service

import (
	"editorApi/repository"
	"editorApi/requests"
	"github.com/gin-gonic/gin"
)

type TeacherAuditService struct{}

func (s *TeacherAuditService) AuditTeacher(ctx *gin.Context, params requests.TeacherAuditRequests) (upserted_id interface{}, err error) {
	var teacher_audit = repository.TeacherAudit{}

	upserted_id, err = teacher_audit.AuditTeacher(ctx, params)
	return
}