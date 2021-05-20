package repository

import (
	"context"
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"editorApi/requests"
	"editorApi/tools/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type ContentExports struct {
	ID        string    `bson:"id" json:"id"`       // 下载数据ID
	UUID      string    `bson:"uuid" json:"uuid"`   // catalogs UUID
	Level     string    `bson:"level" json:"level"` //级别
	Name      string    `bson:"name" json:"name"`
	Code      string    `bson:"code" json:"code"`
	Url       string    `bson:"url" json:"url"`
	Status    int64     `bson:"status" json:"status"`         //1 代表正在处理，2；处理成功
	UserName  string    `bson:"user_name" json:"user_name"`   //操作人
	CreatedOn time.Time `bson:"created_on" json:"created_on"` //创建时间
}

func (m *ContentExports) AddContentExports(ctx context.Context, params requests.ContentExports) (inserted_id interface{}, err error) {
	collection := mgdb.MongoClient.Database(mgdb.DbEditor).Collection(editor.TbContentExports)

	m.ID = params.ID
	m.Code = params.Code
	m.Name = params.Name
	m.CreatedOn = time.Now()
	m.Status = 1
	insertOneResult, err := collection.InsertOne(ctx, m)
	inserted_id = insertOneResult.InsertedID
	return
}

func (m *ContentExports) UpdateContentExports(ctx context.Context, params requests.ContentExports) (upserted_id interface{}, err error) {
	collection := mgdb.MongoClient.Database(mgdb.DbEditor).Collection(editor.TbContentExports)

	var filter = bson.D{}
	if !helpers.Empty(params.ID) {
		filter = append(filter, bson.E{"id", params.ID})
	}

	var update = bson.M{
		"url":    params.Url,
		"status": 2,
	}

	updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return
	}
	upserted_id = updateResult.UpsertedID
	return
}
