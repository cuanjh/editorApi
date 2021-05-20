package repository

import (
	"context"
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Catalogs struct {
}

func (m *Catalogs) FindAll(ctx context.Context, params requests.Catalogs) (result []responses.Catalogs, err error) {
	collection := mgdb.MongoClient.Database(mgdb.DbEditor).Collection(editor.TbCatalogs)
	var filter = bson.D{}
	filter = append(filter, bson.E{"has_del", false})
	if !helpers.Empty(params.Code) {
		filter = append(filter, bson.E{"code", params.Code})
	}

	if !helpers.Empty(params.Parent_uuid) {
		filter = append(filter, bson.E{"parent_uuid", params.Parent_uuid})
	}

	cusor, err := collection.Find(ctx, filter, options.Find().SetSort(bson.M{"list_order": 1}))

	defer cusor.Close(ctx)
	err = cusor.All(ctx, &result)
	return
}
