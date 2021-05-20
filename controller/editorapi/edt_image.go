package editorapi

import (
	"editorApi/controller/servers"
	"editorApi/init/mgdb"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var tblImageSys string = "image_sys"
var tblImageSysTags string = "image_sys_tags"

type imageDownloadParam struct {
	TagKey string `json:"tagKey"`
}

// @Tags EditorImageAPI （图片系统接口）
// @Summary 图片信息下载
// @Security ApiKeyAuth
// @accept mpfd
// @Produce application/json
// @Param tagKey formData string true "图片标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/image/download [get]
func ImageDownload(ctx *gin.Context) {
	tagKey := ctx.Query("tagKey")

	var images []*struct {
		ID       primitive.ObjectID `bson:"_id" json:"_id"`
		ImageID  string             `json:"image_id"`
		ImageUrl string             `bson:"image_url" json:"image_url"`
		Desc     string             `bson:"desc" json:"desc"`
		TagKeys  []string           `bson:"tagKeys" json:"tagKeys"`
	}

	excelFile := xlsx.NewFile()
	sheet, _ := excelFile.AddSheet(tagKey)
	header := sheet.AddRow()
	descH := header.AddCell()
	descH.Value = "图片描述"
	urlH := header.AddCell()
	urlH.Value = "图片URL地址"

	var offset int64 = 0
	var pageSize int64 = 500

	where := bson.M{
		"tagKeys": tagKey,
	}
	for {
		mgdb.Find(
			mgdb.EnvEditor,
			mgdb.DbEditor,
			tblImageSys,
			where,
			map[string]int{
				"desc": 1,
			},
			nil,
			offset,
			pageSize,
			&images,
		)
		for _, img := range images {
			row := sheet.AddRow()
			descH := row.AddCell()
			descH.Value = img.Desc
			urlH := row.AddCell()
			urlH.Value = img.ImageUrl
		}

		offset += pageSize

		if offset == 10000 && len(images) < int(pageSize) {
			break
		}
	}

	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+time.Now().Format("2006-01-02")+".xlsx")
	ctx.Header("Content-Transfer-Encoding", "binary")
	excelFile.Write(ctx.Writer)
}

type imageSearchParam struct {
	SearchType int    `json:"searchType"` // 0模糊搜索，1精确搜索
	TagKey     string `json:"tagKey"`
	Words      string `json:"words"`
	PageSize   int64  `json:"pageSize"`
	Page       int64  `json:"page"`
}

// @Tags EditorImageAPI （图片系统接口）
// @Summary 图片搜索
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.imageSearchParam true "图片搜索参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/image/search [post]
func ImageSearch(ctx *gin.Context) {
	var param imageSearchParam
	ctx.BindJSON(&param)
	if param.PageSize == 0 {
		param.PageSize = 100
	}
	if param.Page == 0 {
		param.Page = 1
	}
	offset := (param.Page - 1) * param.PageSize

	var images []*struct {
		ID       primitive.ObjectID `bson:"_id" json:"_id"`
		ImageID  string             `json:"image_id"`
		ImageUrl string             `bson:"image_url" json:"image_url"`
		Desc     string             `bson:"desc" json:"desc"`
		TagKeys  []string           `bson:"tagKeys" json:"tagKeys"`
	}
	where := bson.M{}

	if param.SearchType == 0 {
		if param.Words != "" {
			where["desc"] = primitive.Regex{
				Pattern: param.Words,
				Options: "i",
			}
		}
		if param.TagKey != "" {
			where["tagKeys"] = param.TagKey
		}
	} else {
		if param.Words != "" {
			where["desc"] = param.Words
		}
		if param.TagKey != "" {
			where["tagKeys"] = param.TagKey
		}
		//where["desc"] = param.Words
	}

	mgdb.Find(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tblImageSys,
		where,
		nil,
		nil,
		offset,
		param.PageSize,
		&images,
	)
	for k, _ := range images {
		images[k].ImageID = images[k].ID.Hex()
	}
	servers.ReportFormat(ctx, true, "图片列表", gin.H{
		"images": images,
	})
}

type imageAddParam struct {
	TagKeys  []string `json:"tagKeys"`
	ImageUrl string   `json:"image_url"`
}

// @Tags EditorImageAPI （图片系统接口）
// @Summary 图片添加
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.imageAddParam true "图片添加"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/image/add [post]
func ImageAdd(ctx *gin.Context) {
	var param imageAddParam
	ctx.BindJSON(&param)

	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblImageSys)
	collection.InsertOne(ctx, bson.M{
		"desc":      strings.Join(param.TagKeys, "/"),
		"image_url": param.ImageUrl,
		"tagKeys":   param.TagKeys,
	})

	//更新图片tag
	mgdb.UpdateMany(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tblImageSysTags,
		bson.M{
			"tagKey": bson.M{"$in": param.TagKeys},
		},
		bson.M{
			"$set": bson.M{
				"conNum": bson.M{"$inc": 1},
				"del":    false,
			},
		},
		true,
	)
	servers.ReportFormat(ctx, true, "添加成功", gin.H{})
}

type imageDelParam struct {
	ImageID string `json:"image_id"`
}

// @Tags EditorImageAPI （图片系统接口）
// @Summary 图片删除
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.imageDelParam true "图片删除"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/image/del [post]
func ImageDel(ctx *gin.Context) {
	var param imageDelParam
	ctx.BindJSON(&param)
	_id, _ := primitive.ObjectIDFromHex(param.ImageID)
	mgdb.MongoClient.Database(EDITOR_DB).Collection(tblImageSys).DeleteOne(ctx, bson.M{
		"_id": _id,
	})
	servers.ReportFormat(ctx, true, "成功", gin.H{})
}

type imageAddMoreParam struct {
	TagKeys   []string `json:"tagKeys"`
	ImageUrls []string `json:"image_urls"`
	Names     []string `json:"names"`
}

// @Tags EditorImageAPI （图片系统接口）
// @Summary 图片添加
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.imageAddMoreParam true "图片添加"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/image/add/more [post]
func ImageAddMore(ctx *gin.Context) {
	var param imageAddMoreParam
	ctx.BindJSON(&param)

	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblImageSys)
	images := make([]interface{}, len(param.ImageUrls))
	for k, ig := range param.ImageUrls {
		images[k] = bson.M{
			"desc":      param.Names[k] + "/" + strings.Join(param.TagKeys, "/"),
			"image_url": ig,
			"tagKeys":   param.TagKeys,
		}
	}

	collection.InsertMany(ctx, images)

	//更新图片tag
	mgdb.UpdateMany(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tblImageSysTags,
		bson.M{
			"tagKey": bson.M{"$in": param.TagKeys},
		},
		bson.M{
			"$set": bson.M{
				"conNum": bson.M{"$inc": 1},
				"del":    false,
			},
		},
		true,
	)
	servers.ReportFormat(ctx, true, "添加成功", gin.H{})
}

type imageEditParam struct {
	ImageID  string   `json:"image_id"`
	TagKeys  []string `json:"tagKeys"`
	Desc     string   `json:"desc"`
	ImageUrl string   `json:"image_url"`
}

// @Tags EditorImageAPI （图片系统接口）
// @Summary 图片编辑
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.imageEditParam true "图片添加"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/image/edit [post]
func ImageEdit(ctx *gin.Context) {
	var param imageEditParam
	ctx.BindJSON(&param)

	_id, _ := primitive.ObjectIDFromHex(param.ImageID)
	mgdb.UpdateOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tblImageSys,
		bson.M{
			"_id": _id,
		},
		bson.M{
			"$set": bson.M{
				"desc":      param.Desc + "/" + strings.Join(param.TagKeys, "/"),
				"image_url": param.ImageUrl,
				"tagKeys":   param.TagKeys,
			},
		},
		false,
	)
	//更新图片tag
	mgdb.UpdateMany(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tblImageSysTags,
		bson.M{
			"tagKey": bson.M{"$in": param.TagKeys},
		},
		bson.M{
			"$set": bson.M{
				"conNum": bson.M{"$inc": 1},
				"del":    false,
			},
		},
		true,
	)
	servers.ReportFormat(ctx, true, "添加成功", gin.H{})
}

type imageTagParam struct {
	TagKey string `bson:"tagKey"`
}

// @Tags EditorImageAPI （图片系统接口）
// @Summary 图片标签列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/image/tags [post]
func ImageTags(ctx *gin.Context) {
	var tags []imageTagParam
	mgdb.Find(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tblImageSysTags,
		bson.M{
			"del": false,
		},
		nil,
		nil,
		0,
		500,
		&tags,
	)
	servers.ReportFormat(ctx, true, "图片标签列表", gin.H{
		"tags": tags,
	})
}

// @Tags EditorImageAPI （图片系统接口）
// @Summary 图片标签添加
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.imageTagParam true "图片标签添加"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/image/tag/add [post]
func ImageTagAdd(ctx *gin.Context) {
	var param imageTagParam
	ctx.BindJSON(&param)

	//更新图片tag
	mgdb.UpdateOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tblImageSysTags,
		bson.M{
			"tagKey": param.TagKey,
		},
		bson.M{
			"$set": bson.M{
				"del":       false,
				"createdOn": time.Now(),
			},
		},
		true,
	)
	servers.ReportFormat(ctx, true, "添加成功", gin.H{})
}

// @Tags EditorImageAPI （图片系统接口）
// @Summary 图片标签删除
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.imageTagParam true "图片标签添加"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/image/tag/del [post]
func ImageTagDel(ctx *gin.Context) {
	var param imageTagParam
	ctx.BindJSON(&param)

	//更新图片tag
	mgdb.UpdateOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tblImageSysTags,
		bson.M{
			"tagKey": param.TagKey,
		},
		bson.M{
			"$set": bson.M{
				"del": true,
			},
		},
		false,
	)
	servers.ReportFormat(ctx, true, "删除成功", gin.H{})
}
