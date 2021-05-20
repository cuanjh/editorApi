package service

import (
	"editorApi/repository"
	"editorApi/requests"
	"editorApi/responses"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeacherService struct{}

// 列表数据
func (s *TeacherService) List(ctx *gin.Context, params requests.TeacherListRequests) (result []responses.TeacherResponse, err error) {
	var teacher = repository.Teacher{}
	result, err = teacher.TeacherList(ctx, params)

	if err != nil {
		return nil, err
	}

	if result != nil {
		var user_ids []primitive.ObjectID
		for _, item := range result {
			_id, _ := primitive.ObjectIDFromHex(item.UserId)
			user_ids = append(user_ids, _id)
		}

		var users = repository.Users{}
		userList, err := users.TeacherList(ctx, user_ids)
		if err != nil {
			return nil, err
		}
		// 数据处理
		for _, user := range userList {
			for key, teacher := range result {
				if user.Id.Hex() == teacher.UserId {
					result[key].Photo = user.Photo
					result[key].Nickname = user.Nickname
					result[key].Phonenumber = user.Phonenumber
					result[key].TalkmateId = user.TalkmateId
				}
			}
		}
	}
	return
}
