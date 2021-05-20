package editorapi

import (
	"context"
	"editorApi/controller/servers"
	"editorApi/init/initNats"
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"editorApi/middleware"
	"editorApi/tools/helpers"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
	uuid "github.com/satori/go.uuid"
	"github.com/tealeg/xlsx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strconv"
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

type ContentExportParams struct {
	ID          string   `json:"id"`    // 下载数据ID
	UUID        string   `json:"uuid"`  // UUID
	Level       string   `json:"level"` //级别
	ParentUuids []string `json:"parent_uuids"`
	AttrTag     string   `json:"attr_tag"`
	Name        string   `json:"name"`
	Code        string   `json:"code"`
	Merge       bool     `json:"merge"`
}

var tblContentExports string = "content_exports"
var exportCatalogsName map[int]string

type ContentExportListParams struct {
	PageNo   int64  `json:"pageNo"`
	PageSize int64  `json:"pageSize"`
	Code     string `json:"code"`
}

// @Tags EditorContentAPI（内容接口）
// @Summary 导出课程内容列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param file query string true "文件名称"
// @Param data body ContentExportListParams true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/content/export_list [post]
func ContentExportList(ctx *gin.Context) {
	param := ContentExportListParams{}
	ctx.BindJSON(&param)
	var limit int64
	var skip int64

	limit = param.PageSize
	if helpers.Empty(limit) {
		limit = 40
	}
	skip = param.PageNo * limit

	if helpers.Empty(param.Code) {
		servers.ReportFormat(ctx, false, "课程编码不能为空！", gin.H{})
		return
	}

	contentExports := []*ContentExports{}

	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentExports)
	cusor, err := collection.Find(
		ctx,
		bson.M{
			"code": param.Code,
		},
		options.Find().SetSort(bson.M{"created_on": -1}).SetLimit(limit).SetSkip(skip),
	)
	if !helpers.Empty(err) {
		checkErr(ctx, err)
		return
	}
	defer cusor.Close(ctx)
	cusor.All(ctx, &contentExports)
	servers.ReportFormat(ctx, true, "下载列表", gin.H{
		"data": contentExports,
	})
	return
}

// @Tags EditorContentAPI（内容接口）
// @Summary 导出课程内容
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param file query string true "文件名称"
// @Param data body ContentExportParams true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/content/export [post]
func ContentExport(ctx *gin.Context) {
	var param ContentExportParams
	ctx.BindJSON(&param)
	token := ctx.Request.Header.Get("x-token")
	j := middleware.NewJWT()
	claims, _ := j.ParseToken(token)

	if helpers.Empty(param.UUID) || helpers.Empty(param.Level) {
		servers.ReportFormat(ctx, false, "必填参数不能为空！", gin.H{})
		return
	}

	//HanderContentExport(&param)
	//ctx.Header("Content-Type", "application/octet-stream")
	//ctx.Header("Content-Disposition", "attachment; filename="+time.Now().Format("2006-01-02")+".xlsx")
	//ctx.Header("Content-Transfer-Encoding", "binary")
	//回写到web 流媒体 形成下载
	//_ = file.Write(ctx.Writer)

	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentExports)

	id := uuid.NewV4().String()
	var contentExports ContentExports
	contentExports.ID = id
	contentExports.CreatedOn = time.Now()
	contentExports.UUID = param.UUID
	contentExports.Level = param.Level
	contentExports.Code = param.Code
	contentExports.Status = 1
	contentExports.UserName = claims.NickName

	collection.InsertOne(ctx, contentExports)

	// 异步操作
	initNats.NatsConn.Publish("ExportCourseContent",
		&ContentExportParams{
			UUID:  param.UUID,
			Level: param.Level,
			ID:    id,
			Merge: param.Merge,
		},
	)

	servers.ReportFormat(ctx, true, "操作成功！", nil)
}

// 处理生成任务
func HanderContentExport(param *ContentExportParams) {
	// 生成一个新的文件
	file := xlsx.NewFile()
	var SheetMap = make(map[string][]string)
	var jobName string
	if param.Level == "1" {
		catalogs1 := getCatalogs(param)
		var param1 ContentExportParams
		// 1 层
		for _, catalog := range catalogs1 {
			param1.ParentUuids = append(param1.ParentUuids, catalog.Uuid)
			jobName = catalog.Name
		}

		// 2 层
		catalogs2 := getCatalogs(&param1)

		for k, catalog := range catalogs2 {
			k = k + 1
			name := strconv.Itoa(k)
			if k < 10 {
				name = "0" + name
			}
			SheetMap[name+"-"+catalog.Name] = nil
		}

		for key, _ := range SheetMap {
			var tmp []string
			for k, catalog := range catalogs2 {
				k = k + 1
				name := strconv.Itoa(k)
				if k < 10 {
					name = "0" + name
				}
				if key == name+"-"+catalog.Name {
					tmp = append(tmp, catalog.Uuid)
				}
			}
			SheetMap[key] = tmp
		}
	} else {
		// 2 层
		catalogs2 := getCatalogs(param)
		for k, catalog := range catalogs2 {
			k = k + 1
			name := strconv.Itoa(k)
			if k < 10 {
				name = "0" + name
			}
			SheetMap[name+"-"+catalog.Name] = nil
			jobName = catalog.Name
			var paramTmp ContentExportParams
			paramTmp.UUID = catalog.Parent_uuid
			catalogsTmps := getCatalogs(&paramTmp)
			for _, catalogsTmp := range catalogsTmps {
				jobName = catalogsTmp.Name + "-" + jobName
			}
		}

		for key, _ := range SheetMap {
			var tmp []string
			for k, catalog := range catalogs2 {
				k = k + 1
				name := strconv.Itoa(k)
				if k < 10 {
					name = "0" + name
				}
				if key == name+"-"+catalog.Name {
					tmp = append(tmp, catalog.Uuid)
				}
			}
			SheetMap[key] = tmp
		}
	}

	// 添加sheet页
	if param.Merge == true {
		titleList := []string{"catalog", "content", "slideLayout", "uuid", "model"}
		sheet, _ := file.AddSheet("All")
		// 插入表头
		titleRow := sheet.AddRow()
		for _, v := range titleList {
			cell := titleRow.AddCell()
			cell.Value = v
			//表头字体颜色
			cell.GetStyle().Font.Color = "00FF0000"
			//居中显示
			cell.GetStyle().Alignment.Horizontal = "center"
			cell.GetStyle().Alignment.Vertical = "center"
		}

		var contents []resultContent
		helpers.EachMap(SheetMap, func(sheetName string, item []string) {
			var param2 ContentExportParams
			param2.ParentUuids = item

			// 3 层
			catalogs3 := getCatalogs(&param2)
			for _, catalog3 := range catalogs3 {
				if !helpers.Empty(catalog3.Content_model) {
					//todo
					var contentParam content
					contentParam.ParentUuid = catalog3.Uuid
					contentParam.ContentModel = catalog3.Content_model
					data := getContent(contentParam)
					for _, itme := range data {
						itme.ContentModel = contentParam.ContentModel
						contents = append(contents, itme)
					}
				} else {
					catalogs4 := getExportCatalogsByParentUuid(catalog3.Uuid)
					for _, catalog4 := range catalogs4 {
						if !helpers.Empty(catalog4.Content_model) {
							//todo
							var contentParam content
							contentParam.ParentUuid = catalog4.Uuid
							contentParam.ContentModel = catalog4.Content_model
							data := getContent(contentParam)
							for _, itme := range data {
								itme.ContentModel = contentParam.ContentModel
								contents = append(contents, itme)
							}
						} else {
							catalogs5 := getExportCatalogsByParentUuid(catalog4.Uuid)
							for _, catalog5 := range catalogs5 {
								var contentParam content
								contentParam.ParentUuid = catalog5.Uuid
								contentParam.ContentModel = catalog5.Content_model
								data := getContent(contentParam)
								for _, itme := range data {
									itme.ContentModel = contentParam.ContentModel
									contents = append(contents, itme)
								}
							}
						}
					}
				}
			}
		})

		for key, itme := range contents {
			row := sheet.AddRow()
			tmpFormNo := row.AddCell()
			tmpFormNo.SetValue(getCatalogsName(itme.ParentUuid) + "->" + strconv.Itoa(key+1))

			tmpSentence := row.AddCell()
			tmpSentence.SetValue(itme.Sentence)

			tmpType := row.AddCell()
			tmpType.SetValue(getTypeContent(itme.Type))

			tmpUuid := row.AddCell()
			tmpUuid.SetValue(itme.Uuid)

			tmpModel := row.AddCell()
			tmpModel.SetValue(itme.ContentModel)
		}
	} else {
		helpers.EachMap(SheetMap, func(sheetName string, item []string) {
			var param2 ContentExportParams
			param2.ParentUuids = item

			titleList := []string{"catalog", "content", "slideLayout", "uuid", "model"}

			sheet, _ := file.AddSheet(sheetName)

			// 插入表头
			titleRow := sheet.AddRow()
			for _, v := range titleList {
				cell := titleRow.AddCell()
				cell.Value = v
				//表头字体颜色
				cell.GetStyle().Font.Color = "00FF0000"
				//居中显示
				cell.GetStyle().Alignment.Horizontal = "center"
				cell.GetStyle().Alignment.Vertical = "center"
			}

			// 3 层
			catalogs3 := getCatalogs(&param2)
			for _, catalog3 := range catalogs3 {
				if !helpers.Empty(catalog3.Content_model) {
					//todo
					var contentParam content
					contentParam.ParentUuid = catalog3.Uuid
					contentParam.ContentModel = catalog3.Content_model
					data := getContent(contentParam)
					for key, itme := range data {
						row := sheet.AddRow()
						tmpFormNo := row.AddCell()
						tmpFormNo.SetValue(getCatalogsName(itme.ParentUuid) + "->" + strconv.Itoa(key+1))

						tmpSentence := row.AddCell()
						tmpSentence.SetValue(itme.Sentence)

						tmpType := row.AddCell()
						tmpType.SetValue(getTypeContent(itme.Type))

						tmpUuid := row.AddCell()
						tmpUuid.SetValue(itme.Uuid)

						tmpModel := row.AddCell()
						tmpModel.SetValue(contentParam.ContentModel)
					}
				} else {
					catalogs4 := getExportCatalogsByParentUuid(catalog3.Uuid)
					for _, catalog4 := range catalogs4 {
						if !helpers.Empty(catalog4.Content_model) {
							//todo
							var contentParam content
							contentParam.ParentUuid = catalog4.Uuid
							contentParam.ContentModel = catalog4.Content_model
							data := getContent(contentParam)
							for key, itme := range data {
								row := sheet.AddRow()
								tmpFormNo := row.AddCell()
								tmpFormNo.SetValue(getCatalogsName(itme.ParentUuid) + "->" + strconv.Itoa(key+1))

								tmpSentence := row.AddCell()
								tmpSentence.SetValue(itme.Sentence)

								tmpType := row.AddCell()
								tmpType.SetValue(getTypeContent(itme.Type))

								tmpUuid := row.AddCell()
								tmpUuid.SetValue(itme.Uuid)

								tmpModel := row.AddCell()
								tmpModel.SetValue(contentParam.ContentModel)
							}
						} else {
							catalogs5 := getExportCatalogsByParentUuid(catalog4.Uuid)
							for _, catalog5 := range catalogs5 {
								var contentParam content
								contentParam.ParentUuid = catalog5.Uuid
								contentParam.ContentModel = catalog5.Content_model
								data := getContent(contentParam)
								for key, itme := range data {
									row := sheet.AddRow()
									tmpFormNo := row.AddCell()
									tmpFormNo.SetValue(getCatalogsName(itme.ParentUuid) + "->" + strconv.Itoa(key+1))

									tmpSentence := row.AddCell()
									tmpSentence.SetValue(itme.Sentence)

									tmpType := row.AddCell()
									tmpType.SetValue(getTypeContent(itme.Type))

									tmpUuid := row.AddCell()
									tmpUuid.SetValue(itme.Uuid)

									tmpModel := row.AddCell()
									tmpModel.SetValue(contentParam.ContentModel)
								}
							}
						}
					}
				}
			}
		})
	}

	os.MkdirAll("data/exports/", os.ModePerm)
	filename := fmt.Sprintf("data/exports/%v", uuid.NewV4().String()+"-"+time.Now().Format("2006-01-02")+".xlsx")
	file.Save(filename)

	// 数据修改
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentExports)
	collection.UpdateOne(
		context.TODO(),
		bson.M{
			"id": param.ID,
		},
		bson.M{"$set": bson.M{
			"url":    "/editor/" + filename,
			"name":   jobName,
			"status": 2,
		}},
	)
}

type content struct {
	ContentModel string `json:"content_model"`
	ParentUuid   string `json:"parent_uuid"`
}

type resultContent struct {
	Sentence     string `bson:"sentence" json:"sentence"`
	Type         string `bson:"type" json:"type"`
	ParentUuid   string `bson:"parent_uuid" json:"parent_uuid"`
	Uuid         string `bson:"uuid" json:"uuid"`
	ContentModel string `bson:"content_model" json:"content_model"`
}

func getContent(param content) (contents []resultContent) {
	ctx := context.TODO()
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(param.ContentModel)
	var filter = bson.D{}

	filter = append(filter, bson.E{"has_del", false})

	if param.ParentUuid != "" {
		filter = append(filter, bson.E{"parent_uuid", param.ParentUuid})
	}

	cusor, _ := collection.Find(
		ctx,
		filter,
		options.Find().SetSort(bson.M{
			"list_order": 1,
		}).SetProjection(bson.M{
			"_id":         0,
			"sentence":    1,
			"type":        1,
			"parent_uuid": 1,
			"uuid":        1,
		}).SetLimit(200),
	)
	defer cusor.Close(ctx)

	cusor.All(ctx, &contents)
	return
}

func getExportCatalogsByParentUuid(parent_uuid string) (catalogs []*editor.Catalogs) {
	ctx := context.TODO()
	catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)
	cusor, _ := catalogsCollection.Find(
		ctx,
		bson.M{
			"parent_uuid": parent_uuid,
			"has_del":     false,
			"is_show":     true,
		},
		options.Find().SetSort(bson.M{"list_order": 1}),
	)
	defer cusor.Close(ctx)
	cusor.All(ctx, &catalogs)
	return
}

func getCatalogs(param *ContentExportParams) (catalogs []*editor.Catalogs) {
	ctx := context.TODO()
	catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)
	var (
		cusor *mongo.Cursor
	)

	var filter = bson.D{}

	filter = append(filter, bson.E{"has_del", false})

	if param.UUID != "" {
		filter = append(filter, bson.E{"uuid", param.UUID})
	}

	if param.AttrTag != "" {
		filter = append(filter, bson.E{"attr_tag", param.AttrTag})
	}

	if param.ParentUuids != nil {
		filter = append(filter, bson.E{"parent_uuid", bson.M{"$in": param.ParentUuids}})
	}

	if param.Name != "" {
		filter = append(filter, bson.E{"name", param.Name})
	}

	cusor, _ = catalogsCollection.Find(
		ctx,
		filter,
		options.Find().SetSort(bson.M{"list_order": 1}).SetProjection(bson.M{
			"_id":     0,
			"has_del": 0,
		}),
	)

	defer cusor.Close(ctx)
	cusor.All(ctx, &catalogs)
	return
}

func getCatalogsName(parent_uuid string) (str string) {
	exportCatalogsName = make(map[int]string)
	getAllCatalogsName(parent_uuid, 100)
	helpers.EachMap(exportCatalogsName, func(key int, item string) {
		if str != "" {
			str = str + "->" + item
		} else {
			str = item
		}
	})
	return
}

func getAllCatalogsName(parent_uuid string, i int) {
	var catalog editor.Catalogs
	mgdb.FindOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tblCatalogs,
		bson.M{
			"uuid":    parent_uuid,
			"has_del": false,
		},
		nil,
		&catalog,
	)

	if !helpers.Empty(catalog) && !helpers.Empty(catalog.Uuid) {
		exportCatalogsName[i] = catalog.Name
		i = i - 1
		getAllCatalogsName(catalog.Parent_uuid, i)
	}
	return
}
