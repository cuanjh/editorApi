package repository

import (
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct{}

func (m *Users) TeacherList(ctx *gin.Context, params []primitive.ObjectID) (result []responses.UsersResponse, err error) {
	collection := mgdb.OnlineClient.Database(mgdb.DbKuyu).Collection(editor.TbUsers)

	var filter = bson.D{}
	if !helpers.Empty(params) {
		filter = append(filter, bson.E{"_id", bson.M{"$in": params}})
	}

	cusor, err := collection.Find(
		ctx,
		filter,
		nil,
	)
	defer cusor.Close(ctx)
	err = cusor.All(ctx, &result)
	return
}
