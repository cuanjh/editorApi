package editorapi

import (
	"context"
	"crypto/md5"
	"editorApi/commons"
	"editorApi/config"
	"editorApi/controller/servers"
	"editorApi/init/initNats"
	"editorApi/init/mgdb"
	"editorApi/init/qmlog"
	"editorApi/mdbmodel/editor"
	"editorApi/middleware"
	"editorApi/requests"
	"editorApi/tools"
	"editorApi/tools/helpers"
	"editorApi/tools/utils"
	"encoding/hex"
	"fmt"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"encoding/json"

	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"

	"github.com/mongodb/mongo-go-driver/mongo"
	uuid "github.com/satori/go.uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"grpcSrv/proto/chatroom"
	chatroomProto "grpcSrv/proto/chatroom"
	imProto "grpcSrv/proto/im"

	"tkCommon/cmfunc"

	client "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

const (
	// LIVE_PUSH_DOMAIN     = "livepush.talkmate.com"
	// LIVE_PUSH_DOMAIN_KEY = "42e30e51fc72571a5d7010c3d4dc450c"
	// LIVE_PULL_DOMAIN     = "livepull.talkmate.com"
	// LIVE_PULL_DOMAIN_KEY = "42e30e51fc72571a5d7010c3d4dc450c"

	UPLOADFILE_DOMAIN = "https://uploadfile1.talkmate.com"
	UPLOADFILE_BUCKET = "uploadfiles"

	KUYU             = "kuyu"
	COURSE_TYPE_LIVE = 6
)

var tblLiveRoom string = "course_module"
var tblLiveCourse string = "liveCourse"
var tblUsers string = "users"
var tblUsersSubscribeCourse string = "users_subscribe_course"
var imSrv imProto.MsgService
var roomSrv chatroom.ChatroomService

func init() {

	imSrv = imProto.NewMsgService(
		"go.micro.srv.talkmateSrv",
		client.NewClient(
			client.Registry(consul.NewRegistry(func(op *registry.Options) {
				op.Addrs = []string{
					"127.0.0.1:8500",
				}
			})),
		),
	)

	roomSrv = chatroomProto.NewChatroomService(
		"go.micro.srv.talkmateSrv",
		client.NewClient(
			client.Registry(consul.NewRegistry(func(op *registry.Options) {
				op.Addrs = []string{
					"127.0.0.1:8500",
				}
			})),
		),
	)

}

type liveSubParam struct {
	TalkmateId string `json:"talkmate_id"`
	CourseCode string `json:"course_code"`
}

// @Tags LiveAPI????????????????????????
// @Summary ??????????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.liveSubParam true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"????????????"}"
// @Router /live/sub [post]
func LiveSub(ctx *gin.Context) {
	var liveParam *liveSubParam
	ctx.BindJSON(&liveParam)
	var user *struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	mgdb.FindOne(
		mgdb.EnvOnline,
		mgdb.DbKuyu,
		"users",
		bson.M{
			"talkmate_id": liveParam.TalkmateId,
		},
		bson.M{
			"_id": 1,
		},
		&user,
	)
	if user == nil {
		servers.ReportFormat(ctx, false, "???????????????", gin.H{})
		return
	}

	mgdb.UpdateOne(
		mgdb.EnvOnline,
		mgdb.DbKuyu,
		tblUsersSubscribeCourse,
		bson.M{
			"user_id":     user.ID.Hex(),
			"course_code": liveParam.CourseCode,
		},
		bson.M{
			"$set": bson.M{
				"purchase_time": time.Now().Unix(),
				"over_date":     0,
				"has_purchased": 0,
				"start_date":    0,
				"del":           0,
				"course_type":   6,
			},
		},
		true,
	)

	servers.ReportFormat(ctx, true, "????????????", gin.H{})
}

// @Tags LiveAPI????????????????????????
// @Summary ??????????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.modelListsParams true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"????????????"}"
// @Router /live/list [post]
func LiveList(ctx *gin.Context) {
	paras := modelListsParams{}
	ctx.BindJSON(&paras)
	var limit int64
	var skip int64

	limit = paras.PageSize
	if limit == 0 {
		limit = 40
	}
	skip = paras.PageNo * limit

	var (
		err   error
		cusor *mongo.Cursor
	)

	rooms := []*editor.LiveRoom{}

	roomsCollection := toClient.Database(KUYU).Collection(tblLiveRoom)

	var filter = bson.M{
		"course_type": COURSE_TYPE_LIVE,
		"is_del":      0,
	}

	claims, _ := middleware.GetClaims(ctx)
	if claims.AuthorityId != "1" && claims.AuthorityId != "13" {
		filter["liveInfo.live_user_uuid"] = claims.UUID.String()
	}

	if cusor, err = roomsCollection.Find(
		ctx,
		filter,
		options.Find().SetSort(bson.M{"created_on": -1}),
		options.Find().SetLimit(limit),
		options.Find().SetSkip(skip),
	); err != nil {
		checkErr(ctx, err)
		return
	}
	defer cusor.Close(ctx)
	cusor.All(ctx, &rooms)
	retData := []bson.M{}
	for _, r := range rooms {
		courses := []*editor.LiveCourse{}
		qmlog.QMLog.Info(r.Code)
		cusor, err := toClient.Database(KUYU).Collection(tblLiveCourse).Find(ctx, bson.M{
			"courseCode": r.Code,
			"isDel":      false,
		}, options.Find().SetSort(map[string]int{
			"listOrder": 1,
		}))
		if err != nil {
			checkErr(ctx, err)
			return
		}
		defer cusor.Close(ctx)
		cusor.All(ctx, &courses)
		retData = append(retData, bson.M{
			"room":    r,
			"courses": courses,
		})
	}

	servers.ReportFormat(ctx, true, "??????????????????", gin.H{
		"rooms": retData,
	})
}

type liveAddParam struct {
	Room    editor.LiveRoom     `json:"room"`
	Courses []editor.LiveCourse `json:"courses"`
}
type user struct {
	ID primitive.ObjectID `bson:"_id"`
}

// @Tags LiveAPI????????????????????????
// @Summary ????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.liveAddParam true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"????????????"}"
// @Router /live/add [post]
func LiveAdd(ctx *gin.Context) {
	var live liveAddParam
	ctx.BindJSON(&live)
	live.Room.Code = live.Room.LanCode + "-" + strconv.Itoa(int(time.Now().UnixNano()))
	live.Room.CreatedOn = time.Now()
	live.Room.UpdateTime = time.Now()
	live.Room.Published = "N"
	live.Room.IsDel = 0

	courses := make([]interface{}, len(live.Courses))
	for k, r := range live.Courses {
		r.CourseCode = live.Room.Code
		r.LanCode = live.Room.LanCode
		r.UUID = uuid.NewV4().String()
		r.Cover = live.Room.Cover
		r.IsDel = false
		courses[k] = r
	}
	userId := "54a2c128f8441bd93dd06647"
	if !helpers.Empty(live.Room.UserID) {
		userId = live.Room.UserID
	}
	majiaUsers := []string{}
	users := []*user{}
	cusor, _ := toClient.Database(KUYU).Collection(tblUsers).Find(ctx, bson.M{
		"role": "5050",
	})
	defer cusor.Close(ctx)
	cusor.All(ctx, &users)
	for _, u := range users {
		majiaUsers = append(
			majiaUsers,
			u.ID.Hex(),
		)
	}

	if len(majiaUsers) > 0 {
		i := rand.Intn(len(majiaUsers))
		userId = majiaUsers[i]
	}

	live.Room.UserID = userId
	live.Room.LiveInfo.CourseNum = len(courses)
	// ???????????????
	live.Room.LiveInfo.BaseRand = utils.GetBaseRand(live.Room.LiveInfo.Level)
	toClient.Database(KUYU).Collection(tblLiveRoom).InsertOne(ctx, live.Room)
	toClient.Database(KUYU).Collection(tblLiveCourse).InsertMany(ctx, courses)

	servers.ReportFormat(ctx, true, "????????????", gin.H{})
}

type liveEditParam struct {
	Room    editor.LiveEditRoom `json:"room"`
	Courses []editor.LiveCourse `json:"courses"`
}

// @Tags LiveAPI????????????????????????
// @Summary ????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.liveEditParam true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"????????????"}"
// @Router /live/edit [post]
func LiveEdit(ctx *gin.Context) {
	var live liveEditParam
	ctx.BindJSON(&live)
	live.Room.UpdateTime = time.Now()
	live.Room.LiveInfo.CourseNum = len(live.Courses)
	live.Room.LiveInfo.BaseRand = utils.GetBaseRand(live.Room.LiveInfo.Level)

	toClient.Database(KUYU).Collection(tblLiveRoom).UpdateOne(ctx, bson.M{
		"code": live.Room.Code,
	}, bson.M{
		"$set": bson.M{
			"cover_v2":                     live.Room.CoverV2,
			"cover":                        live.Room.Cover,
			"free_for_member":              live.Room.FreeForMember,
			"description":                  live.Room.Description,
			"tag_keys":                     live.Room.TagKeys,
			"money_type":                   live.Room.MoneyType,
			"money":                        live.Room.Money,
			"moneyDiscount":                live.Room.MoneyDiscount,
			"lan_code":                     live.Room.LanCode,
			"module_name":                  live.Room.ModuleName,
			"user_id":                      live.Room.UserId,
			"update_time":                  time.Now(),
			"liveInfo.tech_desc":           live.Room.LiveInfo.TechDesc,
			"liveInfo.level":               live.Room.LiveInfo.Level,
			"liveInfo.base_rand":           live.Room.LiveInfo.BaseRand,
			"liveInfo.tech_name":           live.Room.LiveInfo.TechName,
			"liveInfo.tech_photo":          live.Room.LiveInfo.TechPhoto,
			"liveInfo.startDate":           live.Room.LiveInfo.StartDate,
			"liveInfo.endDate":             live.Room.LiveInfo.EndDate,
			"liveInfo.startTime":           live.Room.LiveInfo.StartTime,
			"liveInfo.endTime":             live.Room.LiveInfo.EndTime,
			"liveInfo.weekDays":            live.Room.LiveInfo.WeekDays,
			"liveInfo.posters":             live.Room.LiveInfo.Posters,
			"liveInfo.courseNum":           len(live.Courses),
			"liveInfo.exclude_dates":       live.Room.LiveInfo.ExcludeDates,
			"liveInfo.videoUrl":            live.Room.LiveInfo.VideoUrl,
			"liveInfo.videoCoverUrl":       live.Room.LiveInfo.VideoCoverUrl,
			"liveInfo.finishTitle":         live.Room.LiveInfo.FinishTitle,
			"liveInfo.finishInfo":          live.Room.LiveInfo.FinishInfo,
			"liveInfo.weixinNo":            live.Room.LiveInfo.WeixinNo,
			"liveInfo.date_notice":         live.Room.LiveInfo.DateNotice,
			"liveInfo.basic_course_code":   live.Room.LiveInfo.BasicCourseCode,
			"liveInfo.basic_content_level": live.Room.LiveInfo.BasicContentLevel,
			"liveInfo.basic_chapter_cover": live.Room.LiveInfo.BasicChapterCover,
			"liveInfo.basic_profile_photo": live.Room.LiveInfo.BasicProfilePhoto, // ?????????????????????
			"liveInfo.live_user_uuid":      live.Room.LiveInfo.LiveUserUUID,      // ?????????????????????
			"liveInfo.dis_tech_photo":      live.Room.LiveInfo.DisTechPhoto,      // ???????????????????????????

		},
	})
	toClient.Database(KUYU).Collection(tblLiveCourse).UpdateMany(ctx, bson.M{
		"courseCode": live.Room.Code,
	}, bson.M{
		"$set": bson.M{
			"isDel": true,
		},
	})
	for _, r := range live.Courses {

		r.Cover = live.Room.Cover
		r.IsDel = false
		if r.UUID == "" {
			r.CourseCode = live.Room.Code
			r.LanCode = live.Room.LanCode
			r.UUID = uuid.NewV4().String()
			r.Cover = live.Room.Cover
			r.IsDel = false
		}
		toClient.Database(KUYU).Collection(tblLiveCourse).UpdateOne(
			ctx,
			bson.M{
				"uuid": r.UUID,
			},
			bson.M{
				"$set": bson.M{
					"listOrder":  r.ListOrder,
					"title":      r.Title,
					"lanCode":    r.LanCode,
					"courseCode": r.CourseCode,
					"date":       r.Date,
					"startTime":  r.StartTime,
					"endTime":    r.EndTime,
					"cover":      r.Cover,
					// "state":      r.State,
					"isDel": false,
				},
			},
			options.Update().SetUpsert(true),
		)
	}

	servers.ReportFormat(ctx, true, "????????????", gin.H{})
}

type liveDelPara struct {
	Code string `json:"code"`
}

// @Tags LiveAPI????????????????????????
// @Summary ????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.liveDelPara true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"????????????"}"
// @Router /live/del [post]
func LiveDel(ctx *gin.Context) {
	var param liveDelPara
	ctx.BindJSON(&param)
	toClient.Database(KUYU).Collection(tblLiveRoom).UpdateOne(ctx, bson.M{
		"code": param.Code,
	}, bson.M{
		"$set": bson.M{
			"is_del": 1,
		},
	})
	servers.ReportFormat(ctx, true, "????????????", gin.H{})
}

// @Tags LiveAPI????????????????????????
// @Summary ????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.liveDelPara true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"????????????"}"
// @Router /live/online [post]
func LiveOnline(ctx *gin.Context) {
	var param liveDelPara
	ctx.BindJSON(&param)

	toClient.Database(KUYU).Collection(tblLiveRoom).UpdateOne(ctx, bson.M{
		"code": param.Code,
	}, bson.M{
		"$set": bson.M{
			"published": "Y",
		},
	})
	servers.ReportFormat(ctx, true, "????????????", gin.H{})
}

// @Tags LiveAPI????????????????????????
// @Summary ????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.liveDelPara true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"????????????"}"
// @Router /live/offline [post]
func LiveOffline(ctx *gin.Context) {
	var param liveDelPara
	ctx.BindJSON(&param)
	toClient.Database(KUYU).Collection(tblLiveRoom).UpdateOne(ctx, bson.M{
		"code": param.Code,
	}, bson.M{
		"$set": bson.M{
			"published": "N",
		},
	})
	servers.ReportFormat(ctx, true, "????????????", gin.H{})
}

type liveCoursePara struct {
	UUID string `json:"uuid"` //??????UUID
}

// @Tags LiveAPI????????????????????????
// @Summary ????????????????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body liveCoursePara true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"??????"}"
// @Router /live/pushurl [post]
func LivePushUrl(ctx *gin.Context) {
	var param *liveCoursePara

	ctx.BindJSON(&param)
	expireTime := time.Now().AddDate(0, 1, 0).Unix()
	pushURL := GetLivePushUrl(param.UUID, expireTime)

	toClient.Database(KUYU).Collection(tblLiveCourse).UpdateOne(ctx, bson.M{
		"uuid": param.UUID,
	}, bson.M{
		"$set": bson.M{
			"livePushUrl": pushURL,
		},
	})
	servers.ReportFormat(ctx, true, "??????", gin.H{
		"livePushUrl": pushURL,
	})
}

// @Tags LiveAPI????????????????????????
// @Summary ????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body liveCoursePara true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"????????????"}"
// @Router /live/course/online [post]
func LiveCourseOnline(ctx *gin.Context) {
	var param liveCoursePara
	ctx.BindJSON(&param)

	var course editor.LiveCourse

	ret := toClient.Database(KUYU).Collection(tblLiveCourse).FindOne(ctx, bson.M{
		"uuid": param.UUID,
	})
	ret.Decode(&course)

	initNats.NatsConn.Publish("liveCourseOnlineSub",
		&CourseOnlineMsg{
			UUID:       course.UUID,
			CourseCode: course.CourseCode,
			Title:      course.Title,
		},
	)

	toClient.Database(KUYU).Collection(tblLiveCourse).UpdateOne(ctx, bson.M{
		"uuid": param.UUID,
	}, bson.M{
		"$set": bson.M{
			"state":         1,
			"realStartTime": time.Now().Unix(),
		},
	})
	toClient.Database(KUYU).Collection(tblLiveRoom).UpdateOne(ctx, bson.M{
		"code": course.CourseCode,
	}, bson.M{
		"$set": bson.M{
			"liveInfo.online":        true,
			"liveInfo.uuid":          course.UUID,
			"liveInfo.realStartTime": time.Now().Unix(),
		},
	})

	servers.ReportFormat(ctx, true, "????????????", gin.H{})
}

// @Tags LiveAPI????????????????????????
// @Summary ????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body liveCoursePara true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"????????????"}"
// @Router /live/course/offline [post]
func LiveCourseOffline(ctx *gin.Context) {
	var param liveCoursePara
	ctx.BindJSON(&param)

	var course editor.LiveCourse
	ret := toClient.Database(KUYU).Collection(tblLiveCourse).FindOneAndUpdate(ctx, bson.M{
		"uuid": param.UUID,
	}, bson.M{
		"$set": bson.M{
			"state":       -1,
			"realEndTime": time.Now().Unix(),
		},
	})
	ret.Decode(&course)
	var room editor.LiveRoom

	retRoom := toClient.Database(KUYU).Collection(tblLiveRoom).FindOneAndUpdate(ctx, bson.M{
		"code": course.CourseCode,
	}, bson.M{
		"$set": bson.M{
			"liveInfo.online":      false,
			"liveInfo.realEndTime": time.Now().Unix(),
		},
	})

	retRoom.Decode(&room)

	finishInfo := room.LiveInfo.FinishInfo
	wxNo := room.LiveInfo.WeixinNo
	finishTitle := room.LiveInfo.FinishTitle
	if course.FinishInfo != "" {
		finishInfo = course.FinishInfo
	}
	if course.WeixinNo != "" {
		wxNo = course.WeixinNo
	}
	if course.FinishTitle != "" {
		finishTitle = course.FinishTitle
	}

	initNats.NatsConn.Publish("liveCourseOfflineSub",
		&CourseOfflineMsg{
			UUID:        course.UUID,
			CourseCode:  course.CourseCode,
			Title:       course.Title,
			FinishInfo:  finishInfo,
			FinishTitle: finishTitle,
			WeixinNo:    wxNo,
		},
	)

	servers.ReportFormat(ctx, true, "????????????", gin.H{})
}

type liveCourseEditPara struct {
	UUID        string `json:"uuid"`        //??????UUID
	VideoUrl    string `json:"videoUrl"`    //??????????????????
	VideoCover  string `json:"videoCover"`  //??????????????????
	FinishInfo  string `json:"finishInfo"`  //????????????????????????
	FinishTitle string `json:"finishTitle"` //????????????????????????
	Date        string `json:"date"`        //????????????
	StartTime   int    `json:"startTime"`   //??????????????????
	EndTime     int    `json:"endTime"`     //??????????????????
	WeixinNo    string `json:"weixinNo"`    //?????????
	VideoTime   int64  `json:"videoTime"`   //??????????????????
}

// @Tags LiveAPI????????????????????????
// @Summary ??????????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body liveCourseEditPara true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"????????????"}"
// @Router /live/course/edit [post]
func LiveCourseEdit(ctx *gin.Context) {
	var param liveCourseEditPara
	ctx.BindJSON(&param)

	toClient.Database(KUYU).Collection(tblLiveCourse).UpdateOne(ctx, bson.M{
		"uuid": param.UUID,
	}, bson.M{
		"$set": bson.M{
			"videoUrl":    param.VideoUrl,
			"videoCover":  param.VideoCover,
			"finishTitle": param.FinishTitle,
			"finishInfo":  param.FinishInfo,
			"weixinNo":    param.WeixinNo,
			"date":        param.Date,
			"endTime":     param.EndTime,
			"startTime":   param.StartTime,
			"videoTime":   param.VideoTime,
		},
	})

	servers.ReportFormat(ctx, true, "??????", gin.H{})
}

// @Tags LiveAPI????????????????????????
// @Summary ??????????????????????????????
// @Security ApiKeyAuth
// @accept mpfd
// @Produce application/json
// @Param courseCode formData string true "????????????"
// @Param jumpUrl formData string true "?????????????????????"
// @Param shareTitle formData string true "????????????"
// @Param shareDesc formData string true "????????????"
// @Param sharePoster formData string true "???????????????????????????????????????,??????????????????"
// @Param qrCodeX formData int true "???????????? px"
// @Param qrCodeY formData int true "???????????? px"
// @Param scaleX formData int true "???????????????????????????????????????????????? px"
// @Param scaleY formData int true "???????????????????????????????????????????????? px"
// @Param bgImg formData file true "??????????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"??????"}"
// @Router /live/shareinfo [post]
func ShareInfo(ctx *gin.Context) {
	code := ctx.PostForm("courseCode")

	jumpUrl := ctx.PostForm("jumpUrl")

	shareTitle := ctx.PostForm("shareTitle")
	shareDesc := ctx.PostForm("shareDesc")

	qrCodeX, _ := strconv.Atoi(ctx.PostForm("qrCodeX"))
	qrCodeY, _ := strconv.Atoi(ctx.PostForm("qrCodeY"))
	scaleX, _ := strconv.Atoi(ctx.PostForm("scaleX"))
	scaleY, _ := strconv.Atoi(ctx.PostForm("scaleY"))

	sharePoster := ctx.PostFormArray("sharePoster")

	//??????????????????
	bgImg, err := ctx.FormFile("bgImg")
	set := bson.M{
		"liveInfo.shareTitle": shareTitle,
		"liveInfo.shareDesc":  shareDesc,
	}
	shareBgUrl := ""
	if err == nil {
		bgFilePath := "qrcode/bg/"
		bgFileName := code + bgImg.Filename
		ctx.SaveUploadedFile(bgImg, bgFilePath+bgFileName)
		//???????????????
		qrop := tools.NewQrCode(
			jumpUrl+"?courseCode="+code+"&t="+strconv.Itoa(int(time.Now().Unix())),
			qrCodeX,
			qrCodeY,
			qr.M,
			qr.Auto,
		)
		qrCodeFilePath := "qrcode/"
		qrCodeFileName, _, _ := qrop.Encode(qrCodeFilePath)

		//?????????????????????????????????????????????????????????
		mergeFilePath := "qrcode/merge/"
		mergeFileName := code + "shareImg.jpg"
		m := tools.NewMerge(
			bgFilePath,
			bgFileName,
			qrCodeFilePath,
			qrCodeFileName,
			mergeFilePath,
			scaleX,
			scaleY,
		)

		m.Generate(mergeFileName)
		//??????????????????
		_, filePath, _ := servers.UploadLocal(
			mergeFilePath+mergeFileName,
			UPLOADFILE_BUCKET,
			COURSE_UPLOADFILE_DOMAIN,
			"qrcode/merge/"+code+"shareImg.jpg",
		)
		shareBgUrl = filePath + "?t=" + strconv.Itoa(int(time.Now().Unix()))

		set["liveInfo.shareBgUrl"] = shareBgUrl

		os.Remove(bgFilePath + bgFileName)
		os.Remove(qrCodeFilePath + qrCodeFileName)
		os.Remove(mergeFilePath + mergeFileName)
	}

	if sharePoster != nil && len(sharePoster) > 0 {
		set["liveInfo.sharePoster"] = sharePoster
	}
	//???????????????????????????????????????
	toClient.Database(KUYU).Collection(tblLiveRoom).UpdateOne(ctx, bson.M{
		"code": code,
	}, bson.M{
		"$set": set,
	})

	servers.ReportFormat(ctx, true, "??????", gin.H{
		"shareBgUrl": shareBgUrl,
	})
}
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {

	}
	return strings.Replace(dir, "\\", "/", -1)
}

type roomMsg struct {
	UUID    string `json:"uuid"`    //??????UUID
	Content string `json:"content"` //????????????
	UserId  string `json:"user_id"` //???????????????ID
	Role    string `json:"role"`    //??????
	Coins   int    `json:"coins"`   //???????????????
}

// @Tags LiveAPI????????????????????????
// @Summary ?????????????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body roomMsg true "????????????,UUID??????????????????UUID??????Role(?????????teacher,student),Coins ????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"????????????"}"
// @Router /live/chatroom/sendmsg [post]
func ChatroomSendMsg(ctx *gin.Context) {
	var msg roomMsg
	ctx.BindJSON(&msg)

	var user *struct {
		Nickname string `bson:"nickname"`
		Photo    string `bson:"photo"`
	}
	_id, _ := primitive.ObjectIDFromHex(msg.UserId)
	r := toClient.Database(KUYU).Collection(tblUsers).FindOne(
		ctx,
		bson.M{
			"_id": _id,
		},
	)
	r.Decode(&user)

	extraMap := map[string]interface{}{
		"role":        msg.Role,
		"userName":    user.Nickname,
		"photo":       UPLOADFILE_DOMAIN + "/" + user.Photo,
		"isMember":    true,
		"coins":       msg.Coins,
		"commentTime": time.Now(),
	}
	byts, _ := json.Marshal(&extraMap)

	rsp, _ := imSrv.ChatroomCustome(ctx, &imProto.Request{
		MsgType:    "app:ChatroomMsgv1",
		FromUserId: msg.UserId,
		ToUids:     []string{msg.UUID},
		Content:    msg.Content,
		ExtraData:  string(byts),
	})
	if rsp.Code == 200 {
		toClient.Database(KUYU).Collection("liveCourseComment").InsertOne(
			ctx,
			bson.M{
				"user_id":      msg.UserId,
				"isMember":     true,
				"comment_time": time.Now().UnixNano(),
				"comment":      msg.Content,
				"coins":        msg.Coins,
				"course_uuid":  msg.UUID,
				"role":         msg.Role,
				"isMajia":      true,
				"isGag":        false,
			},
		)
		servers.ReportFormat(ctx, true, "????????????", gin.H{})
	} else {
		servers.ReportFormat(ctx, true, "????????????"+rsp.Msg, gin.H{
			"err": rsp.Msg,
		})
	}
}

type commentsParam struct {
	UUID string `json:"uuid"` //????????????UUID
}

// @Tags LiveAPI????????????????????????
// @Summary ???????????????????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body commentsParam true "????????????UUID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"????????????"}"
// @Router /live/chatroom/comments [post]
func ChatroomComments(ctx *gin.Context) {
	var param commentsParam
	ctx.BindJSON(&param)

	var comments []*struct {
		UserId      string `bson:"user_id" json:"user_id"`
		Nickname    string `json:"nickname"`
		Photo       string `json:"photo"`
		Comment     string `bson:"comment" json:"comment"`
		CommentTime int64  `bson:"comment_time" json:"comment_time"`
		Coins       int64  `bson:"coins" json:"coins"`
		Role        string `bson:"role" json:"role"`
		IsGag       bool   `bson:"isGag" json:"isGag"` // ??????????????????
		IsMajia     bool   `bson:"isMajia" json:"isMajia"`
	}

	cusor, _ := toClient.Database(KUYU).Collection("liveCourseComment").Find(ctx, bson.M{
		"course_uuid": param.UUID,
	})
	defer cusor.Close(ctx)
	cusor.All(ctx, &comments)

	for _, c := range comments {
		user := &struct {
			Nickname string `bson:"nickname"`
			Photo    string `bson:"photo"`
		}{}
		_id, _ := primitive.ObjectIDFromHex(c.UserId)
		r := toClient.Database(KUYU).Collection(tblUsers).FindOne(ctx, bson.M{
			"_id": _id,
		})
		r.Decode(&user)
		c.Nickname = user.Nickname
		c.Photo = cmfunc.Photo(user.Photo, "")
	}
	servers.ReportFormat(ctx, true, "????????????", gin.H{
		"comments": comments,
	})
}

// @Tags LiveAPI????????????????????????
// @Summary ????????????????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"??????"}"
// @Router /live/chatroom/majia [get]
func ChatroomMajia(ctx *gin.Context) {
	var users []*struct {
		ID       primitive.ObjectID `bson:"_id"`
		Photo    string             `bson:"photo"`
		Nickname string             `bson:"nickname"`
	}
	cusor, _ := toClient.Database(KUYU).Collection(tblUsers).Find(ctx, bson.M{"role": "5050"})
	defer cusor.Close(ctx)
	cusor.All(ctx, &users)
	userInfos := []bson.M{}
	for _, u := range users {
		userInfos = append(userInfos, bson.M{
			"user_id":  u.ID.Hex(),
			"nickname": u.Nickname,
			"photo":    "https://uploadfile1.talkmate.com/" + u.Photo,
		})
	}
	servers.ReportFormat(ctx, true, "??????????????????", gin.H{
		"userInfos": userInfos,
	})
}

// @Tags LiveAPI????????????????????????
// @Summary ????????????token
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"??????"}"
// @Router /live/wxtoken [get]
func WXToken(ctx *gin.Context) {

	ticket := ""
	cmfunc.CacheGetV2("weixinShareTicketv1", &ticket)
	if ticket == "" {
		url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=wx32c126b96bed2cbc&secret=cdcf4602af8dcf67683f612058442324"
		rps, err := http.Get(url)
		defer rps.Body.Close()
		if err != nil {
			ctx.JSONP(http.StatusOK, gin.H{
				"success": true,
				"msg":     err,
				"data":    "",
			})
		} else {
			body, _ := ioutil.ReadAll(rps.Body)

			var tokenInfo *struct {
				AccessToken string `json:"access_token"`
				ExpiresIn   string `json:"expires_in"`
			}
			json.Unmarshal(body, &tokenInfo)

			if tokenInfo != nil {
				url := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=" + tokenInfo.AccessToken + "&type=jsapi"
				rspTik, _ := http.Get(url)
				defer rspTik.Body.Close()
				body, _ = ioutil.ReadAll(rspTik.Body)
				var ticketInfo *struct {
					Ticket    string `json:"ticket"`
					ExpiresIn string `json:"expires_in"`
				}
				json.Unmarshal(body, &ticketInfo)
				if ticketInfo != nil {
					ticket = ticketInfo.Ticket
				}
			}

		}
		cmfunc.CacheSaveV2("weixinShareTicketv1", &ticket, time.Minute*59)
	}
	ctx.JSONP(http.StatusOK, gin.H{
		"success": true,
		"msg":     "",
		"data":    ticket,
	})
}

//??????????????????
func GetLivePushUrl(streamName string, txTime int64) string {

	txTimeStr := fmt.Sprintf("%X", txTime)
	lcf := config.GinVueAdminconfig.LiveCourseConfig
	txt := lcf.LivePushDomainKey + streamName + txTimeStr
	txtBytes := md5.Sum([]byte(txt))
	txSecret := hex.EncodeToString(txtBytes[0:])

	return "rtmp://" + lcf.LivePushDomain + "/live/" + streamName + "?txSecret=" + txSecret + "&txTime=" + txTimeStr
}

//??????????????????
func GetLivePullUrl(streamName string, txTime int64) string {

	txTimeStr := fmt.Sprintf("%X", txTime)
	lcf := config.GinVueAdminconfig.LiveCourseConfig
	txt := lcf.LivePullDomainKey + streamName + txTimeStr
	txtBytes := md5.Sum([]byte(txt))
	txSecret := hex.EncodeToString(txtBytes[0:])

	return "http://" + lcf.LivePullDomain + "/live/" + streamName + ".flv?txSecret=" + txSecret + "&txTime=" + txTimeStr
}

type GagParam struct {
	UserIds    []string `json:"userIds"`    //?????????????????????ID
	ChatroomId string   `json:"courseUUID"` //??????ChatroomID
	Minute     string   `json:"minute"`     //???????????????
}

// @Tags LiveAPI????????????????????????
// @Summary ???????????????????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body GagParam true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"??????"}"
// @Router /live/gagadd [post]
func GagAdd(ctx *gin.Context) {
	var param GagParam
	ctx.BindJSON(&param)

	if helpers.Empty(param.ChatroomId) && helpers.Empty(param.UserIds) {
		servers.ReportFormat(ctx, false, "????????????????????????", gin.H{})
		return
	}

	roomSrv := chatroomProto.NewChatroomService(
		"go.micro.srv.talkmateSrv",
		client.NewClient(
			client.Registry(consul.NewRegistry(func(op *registry.Options) {
				op.Addrs = []string{
					"127.0.0.1:8500",
				}
			})),
		),
	)
	rst := &chatroomProto.GagAddRequest{
		UserIds:    param.UserIds,
		ChatroomId: param.ChatroomId,
		Minute:     60,
	}

	rsp, _ := roomSrv.GagAdd(context.Background(), rst)
	if rsp.Code != 200 {
		servers.ReportFormat(ctx, false, "????????????????????????????????????"+rsp.Msg, gin.H{})
	} else {
		where := bson.M{
			"user_id":     bson.M{"$in": param.UserIds},
			"course_uuid": param.ChatroomId,
		}

		updata := bson.M{
			"$set": bson.M{
				"isGag": true,
			},
		}

		toClient.Database(KUYU).Collection("liveCourseComment").UpdateMany(ctx, where, updata)

		servers.ReportFormat(ctx, true, "?????????????????????????????????", gin.H{})
	}
}

// @Tags LiveAPI????????????????????????
// @Summary ???????????????????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body GagParam true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"??????"}"
// @Router /live/gagremove [post]
func GagRemove(ctx *gin.Context) {
	var param GagParam
	ctx.BindJSON(&param)

	if helpers.Empty(param.ChatroomId) && helpers.Empty(param.UserIds) {
		servers.ReportFormat(ctx, false, "????????????????????????", gin.H{})
		return
	}

	roomSrv := chatroomProto.NewChatroomService(
		"go.micro.srv.talkmateSrv",
		client.NewClient(
			client.Registry(consul.NewRegistry(func(op *registry.Options) {
				op.Addrs = []string{
					"127.0.0.1:8500",
				}
			})),
		),
	)
	rst := &chatroomProto.GagRemoveRequest{
		UserIds:    param.UserIds,
		ChatroomId: param.ChatroomId,
	}

	rsp, _ := roomSrv.GagRemove(context.Background(), rst)
	if rsp.Code != 200 {
		servers.ReportFormat(ctx, false, "????????????????????????????????????"+rsp.Msg, gin.H{})
	} else {

		where := bson.M{
			"user_id":     bson.M{"$in": param.UserIds},
			"course_uuid": param.ChatroomId,
		}

		updata := bson.M{
			"$set": bson.M{
				"isGag": false,
			},
		}

		toClient.Database(KUYU).Collection("liveCourseComment").UpdateMany(ctx, where, updata)

		servers.ReportFormat(ctx, true, "?????????????????????????????????", gin.H{})
	}
}

// @Tags LiveAPI????????????????????????
// @Summary ???????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"??????"}"
// @Router /live/dbtest [post]
func DbTest(ctx *gin.Context) {
	var rst []*editor.Catalogs
	mgdb.Find(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		"catalogs",
		bson.M{"uuid": "b98cb234-d8d9-4ed4-8a32-1582f274d157"},
		nil,
		nil,
		0, 0, &rst,
	)
	servers.ReportFormat(ctx, true, "??????", gin.H{"rst": rst})
}

type UserCountParam struct {
	CourseUuid string `bson:"course_uuid" json:"course_uuid"` //??????ChatroomID
}

// @Tags LiveAPI????????????????????????
// @Summary ??????????????????????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body UserCountParam true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"??????"}"
// @Router /live/usercount [post]
func UserCount(ctx *gin.Context) {
	var param UserCountParam
	ctx.BindJSON(&param)
	if helpers.Empty(param.CourseUuid) {
		servers.ReportFormat(ctx, false, "????????????????????????", gin.H{})
		return
	}
	var result []*editor.LiveCourseUserCount
	cusor, _ := toClient.Database(KUYU).Collection("live_course_user_count").Find(ctx, bson.M{
		"course_uuid": param.CourseUuid,
	}, options.Find().SetSort(map[string]int{
		"created_on": 1,
	}))
	defer cusor.Close(ctx)
	cusor.All(ctx, &result)

	servers.ReportFormat(ctx, true, "??????", gin.H{"result": result})
}

type ListOrderParam struct {
	Code      string `json:"code"`
	ListOrder int    `json:"list_order"`
}

// @Tags LiveAPI????????????????????????
// @Summary ?????????????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ListOrderParam true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"??????"}"
// @Router /live/edit_list_order [post]
func EditListOrder(ctx *gin.Context) {
	var param ListOrderParam
	ctx.BindJSON(&param)
	if helpers.Empty(param.Code) || helpers.Empty(param.ListOrder) {
		servers.ReportFormat(ctx, false, "????????????????????????", gin.H{})
		return
	}

	toClient.Database(KUYU).Collection(tblLiveRoom).UpdateOne(ctx, bson.M{
		"code": param.Code,
	}, bson.M{
		"$set": bson.M{
			"list_order": param.ListOrder,
		},
	})

	servers.ReportFormat(ctx, true, "????????????", nil)
}

// @Tags LiveAPI????????????????????????
// @Summary ????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ListOrderParam true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"??????"}"
// @Router /live/comments_upload [post]
func CommentsUpload(ctx *gin.Context) {
	header, err := ctx.FormFile("filename")
	courseUuid := ctx.PostForm("course_uuid")
	if err != nil {
		commons.Error(ctx, 500, err, "????????????????????????")
	}

	os.MkdirAll("data/comments/", os.ModePerm)
	dst := "data/comments/" + uuid.NewV4().String() + ".xlsx"

	// gin ??????????????????,??????????????????
	if err := ctx.SaveUploadedFile(header, dst); err != nil {
		commons.Error(ctx, 500, err, "?????????????????????")
	}

	// ????????????
	initNats.NatsConn.Publish("CommentsUpload",
		&requests.CommentsUploadHanderParams{
			FilePath:   dst,
			CourseUuid: courseUuid,
		},
	)

	commons.Success(ctx, nil, "?????????", nil)
}

// @Tags LiveAPI????????????????????????
// @Summary ??????????????????
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ListOrderParam true "????????????"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"??????"}"
// @Router /live/send_live_common [post]
func SendLiveCommon(ctx *gin.Context) {
	liveCommentsTmp := []LiveCommentsTmp{}
	var filter = bson.D{}
	filter = append(filter, bson.E{"create_time", bson.M{"$lt": time.Now().UnixNano()}})
	filter = append(filter, bson.E{"status", 1})
	cusor, _ := toClient.Database(KUYU).Collection("live_comments_tmp").Find(ctx, filter)
	cusor.All(ctx, &liveCommentsTmp)
	defer cusor.Close(ctx)

	for _, liveComment := range liveCommentsTmp {
		userId := GetRandomUserId()
		var user *struct {
			Nickname string `bson:"nickname"`
			Photo    string `bson:"photo"`
		}
		_id, _ := primitive.ObjectIDFromHex(userId)
		r := toClient.Database(KUYU).Collection(tblUsers).FindOne(
			ctx,
			bson.M{
				"_id": _id,
			},
		)
		r.Decode(&user)

		extraMap := map[string]interface{}{
			"role":        "student",
			"userName":    user.Nickname,
			"photo":       UPLOADFILE_DOMAIN + "/" + user.Photo,
			"isMember":    true,
			"coins":       0,
			"commentTime": time.Now(),
		}
		byts, _ := json.Marshal(&extraMap)

		imSrv.ChatroomCustome(ctx, &imProto.Request{
			MsgType:    "app:ChatroomMsgv1",
			FromUserId: userId,
			ToUids:     []string{liveComment.CourseUuid},
			Content:    liveComment.Content,
			ExtraData:  string(byts),
		})
		//if rsp.Code == 200 {
		//
		//}
		toClient.Database(KUYU).Collection("liveCourseComment").InsertOne(
			ctx,
			bson.M{
				"user_id":      userId,
				"isMember":     true,
				"comment_time": liveComment.CreateTime,
				"comment":      liveComment.Content,
				"coins":        0,
				"course_uuid":  liveComment.CourseUuid,
				"role":         "student",
				"isMajia":      true,
				"isGag":        false,
			},
		)

		toClient.Database(KUYU).Collection("live_comments_tmp").UpdateOne(ctx, bson.M{
			"uuid": liveComment.Uuid,
		}, bson.M{
			"$set": bson.M{
				"status": 2,
			},
		})
	}

	commons.Success(ctx, nil, "?????????", nil)
}

type LiveCommentsTmp struct {
	Uuid       string `bson:"uuid" json:"uuid"`               //??????course_uuid
	CourseUuid string `bson:"course_uuid" json:"course_uuid"` //??????course_uuid
	Content    string `bson:"content" json:"content"`         //??????content
	CreateTime int64  `bson:"create_time" json:"create_time"` //????????????
	Status     int    `bson:"status" json:"status"`           //????????????
}

func CommentsUploadHander(request *requests.CommentsUploadHanderParams) {
	xlFile, err := xlsx.OpenFile(request.FilePath)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	var course editor.LiveCourse
	ret := toClient.Database(KUYU).Collection(tblLiveCourse).FindOne(ctx, bson.M{
		"uuid": request.CourseUuid,
	})
	ret.Decode(&course)

	var lock = sync.Mutex{}
	for _, sheet := range xlFile.Sheets {
		for index, row := range sheet.Rows {
			if index > 0 {
				lock.Lock()
				var liveCommentsTmp LiveCommentsTmp
				for j, cell := range row.Cells {
					if j == 0 {
						times, _ := strconv.ParseInt(cell.Value, 10, 64)
						liveCommentsTmp.CreateTime = times*1000000000 + time.Now().UnixNano()
					}
					if j == 1 {
						liveCommentsTmp.Content = cell.Value
					}
				}
				liveCommentsTmp.CourseUuid = request.CourseUuid
				liveCommentsTmp.Status = 1
				liveCommentsTmp.Uuid = uuid.NewV4().String()
				toClient.Database(KUYU).Collection("live_comments_tmp").InsertOne(context.TODO(), liveCommentsTmp)
				lock.Unlock()
			}
		}
	}
}

func GetRandomUserId() string {
	ctx := context.Background()
	majiaUsers := []string{}
	users := []*user{}
	cusor, _ := toClient.Database(KUYU).Collection(tblUsers).Find(ctx, bson.M{
		"role": "5050",
	})
	cusor.All(ctx, &users)
	defer cusor.Close(ctx)
	for _, u := range users {
		majiaUsers = append(majiaUsers, u.ID.Hex())
	}
	return majiaUsers[rand.Intn(len(majiaUsers))]
}
