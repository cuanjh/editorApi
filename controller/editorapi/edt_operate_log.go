package editorapi

import (
	"editorApi/controller/servers"
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"editorApi/middleware"
	"editorApi/tools/helpers"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"sync"
	"time"
)

var tbOperateLog = "operate_log"
var finalParentUuid string
var finalCatalogs []editor.Catalogs
var lock sync.Mutex

type OperateLogServer struct {
}

type OperateLogParams struct {
	Ctx        *gin.Context
	Model      string
	ParentUuid string
	Pattern    int
	Mold       int
}

/**
开始:数据保存前执行
*/
func (s *OperateLogServer) Start(params OperateLogParams) {
	switch params.Pattern {
	case PATTERN_CONTENT:
		// 判断是否写入过日志
		operateLog := s.checkOperateLog(params.ParentUuid, params.Pattern)
		if !helpers.Empty(operateLog) {
			return
		}
		s.addOperateLog(params)
	case PATTERN_CATALOGS:
		operateLog := s.checkOperateLog(params.ParentUuid, params.Pattern)
		if !helpers.Empty(operateLog) {
			return
		}
		s.addOperateLog(params)
	}
}

/**
结束:数据保存完后执行
*/
func (s OperateLogServer) Finish(params OperateLogParams) {
	switch params.Pattern {
	case PATTERN_CONTENT:
		s.addOperateLog(params)
	// 异步操作
	/**
	initNats.NatsConn.Publish("OperateLog", &OperateLogParams{
		Ctx:        params.Ctx,
		Model:      params.Model,
		ParentUuid: params.ParentUuid,
		Pattern:    params.Pattern,
		Mold:       params.Mold,
	})
	**/
	case PATTERN_CATALOGS:
		s.addOperateLog(params)
	}
}

/**
异步处理
*/
func HanderOperateLog(params *OperateLogParams) {
	switch params.Pattern {
	case PATTERN_CONTENT:
		var operateLogServer OperateLogServer
		operateLogServer.addOperateLog(OperateLogParams{
			Ctx:        params.Ctx,
			Model:      params.Model,
			ParentUuid: params.ParentUuid,
			Pattern:    params.Pattern,
			Mold:       params.Mold,
		})
	}
}

func (s OperateLogServer) checkOperateLog(parent_uuid string, pattern int) (operateLog *editor.OperateLog) {
	filter := bson.M{
		"parent_uuid": parent_uuid,
		"pattern":     pattern,
	}
	mgdb.FindOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tbOperateLog,
		filter,
		nil,
		&operateLog,
	)
	return
}

//添加日志信息
func (s OperateLogServer) addOperateLog(params OperateLogParams) (insertResult *mongo.InsertOneResult, err error) {
	lock.Lock()
	token := params.Ctx.Request.Header.Get("x-token")
	claims, _ := s.getUserInfo(token)

	finalCatalogs = []editor.Catalogs{}
	menusData := s.getMenus(params.ParentUuid)

	var contentsString string
	switch params.Pattern {
	case PATTERN_CONTENT:
		collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(params.Model)
		cusor, _ := collection.Find(
			params.Ctx,
			bson.M{
				"parent_uuid": params.ParentUuid,
			},
			options.Find().SetSort(bson.M{
				"list_order": 1,
			}).SetProjection(bson.M{
				"_id": 0,
			}),
		)
		defer cusor.Close(params.Ctx)

		var contents []bson.M
		cusor.All(params.Ctx, &contents)
		result, _ := json.Marshal(contents)

		contentsString = string(result)
	case PATTERN_CATALOGS:
		var catalog editor.Catalogs
		mgdb.FindOne(
			mgdb.EnvEditor,
			mgdb.DbEditor,
			tblCatalogs,
			bson.M{
				"uuid": params.ParentUuid,
			},
			nil,
			&catalog,
		)

		result, _ := json.Marshal(catalog)
		contentsString = string(result)
	}

	if !helpers.Empty(contentsString) {
		operateLogCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tbOperateLog)
		// 添加时，需要保证当前只有一个显示

		var param editor.OperateLog
		param.CreatedTime = time.Now().UnixNano()
		param.Content = contentsString
		param.Uuid = uuid.NewV4().String()
		param.ParentUuid = params.ParentUuid
		param.OperateDate = time.Now().Format("2006-01-02 15:04:05")
		param.Language = menusData.Lang
		param.LangCode = menusData.LangCode
		param.Version = menusData.Version
		param.VersionId = menusData.VersionId
		param.Course = menusData.Course
		param.CourseID = menusData.CourseID
		param.Catalogs = menusData.Menus
		param.Pattern = params.Pattern
		param.Mold = params.Mold
		param.UserId = claims.UUID.String()
		param.UserName = claims.NickName
		param.Ip = params.Ctx.ClientIP()

		insertResult, err = operateLogCollection.InsertOne(params.Ctx, param)
	}

	lock.Unlock()
	return
}

func (s *OperateLogServer) getUserInfo(token string) (claims *middleware.CustomClaims, err error) {
	j := middleware.NewJWT()
	claims, err = j.ParseToken(token)
	return
}

type MenusData struct {
	Menus     string
	Version   string
	VersionId string
	Course    string
	CourseID  string
	Lang      string
	LangCode  string
}

func (s *OperateLogServer) getMenus(parent_uuid string) (menusData MenusData) {
	s.getCatalogs(parent_uuid)

	for _, item := range finalCatalogs {
		finalParentUuid = item.Parent_uuid
	}
	var menus []string
	for j := len(finalCatalogs) - 1; j >= 0; j-- {
		menus = append(menus, finalCatalogs[j].Name)
	}
	menusData.Menus = strings.Join(menus, ">")

	// 版本
	var courseContentInfos editor.Course_content_infos
	if !helpers.Empty(finalParentUuid) {
		mgdb.FindOne(
			mgdb.EnvEditor,
			mgdb.DbEditor,
			tblContentInfo,
			bson.M{
				"uuid":    finalParentUuid,
				"has_del": false,
			},
			nil,
			&courseContentInfos,
		)
		menusData.Version = courseContentInfos.Name
		menusData.VersionId = courseContentInfos.Uuid
	}

	//course 课程
	var courseInfos editor.Course_infos
	if !helpers.Empty(courseContentInfos) && !helpers.Empty(courseContentInfos.Parent_uuid) {
		mgdb.FindOne(
			mgdb.EnvEditor,
			mgdb.DbEditor,
			tblCourseInfos,
			bson.M{
				"uuid":    courseContentInfos.Parent_uuid,
				"has_del": false,
			},
			nil,
			&courseInfos,
		)
		menusData.Course = courseInfos.Title["zh-CN"]
		menusData.CourseID = courseInfos.Uuid
	}

	//lang 语言
	var courseLangs editor.Course_langs
	if !helpers.Empty(courseInfos) && !helpers.Empty(courseInfos.Lan_code) {
		mgdb.FindOne(
			mgdb.EnvEditor,
			mgdb.DbEditor,
			tblCourseLangs,
			bson.M{
				"lan_code": courseInfos.Lan_code,
				"has_del":  false,
			},
			nil,
			&courseLangs,
		)
		menusData.Lang = courseLangs.Title["zh-CN"]
		menusData.LangCode = courseLangs.Lan_code
	}

	return
}

func (s *OperateLogServer) getCatalogs(parent_uuid string) {
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
		finalCatalogs = append(finalCatalogs, catalog)
		s.getCatalogs(catalog.Parent_uuid)
	}
	return
}

type OperateLogListParams struct {
	PageNo      int64       `json:"page_no"`
	PageSize    int64       `json:"page_size"`
	OperateDate OperateDate `json:"operate_date"`
	UserId      string      `json:"user_id"`
}

type OperateDate struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// @Tags EditorOnlineAPI（日志操作接口）
// @Summary 日志操作列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.OperateLogListParams true "数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/operate_log/list [post]
func OperateLogList(ctx *gin.Context) {
	parameter := OperateLogListParams{}
	ctx.BindJSON(&parameter)
	var limit int64
	var skip int64
	if helpers.Empty(parameter.PageSize) {
		limit = 40
	} else {
		limit = parameter.PageSize
	}
	if parameter.PageNo <= 0 || helpers.Empty(parameter.PageNo) {
		parameter.PageNo = 1
	}
	skip = (parameter.PageNo - 1) * limit
	var (
		err   error
		cusor *mongo.Cursor
	)

	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tbOperateLog)

	// 过滤条件
	var filter = bson.D{}
	if !helpers.Empty(parameter.UserId) {
		filter = append(filter, bson.E{"user_id", parameter.UserId})
	}

	if !helpers.Empty(parameter.OperateDate) {
		filter = append(filter, bson.E{"operate_date", bson.M{"$gte": parameter.OperateDate.Start, "$lt": parameter.OperateDate.End}})
	}

	cusor, err = collection.Find(
		ctx,
		filter,
		options.Find().SetProjection(map[string]int{
			"_id": 0,
		}).SetSort(bson.M{"created_time": -1}).SetLimit(limit).SetSkip(skip),
	)
	defer cusor.Close(ctx)
	if err != nil {
		checkErr(ctx, err)
		return
	}
	var operate_log []bson.M
	cusor.All(ctx, &operate_log)

	// 查询条数
	total, _ := collection.CountDocuments(ctx, filter)

	servers.ReportFormat(ctx, true, "成功", gin.H{
		"operate_log": operate_log,
		"total":       total,
	})
}

type OperateLogDetailsParams struct {
	Pattern    int    `json:"pattern"`
	Uuid       string `json:"uuid"`
	UserId     string `json:"user_id"`
	ParentUuid string `json:"parent_uuid"`
}

// @Tags EditorOnlineAPI（日志操作接口）
// @Summary 日志操作详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.OperateLogDetailsParams true "数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/operate_log/details [post]
func OperateLogDetails(ctx *gin.Context) {
	parameter := OperateLogDetailsParams{}
	ctx.BindJSON(&parameter)

	var operateLog *editor.OperateLog
	filter := bson.M{
		"pattern":     parameter.Pattern,
		"user_id":     parameter.UserId,
		"parent_uuid": parameter.ParentUuid,
		"uuid":        parameter.Uuid,
	}
	mgdb.FindOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tbOperateLog,
		filter,
		nil,
		&operateLog,
	)

	var result []*editor.OperateLog
	result = append(result, operateLog)

	var oldOperateLog []*editor.OperateLog
	oldFilter := bson.M{
		"pattern":     parameter.Pattern,
		"user_id":     parameter.UserId,
		"parent_uuid": parameter.ParentUuid,
		"created_time": bson.M{
			"$lt": operateLog.CreatedTime,
		},
	}
	mgdb.Find(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tbOperateLog,
		oldFilter,
		bson.M{"created_time": -1},
		nil,
		0,
		1,
		&oldOperateLog,
	)
	if !helpers.Empty(oldOperateLog) {
		result = append(result, oldOperateLog[0])
	}

	servers.ReportFormat(ctx, true, "成功", gin.H{
		"operate_log": result,
	})
}

type RollbackParams struct {
	Uuid string `json:"uuid"`
}

// @Tags EditorOnlineAPI（日志操作接口）
// @Summary 日志操作回滚
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.RollbackParams true "数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/operate_log/rollback [post]
func Rollback(ctx *gin.Context) {
	parameter := RollbackParams{}
	ctx.BindJSON(&parameter)

	var operation OperateLogServer
	var operateLog editor.OperateLog
	filter := bson.M{
		"uuid": parameter.Uuid,
	}
	mgdb.FindOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tbOperateLog,
		filter,
		nil,
		&operateLog,
	)

	if !helpers.Empty(operateLog) && !helpers.Empty(operateLog.Pattern) {
		switch operateLog.Pattern {
		case PATTERN_CONTENT:
			if !helpers.Empty(operateLog.Content) {
				// 查询内容的模型
				var catalog editor.Catalogs
				mgdb.FindOne(
					mgdb.EnvEditor,
					mgdb.DbEditor,
					tblCatalogs,
					bson.M{
						"uuid": operateLog.ParentUuid,
					},
					nil,
					&catalog,
				)

				collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(catalog.Content_model)
				var contents []bson.M
				json.Unmarshal([]byte(operateLog.Content), &contents)
				for _, c := range contents {

					if uid, ok := c["uuid"]; !ok || uid == "" {
						c["uuid"] = uuid.NewV4().String()
					}
					c["parent_uuid"] = operateLog.ParentUuid
					c["has_changed"] = true
					c["update_time"] = time.Now().Unix()
					if _, ok := c["has_del"]; !ok {
						c["has_del"] = false
					}

					collection.UpdateOne(
						ctx,
						bson.M{
							"uuid": c["uuid"],
						},
						bson.M{
							"$set": c,
						},
						options.Update().SetUpsert(true),
					)
				}
				updateParentCatalogInfo(ctx, []string{catalog.Parent_uuid})
				//记录操作日志
				var operateLogRollbackParams = OperateLogParams{Ctx: ctx, Model: catalog.Content_model, ParentUuid: operateLog.ParentUuid, Pattern: PATTERN_CONTENT, Mold: MOLD_CONTENT_ROLLBACK}
				defer operation.Finish(operateLogRollbackParams)
			}
		case PATTERN_CATALOGS:
			if !helpers.Empty(operateLog.Content) {
				catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)

				var catalog editor.Catalogs
				json.Unmarshal([]byte(operateLog.Content), &catalog)
				catalog.Has_changed = true
				catalog.Update_time = time.Now().Unix()
				if operateLog.Mold == MOLD_CATALOGS_MOVE {
					tmp_catalog := editor.Catalogs{}
					mgdb.FindOne(
						mgdb.EnvEditor,
						mgdb.DbEditor,
						tblCatalogs,
						bson.M{
							"uuid": catalog.Uuid,
						},
						nil,
						&tmp_catalog,
					)
					if catalog.Parent_uuid != tmp_catalog.Parent_uuid {
						updateParentCatalogInfo(ctx, []string{tmp_catalog.Parent_uuid})
					}
				}

				catalogsCollection.FindOneAndUpdate(ctx,
					bson.M{
						"uuid": catalog.Uuid,
					},
					bson.M{
						"$set": catalog,
					},
				).Decode(&catalog)

				updateParentCatalogInfo(ctx, []string{catalog.Parent_uuid})
				//记录操作日志
				var operateLogRollbackParams = OperateLogParams{Ctx: ctx, Model: "catalogs", ParentUuid: catalog.Uuid, Pattern: PATTERN_CATALOGS, Mold: MOLD_CATALOGS_ROLLBACK}
				defer operation.Finish(operateLogRollbackParams)
			}
		}
	}

	servers.ReportFormat(ctx, true, "回滚成功", nil)
}

// 定义常量
const (
	MOLD_CONTENT_CREATE   int = 11 //内容新建
	MOLD_CONTENT_EDIT     int = 12 //内容编辑
	MOLD_CONTENT_DELETE   int = 13 //内容删除
	MOLD_CONTENT_ROLLBACK int = 14 //内容回滚

	MOLD_CATALOGS_CREATE   int = 20 //目录新建
	MOLD_CATALOGS_EDIT     int = 21 //目录编辑
	MOLD_CATALOGS_DELETE   int = 22 //目录删除
	MOLD_CATALOGS_MOVE     int = 23 //目录移动
	MOLD_CATALOGS_RENAME   int = 24 //目录改名称
	MOLD_CATALOGS_ROLLBACK int = 25 //目录回滚
)

const (
	PATTERN_CONTENT  int = 1 //内容
	PATTERN_CATALOGS int = 2 //目录
)
