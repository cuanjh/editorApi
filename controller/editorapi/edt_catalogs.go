package editorapi

import (
	"context"
	"editorApi/controller/servers"
	"editorApi/init/mgdb"
	"editorApi/init/qmlog"
	"editorApi/mdbmodel/editor"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tblCatalogs string = "catalogs"

type catalogsListsParams struct {
	ParentUuid string `json:"parent_uuid"`
}

// @Tags EditorCatalogAPI(目录接口)
// @Summary 目录列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.catalogsListsParams true "目录列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"目录列表"}"
// @Router /editor/catalog/list [post]
func CatalogList(c *gin.Context) {

	var param catalogsListsParams
	c.BindJSON(&param)
	catalogs := []*editor.Catalogs{}
	catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)
	var (
		err   error
		cusor *mongo.Cursor
	)

	if cusor, err = catalogsCollection.Find(
		c,
		bson.M{
			"parent_uuid": param.ParentUuid,
			"has_del":     false,
		},
		options.Find().SetSort(bson.M{"list_order": 1}).SetProjection(bson.M{
			"_id":     0,
			"has_del": 0,
		}),
	); err != nil {
		checkErr(c, err)
		return
	}
	defer cusor.Close(c)
	cusor.All(c, &catalogs)

	servers.ReportFormat(c, true, "目录列表", gin.H{
		"catalogs": catalogs,
	})
}

type catalogsInfo struct {
	ParentUUID   string            `bson:"parent_uuid" json:"parent_uuid"`
	HasChanged   bool              `bson:"has_changed" json:"has_changed"`
	UpdateTime   int64             `bson:"update_time" json:"update_time"`
	Tags         []string          `bson:"tags" json:"tags"`
	AttrTag      string            `bson:"attr_tag" json:"attr_tag"`
	ListOrder    int               `bson:"list_order" json:"list_order"`
	IsShow       bool              `bson:"is_show" json:"is_show"`
	OnlineState  int8              `bson:"onlineState" json:"onlineState"`
	Title        map[string]string `bson:"title" json:"title"`
	Name         string            `bson:"name" json:"name"`
	HasDel       bool              `bson:"has_del" json:"has_del"`
	Flag         []string          `bson:"flag" json:"flag"`
	Desc         map[string]string `bson:"desc" json:"desc"`
	UUID         string            `bson:"uuid" json:"uuid"`
	Type         string            `bson:"type" json:"type"`
	Cover        []string          `bson:"cover" json:"cover"`
	ContentModel string            `bson:"content_model" json:"content_model"`
}

type catalogAddParam struct {
	Num         int          `json:"num"` //增加的目录或文件的数量
	CatalogInfo catalogsInfo `json:"catalogsInfo"`
}

// @Tags EditorCatalogAPI(目录接口)
// @Summary 添加目录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.catalogAddParam true "添加目录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/catalog/add [post]
func CatalogAdd(c *gin.Context) {
	defer rcv(c)
	var catalog catalogsInfo
	var para *catalogAddParam
	c.BindJSON(&para)

	catalogs := make([]interface{}, para.Num)

	catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)

	catalog = para.CatalogInfo
	parentUuids := []string{
		catalog.ParentUUID,
	}
	if para.Num == 0 {
		para.Num = 1
	}
	catalogUUIDS := []string{}
	if para.Num == 1 {
		tmp_uuid := uuid.NewV4().String()
		catalogUUIDS = append(catalogUUIDS, tmp_uuid)
		catalog.UUID = tmp_uuid
		catalog.HasChanged = true
		catalog.IsShow = true
		catalog.HasDel = false
		catalog.UpdateTime = time.Now().Unix()
		catalogs[0] = catalog
	} else {
		name := catalog.Name
		for i := 0; i < para.Num; i++ {
			tmp_uuid := uuid.NewV4().String()
			catalogUUIDS = append(catalogUUIDS, tmp_uuid)
			catalog.Name = name + strconv.Itoa(i+1)
			catalog.UUID = tmp_uuid
			catalog.HasChanged = true
			catalog.IsShow = true
			catalog.HasDel = false
			catalog.UpdateTime = time.Now().Unix()
			catalogs[i] = catalog
		}
	}
	_, e := catalogsCollection.InsertMany(c, catalogs)

	for _, item := range catalogUUIDS {
		//记录操作日志
		var (
			operation              OperateLogServer
			operateLogCreateParams = OperateLogParams{Ctx: c, Model: "catalogs", ParentUuid: item, Pattern: PATTERN_CATALOGS, Mold: MOLD_CATALOGS_CREATE}
		)
		operation.Start(operateLogCreateParams)
	}

	if e != nil {
		checkErr(c, e)
		return
	}
	//更新父目录信息
	updateParentCatalogInfo(c, parentUuids)
	servers.ReportFormat(c, true, "添加目录或文件成功", gin.H{
		"catalog": catalogs,
	})
}

type catalogEditInfo struct {
	HasChanged  bool              `bson:"has_changed" json:"has_changed"`
	UpdateTime  int64             `bson:"update_time" json:"update_time"`
	Tags        []string          `bson:"tags" json:"tags"`
	AttrTag     string            `bson:"attr_tag" json:"attr_tag"`
	ListOrder   int               `bson:"list_order" json:"list_order"`
	IsShow      bool              `bson:"is_show" json:"is_show"`
	OnlineState int8              `bson:"onlineState" json:"onlineState"`
	Title       map[string]string `bson:"title" json:"title"`
	GoalTitle   string            `bson:"goalTitle" json:"goalTitle"`
	Name        string            `bson:"name" json:"name"`
	Flag        []string          `bson:"flag" json:"flag"`
	Desc        map[string]string `bson:"desc" json:"desc"`
	Cover       []string          `bson:"cover" json:"cover"`
}

type catalogEditParam struct {
	Uuid        string          `json:"uuid"`
	CatalogInfo catalogEditInfo `json:"catalog_info"`
}

// @Tags EditorCatalogAPI(目录接口)
// @Summary 修改目录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.catalogEditParam true "修改目录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/catalog/edit [post]
func CatalogEdit(c *gin.Context) {

	var param catalogEditParam
	c.BindJSON(&param)

	//记录操作日志
	var (
		operation              OperateLogServer
		operateLogCreateParams = OperateLogParams{Ctx: c, Model: "catalogs", ParentUuid: param.Uuid, Pattern: PATTERN_CATALOGS, Mold: MOLD_CATALOGS_CREATE}
		operateLogDelParams    = OperateLogParams{Ctx: c, Model: "catalogs", ParentUuid: param.Uuid, Pattern: PATTERN_CATALOGS, Mold: MOLD_CATALOGS_EDIT}
	)
	operation.Start(operateLogCreateParams)
	defer operation.Finish(operateLogDelParams)

	catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)
	param.CatalogInfo.HasChanged = true
	param.CatalogInfo.UpdateTime = time.Now().Unix()

	var catalog *editor.Catalogs
	catalogsCollection.FindOneAndUpdate(c,
		bson.M{
			"uuid": param.Uuid,
		},
		bson.M{
			"$set": param.CatalogInfo,
		},
	).Decode(&catalog)
	//更新父目录信息
	updateParentCatalogInfo(c, []string{
		catalog.Parent_uuid,
	})

	servers.ReportFormat(c, true, "更新成功", gin.H{})
}

type catalogDelParam struct {
	Uuid string `json:"uuid"`
}

// @Tags EditorCatalogAPI(目录接口)
// @Summary 删除目录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.catalogDelParam true "删除目录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/catalog/del [post]
func CatalogDel(c *gin.Context) {

	var param catalogDelParam
	c.BindJSON(&param)

	//记录操作日志
	var (
		operation              OperateLogServer
		operateLogCreateParams = OperateLogParams{Ctx: c, Model: "catalogs", ParentUuid: param.Uuid, Pattern: PATTERN_CATALOGS, Mold: MOLD_CATALOGS_CREATE}
		operateLogDelParams    = OperateLogParams{Ctx: c, Model: "catalogs", ParentUuid: param.Uuid, Pattern: PATTERN_CATALOGS, Mold: MOLD_CATALOGS_DELETE}
	)
	operation.Start(operateLogCreateParams)
	defer operation.Finish(operateLogDelParams)

	catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)

	var catalog *editor.Catalogs
	catalogsCollection.FindOneAndUpdate(
		c,
		bson.M{
			"uuid": param.Uuid,
		},
		bson.M{
			"$set": bson.M{
				"has_del":     true,
				"has_changed": true,
				"update_time": time.Now().Unix(),
			},
		},
	).Decode(&catalog)
	//更新父目录信息
	updateParentCatalogInfo(c, []string{
		catalog.Parent_uuid,
	})

	servers.ReportFormat(c, true, "删除成功", gin.H{})
}

type catalogShowParam struct {
	Uuid   string `json:"uuid"`
	IsShow bool   `json:"is_show"`
}

// @Tags EditorCatalogAPI(目录接口)
// @Summary 目录显示或隐藏操作
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.catalogShowParam true "是否显示参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/catalog/show [post]
func CatalogShow(c *gin.Context) {

	var param catalogShowParam
	c.BindJSON(&param)
	catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)

	var catalog *editor.Catalogs
	catalogsCollection.FindOneAndUpdate(
		c,
		bson.M{
			"uuid": param.Uuid,
		},
		bson.M{
			"$set": bson.M{
				"is_show":     param.IsShow,
				"has_changed": true,
			},
		},
	).Decode(&catalog)
	//更新父目录信息
	updateParentCatalogInfo(c, []string{
		catalog.Parent_uuid,
	})

	servers.ReportFormat(c, true, "操作成功", gin.H{})
}

type catalogRenameParam struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

// @Tags EditorCatalogAPI(目录接口)
// @Summary 目录重命名
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.catalogRenameParam true "重命名"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/catalog/rename [post]
func CatalogRename(c *gin.Context) {

	var param catalogRenameParam
	c.BindJSON(&param)

	//记录操作日志
	var (
		operation              OperateLogServer
		operateLogCreateParams = OperateLogParams{Ctx: c, Model: "catalogs", ParentUuid: param.Uuid, Pattern: PATTERN_CATALOGS, Mold: MOLD_CATALOGS_CREATE}
		operateLogDelParams    = OperateLogParams{Ctx: c, Model: "catalogs", ParentUuid: param.Uuid, Pattern: PATTERN_CATALOGS, Mold: MOLD_CATALOGS_RENAME}
	)
	operation.Start(operateLogCreateParams)
	defer operation.Finish(operateLogDelParams)

	catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)

	catalogsCollection.UpdateOne(
		c,
		bson.M{
			"uuid": param.Uuid,
		},
		bson.M{
			"$set": bson.M{
				"name": param.Name,
			},
		},
	)

	servers.ReportFormat(c, true, "操作成功", gin.H{})
}

type catalogMoveParam struct {
	Uuid     string `json:"uuid"`
	FromUuid string `json:"from_uuid"`
	ToUuid   string `json:"to_uuid"`
}

func checkCatalogMove(c *gin.Context, param catalogMoveParam) bool {
	if param.ToUuid == param.Uuid {
		return true
	}

	result := findAllCatalogSonUuid(c, param.ToUuid)

	for _, value := range result {
		if value == param.Uuid {
			return true
		}
	}

	return false
}

// 查询所有的子集
func findAllCatalogSonUuid(c *gin.Context, to_uuid string) (uuids []string) {
	var catalogs []editor.Catalogs

	catCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)
	cusor, _ := catCollection.Find(
		c,
		bson.M{
			"uuid": to_uuid,
		},
	)
	cusor.All(c, &catalogs)
	cusor.Close(c)

	if catalogs != nil {
		for _, cat := range catalogs {
			uuids = append(uuids, cat.Parent_uuid)
			data := findAllCatalogSonUuid(c, cat.Parent_uuid)
			for _, item := range data {
				uuids = append(uuids, item)
			}
		}
	}
	return uuids
}

// @Tags EditorCatalogAPI(目录接口)
// @Summary 移动目录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.catalogMoveParam true "移动目录：uuid当前目录，fromUuid是从哪个目录移出，toUuid是移到哪个目录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/catalog/move [post]
func CatalogMove(c *gin.Context) {
	var param catalogMoveParam
	c.BindJSON(&param)

	//检测移动目录不能是自己的目录及子目录
	if check_bool := checkCatalogMove(c, param); check_bool {
		servers.ReportFormat(c, false, "操作失败，移动目录不能是自己及子目录！", gin.H{})
		return
	}

	//记录操作日志
	var (
		operation              OperateLogServer
		operateLogCreateParams = OperateLogParams{Ctx: c, Model: "catalogs", ParentUuid: param.Uuid, Pattern: PATTERN_CATALOGS, Mold: MOLD_CATALOGS_CREATE}
		operateLogDelParams    = OperateLogParams{Ctx: c, Model: "catalogs", ParentUuid: param.Uuid, Pattern: PATTERN_CATALOGS, Mold: MOLD_CATALOGS_MOVE}
	)
	operation.Start(operateLogCreateParams)
	defer operation.Finish(operateLogDelParams)

	catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)
	catalogsCollection.UpdateOne(
		c,
		bson.M{
			"uuid": param.Uuid,
		},
		bson.M{
			"$set": bson.M{
				"parent_uuid": param.ToUuid,
				"has_changed": true,
				"update_time": time.Now().Unix(),
			},
		},
	)
	//更新父目录和课程内容版本信息
	updateParentCatalogInfo(c, []string{param.FromUuid})
	updateParentCatalogInfo(c, []string{param.ToUuid})
	servers.ReportFormat(c, true, "移动成功", gin.H{})
}

type CatalogCopyParam struct {
	Uuids    []string `json:"uuids"`
	ToUuid   string   `json:"to_uuid"`
	SameLang bool     `json:"same_lang"`
}

// @Tags EditorCatalogAPI(目录接口)
// @Summary 复制目录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.CatalogCopyParam true "复制目录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/catalog/copy [post]
func CatalogCopy(c *gin.Context) {
	var param *CatalogCopyParam
	c.BindJSON(&param)
	copyCatalogInfo(param, true) //更新父目录信息
	updateParentCatalogInfo(c, []string{param.ToUuid})
	servers.ReportFormat(c, true, "复制成功", gin.H{})
}

//复制目录
type catalogInfo struct {
	Parent_uuid   string            `bson:"parent_uuid" json:"parent_uuid"`
	Has_changed   bool              `bson:"has_changed" json:"has_changed"`
	Update_time   int64             `bson:"update_time" json:"update_time"`
	Tags          []string          `bson:"tags" json:"tags"`
	List_order    int               `bson:"list_order" json:"list_order"`
	Is_show       bool              `bson:"is_show" json:"is_show"`
	Title         map[string]string `bson:"title" json:"title"`
	Name          string            `bson:"name" json:"name"`
	Has_del       bool              `bson:"has_del" json:"has_del"`
	Flag          []string          `bson:"flag" json:"flag"`
	Desc          map[string]string `bson:"desc" json:"desc"`
	Uuid          string            `bson:"uuid" json:"uuid"`
	Type          string            `bson:"type" json:"type"`
	Cover         []string          `bson:"cover" json:"cover"`
	Content_model string            `bson:"content_model" json:"content_model"`
	AttrTag       string            `bson:"attr_tag" json:"attr_tag"`
}

func copyCatalogInfo(
	param *CatalogCopyParam,
	isCatalog bool,
) {
	qmlog.QMLog.Info("开始复制目录")

	var wg sync.WaitGroup

	ctx := context.Background()
	catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)

	var catalogs []*catalogInfo
	cusor, _ := catalogsCollection.Find(ctx, bson.M{
		"uuid": bson.M{"$in": param.Uuids},
	})
	defer cusor.Close(ctx)

	cusor.All(ctx, &catalogs)
	catalogsInfo := make([]interface{}, 0)

	for _, c := range catalogs {
		oldUuid := c.Uuid
		newUuid := uuid.NewV4().String()
		if isCatalog {
			c.Name = "复制-" + c.Name
		}
		c.Uuid = newUuid
		c.Parent_uuid = param.ToUuid
		c.List_order = getCatalogsNextListOrder(param.ToUuid)
		c.Has_changed = true
		c.Update_time = time.Now().Unix()
		catalogsInfo = append(catalogsInfo, c)

		//如果当前目录是目录类型，则继续往下查找
		if c.Type == "catalog" { //目录类型
			go func(oldUuid, newUuid string, sameLang bool) {

				wg.Add(1)
				defer wg.Done()

				copyChildCatalog(
					oldUuid,
					newUuid,
					sameLang,
				)
			}(oldUuid, newUuid, param.SameLang)

		} else if c.Type == "content" { //文件内容类型,则拷贝文件内容
			var contents []bson.M
			contentCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(c.Content_model)
			contentCusor, _ := contentCollection.Find(
				ctx,
				bson.M{
					"parent_uuid": oldUuid,
				},
			)

			contentCusor.All(ctx, &contents)
			contentCusor.Close(ctx)

			if contents != nil {
				childContents := make([]interface{}, len(contents))

				for kk, cc := range contents {
					delete(cc, "_id")
					delete(cc, "code")
					cc["parent_uuid"] = newUuid
					cc["uuid"] = uuid.NewV4().String()
					cc["has_changed"] = true
					cc["list_order"] = getContentNextListOrder(newUuid, c.Content_model)

					if !param.SameLang {
						if c.Content_model == "content_model_pro_sound" {
							cc["sentence_temp"] = cc["sentence"]
							cc["options_temp"] = cc["options"]
							cc["sentence"] = ""
							cc["options"] = []string{}
							cc["sound"] = ""
						}
					}
					childContents[kk] = cc
				}
				contentCollection.InsertMany(ctx, childContents)
			}
		}
	}
	catalogsCollection.InsertMany(ctx, catalogsInfo)

	wg.Wait()

	qmlog.QMLog.Info("结束复制目录")

}

func getContentNextListOrder(parent_uuid, content_model string) int {
	operations := []bson.D{
		bson.D{{"$match", bson.M{"parent_uuid": parent_uuid}}},
		bson.D{{"$group", bson.D{
			{"_id", "$parent_uuid"}, {"max", bson.M{"$max": "$list_order"}}},
		}},
	}

	type MaxNumber struct {
		Max int64  `bson:"max" json:"max"`
		Id  string `bson:"_id" json:"id"`
	}

	var result []MaxNumber

	cursor, _ := mgdb.Aggregate(
		mgdb.EnvEditor,
		EDITOR_DB,
		content_model,
		operations,
	)
	var ctx context.Context
	defer cursor.Close(ctx)
	cursor.All(ctx, &result)

	step := 10
	if result != nil {
		max := result[0].Max
		step = step + int(max)
	}
	return step
}

// 获取下一个 ListOrder 的值
func getCatalogsNextListOrder(parent_uuid string) int {
	operations := []bson.D{
		bson.D{{"$match", bson.M{"parent_uuid": parent_uuid}}},
		bson.D{{"$group", bson.D{
			{"_id", "$parent_uuid"}, {"max", bson.M{"$max": "$list_order"}}},
		}},
	}

	type MaxNumber struct {
		Max int64  `bson:"max" json:"max"`
		Id  string `bson:"_id" json:"id"`
	}

	var result []MaxNumber

	cursor, _ := mgdb.Aggregate(
		mgdb.EnvEditor,
		EDITOR_DB,
		tblCatalogs,
		operations,
	)
	var ctx context.Context
	defer cursor.Close(ctx)
	cursor.All(ctx, &result)

	step := 10
	if result != nil {
		max := result[0].Max
		step = step + int(max)
	}
	return step
}

//复制子目录
func copyChildCatalog(
	fromParentUuid,
	toParentUuid string,
	sameLang bool,
) {
	ctx := context.Background()
	catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)
	var catalogs []*catalogInfo
	cusor, e := catalogsCollection.Find(ctx, bson.M{
		"parent_uuid": fromParentUuid,
	})

	if e != nil {
		return
	}
	cusor.All(ctx, &catalogs)

	cusor.Close(ctx)
	catalogsInfo := make([]interface{}, 0)

	for _, c := range catalogs {
		oldUuid := c.Uuid
		newUuid := uuid.NewV4().String()

		c.Uuid = newUuid
		c.Parent_uuid = toParentUuid
		c.Has_changed = true
		c.Update_time = time.Now().Unix()
		catalogsInfo = append(catalogsInfo, c)

		//如果当前目录是目录类型，则继续往下查找
		if c.Type == "catalog" { //目录类型
			copyChildCatalog(
				oldUuid,
				newUuid,
				sameLang,
			)
		} else if c.Type == "content" { //文件内容类型,则拷贝文件内容
			var contents []bson.M
			contentCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(c.Content_model)
			contentCusor, _ := contentCollection.Find(
				ctx,
				bson.M{
					"parent_uuid": oldUuid,
				},
			)

			contentCusor.All(ctx, &contents)
			contentCusor.Close(ctx)

			if contents != nil {
				childContents := make([]interface{}, len(contents))

				for kk, cc := range contents {
					delete(cc, "_id")
					delete(cc, "code")
					cc["parent_uuid"] = newUuid
					cc["uuid"] = uuid.NewV4().String()
					cc["has_changed"] = true

					if !sameLang {
						if c.Content_model == "content_model_pro_sound" {
							cc["sentence_temp"] = cc["sentence"]
							cc["options_temp"] = cc["options"]
							cc["sentence"] = ""
							cc["options"] = []string{}
							cc["sound"] = ""
						}
					}
					childContents[kk] = cc
				}
				contentCollection.InsertMany(ctx, childContents)
			}
		}
	}

	catalogsCollection.InsertMany(ctx, catalogsInfo)
	// return catalogsInfo
}

//更新上层的目录相关信息
type parentUuid struct {
	ParentUuid string `bson:"parent_uuid" json:"parent_uuid"`
}

func updateParentCatalogInfo(c *gin.Context, parentUuids []string) {
	var uuids []*parentUuid
	ctx := c
	catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)
	catalogsCollection.UpdateMany(
		ctx,
		bson.M{
			"uuid": bson.M{"$in": parentUuids},
		},
		bson.M{
			"$set": bson.M{
				"has_changed": true,
				"update_time": time.Now().Unix(),
			},
		},
	)
	cusor, e := catalogsCollection.Find(
		ctx,
		bson.M{
			"uuid": bson.M{"$in": parentUuids},
		},
		options.Find().SetProjection(bson.M{
			"_id":         0,
			"parent_uuid": 1,
		}),
	)

	if e != nil {
		checkErr(c, e)
		return
	}
	cusor.All(ctx, &uuids)
	cusor.Close(context.Background())
	//更新目录信息
	if uuids != nil {
		tmpUuids := []string{}
		for _, uid := range uuids {
			tmpUuids = append(tmpUuids, uid.ParentUuid)
		}
		updateParentCatalogInfo(ctx, tmpUuids)
	} else {
		//更新内容版本信息
		contentCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentInfo)
		contentCollection.UpdateMany(
			ctx,
			bson.M{
				"uuid": bson.M{"$in": parentUuids},
			},
			bson.M{
				"$set": bson.M{
					"has_changed": true,
					"update_time": time.Now().Unix(),
				},
			},
		)
	}
}
