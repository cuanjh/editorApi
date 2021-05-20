package editorapi

import (
	"context"
	"editorApi/controller/servers"
	"editorApi/init/initNats"
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"editorApi/tools/helpers"
	"errors"
	"strconv"
	"time"
	"tkCommon/cmfunc"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tblCourseInfos string = "course_infos"

type courseListParams struct {
	PageNo   int64  `json:"pageNo"`
	PageSize int64  `json:"pageSize"`
	LanCode  string `json:"lan_code"`
}

// @Tags EditorCourseAPI（课程接口）
// @Summary 课程列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.courseListParams true "课程列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/course/list [post]
func CourseList(c *gin.Context) {

	var params courseListParams
	c.BindJSON(&params)
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCourseInfos)
	var courseLists []*editor.Course_infos
	filter := bson.M{
		"has_del": false,
	}

	if params.LanCode != "" {
		filter["lan_code"] = params.LanCode
	}
	var limit int64 = 40
	if params.PageSize != 0 {
		limit = params.PageSize
	}
	offset := params.PageNo * limit

	cusor, err := collection.Find(
		c,
		filter,
		options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{
			"list_order": 1,
		}),
	)
	defer cusor.Close(c)

	if err != nil {
		checkErr(c, err)
		return
	}

	if err := cusor.Err(); err != nil {
		checkErr(c, err)
		return
	}
	cusor.All(c, &courseLists)

	servers.ReportFormat(c, true, "课程列表", gin.H{
		"courses": courseLists,
	})
}

type courseAddParam struct {
	Is_show      bool              `json:"is_show"`
	Uuid         string            `json:"uuid"`
	Name         string            `json:"name"`
	Code         string            `json:"code"`
	Flag         []string          `json:"flag"`
	Course_type  int64             `json:"course_type"`
	Cover        []string          `json:"cover"`
	Lan_code     string            `json:"lan_code"`
	Desc         map[string]string `json:"desc"`
	Title        map[string]string `json:"title"`
	Tags         []string          `json:"tags"`
	Has_del      bool              `json:"has_del"`
	HasDict      bool              `bson:"has_dict" json:"has_dict"`
	SoundActors  []soundActor      `bson:"sound_actors" json:"sound_actors"`
	DefaultActor string            `bson:"default_actor" json:"default_actor"`
}

type soundActor struct {
	Role   string `bson:"role" json:"role"`
	Name   string `bson:"name" json:"name"`
	Photo  string `bson:"photo" json:"photo"`
	Gender int    `bson:"gender" json:"gender"`
	Sound  string `bson:"sound" json:"sound"`
	Desc   string `json:"desc" bson:"desc"` // 描述
	City   string `json:"city" bson:"city"` // 城市
}

// @Tags EditorCourseAPI（课程接口）
// @Summary 添加课程
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.courseAddParam true "添加课程"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/course/add [post]
func CourseAdd(c *gin.Context) {

	var param courseAddParam
	c.BindJSON(&param)
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCourseInfos)
	if courseNum, err := collection.CountDocuments(c, bson.M{
		"lan_code": param.Lan_code,
		"code":     param.Code,
		"has_del":  false,
	}); courseNum > 0 || err != nil {
		err = errors.New("编码重复")
		checkErr(c, err)
		return
	}
	param.Has_del = false
	param.Is_show = true
	param.Uuid = uuid.NewV4().String()
	_, e := collection.InsertOne(c, param)
	if e != nil {
		checkErr(c, e)
		return
	}
	servers.ReportFormat(c, true, "成功", gin.H{
		"uuid": param.Uuid,
	})

}

type courseEditInfo struct {
	Is_show      bool              `json:"is_show"`
	Flag         []string          `json:"flag"`
	Name         string            `json:"name"`
	Course_type  int64             `json:"course_type"`
	Cover        []string          `json:"cover"`
	Desc         map[string]string `json:"desc"`
	Title        map[string]string `json:"title"`
	Tags         []string          `json:"tags"`
	HasDict      bool              `bson:"has_dict" json:"has_dict"`
	SoundActors  []soundActor      `bson:"sound_actors" json:"sound_actors"`
	DefaultActor string            `bson:"default_actor" json:"default_actor"`
}

type courseEditParam struct {
	Uuid     string         `json:"uuid"`
	EditInfo courseEditInfo `json:"editInfo"`
}

// @Tags EditorCourseAPI（课程接口）
// @Summary 修改课程
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.courseEditParam true "修改课程"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/course/edit [post]
func CourseEdit(c *gin.Context) {

	var param courseEditParam
	c.BindJSON(&param)
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCourseInfos)
	_, err := collection.UpdateOne(c, bson.M{
		"uuid": param.Uuid,
	}, bson.M{
		"$set": param.EditInfo,
	})
	if err != nil {
		checkErr(c, err)
		return
	}

	servers.ReportFormat(c, true, "更新成功", gin.H{})

}

type coruseDelParam struct {
	Uuid string `json:"uuid"`
}

// @Tags EditorCourseAPI（课程接口）
// @Summary 删除课程
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.coruseDelParam true "删除课程"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/course/del [post]
func CourseDel(c *gin.Context) {

	var param coruseDelParam
	c.BindJSON(&param)

	dc, _ := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentInfo).CountDocuments(c, bson.M{
		"parent_uuid": param.Uuid,
		"has_del":     false,
	})
	if dc > 0 {
		servers.ReportFormat(c, false, "此课程存在内容版本", gin.H{})
	} else {
		collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCourseInfos)
		_, err := collection.UpdateOne(
			c,
			bson.M{
				"uuid": param.Uuid,
			},
			bson.M{
				"$set": bson.M{"has_del": true},
			},
		)
		if err != nil {
			checkErr(c, err)
			return
		}

		servers.ReportFormat(c, true, "删除成功", gin.H{
			"err": err,
		})
	}
}

// @Tags EditorCourseAPI（课程接口）
// @Summary 课程类型
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/course/types [get]
func CourseTypes(c *gin.Context) {

	servers.ReportFormat(c, true, "获取成功", gin.H{
		"types": []bson.M{
			bson.M{
				"name": "Pro官方课程",
				"type": 0,
			},
			bson.M{
				"name": "Kid官方课程",
				"type": 3,
			},
			bson.M{
				"name": "小学官方课程",
				"type": 5,
			},
			bson.M{
				"name": "小学人教课程",
				"type": 7,
			},
		},
	})
}

type CourseGetDetailParam struct {
	CourseCode string `json:"course_code"`
}

type CoursePriDetail struct {
	Title      string `json:"title"`
	Flag       string `json:"flag" bson:"flag"`
	CourseCode string `json:"course_code" bson:"course_code"`
	Desc       string `json:"desc" bson:"desc"`
	Content    struct {
		Lesson   string `json:"lesson" bson:"lesson"`
		Word     string `json:"word" bson:"word"`
		Sentence string `json:"sentence" bson:"sentence"`
	} `json:"content" bson:"content"`
	Features []struct {
		ListOrder int    `json:"list_order" bson:"list_order"`
		Title     string `json:"title" bson:"title"`
		Desc      string `json:"desc" bson:"desc"`
		Img       string `json:"img" bson:"img"`
	} `json:"features" bson:"features"`
	HasDict     bool         `bson:"has_dict" json:"has_dict"`
	SoundActors []soundActor `bson:"sound_actors" json:"sound_actors"`
}

// @Tags EditorCourseAPI（课程接口）
// @Summary 获取小学英语课程详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.CourseGetDetailParam true "课程编码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/course/detail [post]
func CourseDetail(c *gin.Context) {
	var pm CourseGetDetailParam
	c.BindJSON(&pm)
	detail := &CoursePriDetail{}
	mgdb.FindOne(
		mgdb.EnvOnline,
		mgdb.DbContent,
		"course_pri_detail",
		bson.M{
			"course_code": pm.CourseCode,
		},
		nil,
		&detail,
	)

	detail.Flag = cmfunc.GetCourseAssets(detail.Flag)
	for k, p := range detail.Features {
		p.Img = cmfunc.GetCourseAssets(p.Img)
		detail.Features[k] = p
	}
	servers.ReportFormat(c, true, "成功", gin.H{
		"detail": detail,
	})
}

// @Tags EditorCourseAPI（课程接口）
// @Summary 编辑小学英语课程详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.CoursePriDetail true "课程详情"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/course/detail/edit [post]
func CourseDetailEdit(c *gin.Context) {
	var detail CoursePriDetail
	c.BindJSON(&detail)
	mgdb.UpdateOne(
		mgdb.EnvOnline,
		mgdb.DbContent,
		"course_pri_detail",
		bson.M{
			"course_code": detail.CourseCode,
		},
		bson.M{
			"$set": detail,
		},
		true,
	)

	servers.ReportFormat(c, true, "成功", gin.H{
		"detail": detail,
	})
}

type contentUnlockParam struct {
	CourseCode string `bson:"course_code" json:"course_code" validate:"required"`
	TalkmateId string `bson:"talkmate_id" json:"talkmate_id"`
}

// @Tags EditorCourseAPI（课程接口）
// @Summary 课程解锁
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.contentUnlockParam true "参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/course/unlock [post]
func ContentUnlock(ctx *gin.Context) {
	var param contentUnlockParam
	ctx.BindJSON(&param)

	if helpers.Empty(param.CourseCode) || helpers.Empty(param.TalkmateId) {
		servers.ReportFormat(ctx, false, "参数不能为空！", gin.H{})
		return
	}

	var user *struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	mgdb.FindOne(
		mgdb.EnvOnline,
		mgdb.DbKuyu,
		tblUsers,
		bson.M{
			"talkmate_id": param.TalkmateId,
		},
		nil,
		&user,
	)

	if helpers.Empty(user) {
		servers.ReportFormat(ctx, false, "用户不存在！", gin.H{})
		return
	}

	userId := user.ID.Hex()

	//判断是否订阅了
	var usersSubscribeCourse UsersSubscribeCourse
	mgdb.FindOne(
		mgdb.EnvOnline,
		mgdb.DbKuyu,
		"users_subscribe_course",
		bson.M{
			"user_id":     userId,
			"course_code": param.CourseCode,
		},
		nil,
		&usersSubscribeCourse,
	)

	if helpers.Empty(usersSubscribeCourse) {
		servers.ReportFormat(ctx, false, "此用户未订阅此课程！", gin.H{})
		return
	}

	//if usersSubscribeCourse.HasUnlockTransfer {
	//	servers.ReportFormat(ctx, false, "此用户订阅的课程已解锁！", gin.H{})
	//	return
	//}

	// 异步操作
	initNats.NatsConn.Publish("ContentUnlock",
		&ContentUnlockMsg{
			CourseCode: param.CourseCode,
			UserId:     userId,
		},
	)

	servers.ReportFormat(ctx, true, "提交成功，解锁需要等待几分钟！", nil)
}

type ContentUnlockMsg struct {
	CourseCode string `bson:"course_code" json:"course_code"`
	UserId     string `bson:"user_id" json:"user_id"`
}

func HanderContentUnlock(param *ContentUnlockMsg) {
	var ctx context.Context
	var courseCode string

	var courseInfos struct {
		Uuid       string `bson:"uuid" json:"uuid"`
		CourseType int64  `json:"course_type"`
	}
	courseCode = param.CourseCode
	mgdb.FindOne(
		mgdb.EnvOnline,
		courseContentDb,
		tblCourseInfos,
		bson.M{
			"code": courseCode,
		},
		nil,
		&courseInfos,
	)

	filter := bson.M{
		"has_del":     false,
		"parent_uuid": courseInfos.Uuid,
	}
	collection := mgdb.OnlineClient.Database(courseContentDb).Collection(tblContentInfo)
	cusor, _ := collection.Find(ctx, filter, options.Find().SetSort(bson.M{
		"list_order": 1,
	}))
	defer cusor.Close(ctx)
	var contentLists []editor.Course_content_infos
	cusor.All(ctx, &contentLists)

	for _, content := range contentLists {
		catalogs := getCatalogsByParentUuid(ctx, content.Uuid)
		if !helpers.Empty(catalogs) {
			for l, catalog := range catalogs {
				level := courseCode + "-L" + strconv.Itoa(l+1)
				chapterCatalogs := getCatalogsByParentUuid(ctx, catalog.Uuid)
				for c, chapterCatalog := range chapterCatalogs {
					chapterCode := level + "-C" + strconv.Itoa(c+1)
					moduleCatalogs := getCatalogsByParentUuid(ctx, chapterCatalog.Uuid)
					for m, moduleCatalog := range moduleCatalogs {
						moduleCode := chapterCode + "-M" + strconv.Itoa(m+1)
						partCatalogs := getCatalogsByParentUuid(ctx, moduleCatalog.Uuid)
						for p, _ := range partCatalogs {
							partCode := moduleCode + "-P" + strconv.Itoa(p+1)
							mgdb.OnlineClient.Database(courseContentDb).Collection("unlock_infos_part").UpdateOne(
								ctx,
								bson.M{
									"partCode": partCode,
									"userId":   param.UserId,
								},
								bson.M{
									"$set": bson.M{
										"courseCode":  courseCode,
										"chapter":     chapterCode,
										"correctRate": 100,
									},
								},
								options.Update().SetUpsert(true),
							)
						}
					}
					mgdb.OnlineClient.Database(courseContentDb).Collection("unlock_infos_chapter").UpdateOne(
						ctx,
						bson.M{
							"chapter": chapterCode,
							"user_id": param.UserId,
						},
						bson.M{
							"$set": bson.M{
								"course_code": courseCode,
								"level":       level,
							},
						},
						options.Update().SetUpsert(true),
					)
				}
			}
		}

		mgdb.OnlineClient.Database(KUYU).Collection("users_subscribe_course").UpdateOne(
			ctx,
			bson.M{
				"user_id":     param.UserId,
				"course_code": param.CourseCode,
			},
			bson.M{
				"$set": bson.M{
					"course_type":          courseInfos.CourseType,
					"purchase_time":        time.Now(),
					"del":                  0,
					"is_top":               0,
					"keeped":               false,
					"has_unlock_transfer":  true,
					"current_chapter_code": param.CourseCode + "-L1-C1",
				},
			},
			options.Update().SetUpsert(true),
		)
	}
}

type UnlockCatalogInfo struct {
	Uuid string `bson:"uuid" json:"uuid"`
}

func getCatalogsByParentUuid(ctx context.Context, parent_uuid string) (catalogs []UnlockCatalogInfo) {
	catalogsCollection := mgdb.OnlineClient.Database(courseContentDb).Collection(tblCatalogs)
	cusor, _ := catalogsCollection.Find(ctx, bson.M{
		"parent_uuid": parent_uuid,
		"has_del":     false,
		"is_show":     true,
	})
	defer cusor.Close(ctx)
	cusor.All(ctx, &catalogs)
	return
}

type UsersSubscribeCourse struct {
	UserId            string `bson:"user_id" json:"user_id"`
	HasUnlockTransfer bool   `bson:"has_unlock_transfer" json:"has_unlock_transfer"`
}
