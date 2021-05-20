package editorapi

import (
	"context"
	"editorApi/init/mgdb"
	"encoding/json"
	"fmt"
	chatroomProto "grpcSrv/proto/chatroom"
	imProto "grpcSrv/proto/im"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	client "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CourseOnlineMsg struct {
	UUID       string
	CourseCode string
	Title      string
}

type CourseOfflineMsg struct {
	UUID        string
	CourseCode  string
	Title       string
	FinishTitle string
	FinishInfo  string
	WeixinNo    string
}

//直播课程上线程序
func LiveCourseOnlineSub(msg *CourseOnlineMsg) {
	//设置recover，recover只能放在defer后面使用
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("==>%s\n", err)
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Printf("==> %s\n", string(buf[:n]))
		}
	}()
	ctx := context.Background()
	//创建聊天室
	fmt.Println("开始创建聊天室，uuid:", msg.UUID)
	rst := &chatroomProto.Request{
		Rooms: []*chatroomProto.Ctroom{
			&chatroomProto.Ctroom{
				Id:   msg.UUID,
				Name: msg.Title,
			},
		},
	}

	rsp, err := roomSrv.Create(context.Background(), rst)

	if rsp == nil || rsp.Code != 200 {

		fmt.Println("创建聊天室失败,uuid:", msg.UUID, err)

	}
	fmt.Println("结束创建聊天室，uuid:", msg.UUID)
	//获取课程信息

	subUids := []string{}
	subInfos := []map[string]string{}

	rsts, _ := toClient.Database(KUYU).Collection("users_subscribe_course").Find(
		ctx,
		bson.M{
			"course_code": msg.CourseCode,
		},
		options.Find().SetProjection(bson.M{
			"user_id": 1,
			"_id":     0,
		}),
	)

	rsts.All(ctx, &subInfos)
	rsts.Close(ctx)

	userExists := map[string]struct{}{}
	for _, sub := range subInfos {
		subUids = append(
			subUids,
			sub["user_id"],
		)
		userExists[sub["user_id"]] = struct{}{}
	}

	room := struct {
		ModuleName string `bson:"module_name"`
		LanCode    string `bson:"lan_code"`
		LiveInfo   struct {
			BasicCourseCode string `basic_course_code`
		} `liveInfo`
	}{}

	mgdb.FindOne(
		mgdb.EnvOnline,
		mgdb.DbKuyu,
		"course_module",
		bson.M{
			"code": msg.CourseCode,
		},
		nil,
		&room,
	)
	var lang *struct {
		Title map[string]string `bson:"title"`
	}
	mgdb.FindOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tblCourseLangs,
		bson.M{
			"lan_code": room.LanCode,
		},
		nil,
		&lang,
	)
	var lc *struct {
		StartTime int64 `bson:"startTime"`
	}

	mgdb.FindOne(
		mgdb.EnvOnline,
		mgdb.DbKuyu,
		tblLiveCourse,
		bson.M{
			"uuid": msg.UUID,
		},
		nil,
		&lc,
	)

	count, _ := mgdb.Count(
		mgdb.EnvOnline,
		mgdb.DbKuyu,
		tblLiveCourse,
		bson.M{
			"courseCode": msg.CourseCode,
			"isDel":      false,
			"startTime": bson.M{
				"$lt": lc.StartTime,
			},
		},
	)
	showTitle := "叮铃~上课时间到！您订阅的" + lang.Title["zh-CN"] + "直播课程《" + room.ModuleName + "》第" + strconv.Itoa(int(count+1)) + "课已经开始直播啦~点击此处直接进入教室→"
	extraData, _ := json.Marshal(map[string]interface{}{
		"courseCode": msg.CourseCode,
		"uuid":       msg.UUID,
		"remindTime": time.Now().Unix(),
		"state":      "start",
		"showTitle":  showTitle,
	})

	noticeUids := []string{}
	silentUids := []string{}

	if room.LiveInfo.BasicCourseCode == "" { //非配套直播课程
		noticeUids = subUids
		//进入直播课程的用户ID
		var inUserIds *struct {
			UserIds []string `bson:"userIds"`
		}
		st := toClient.Database(KUYU).Collection("liveCourseInUser").FindOne(
			ctx,
			bson.M{
				"courseCode": msg.CourseCode,
			},
		)

		st.Decode(&inUserIds)

		if inUserIds != nil {
			for _, uid := range inUserIds.UserIds {
				if _, ok := userExists[uid]; !ok {
					silentUids = append(silentUids, uid)
				}
			}
		}

	} else { //直播配套课程

		for _, v := range subUids {

			var current *struct {
				CurrentCourse *struct {
					Code string `bson:"code"`
				} `bson:"language_current_study"`
			}
			mgdb.FindOne(
				mgdb.EnvOnline,
				mgdb.DbKuyu,
				"users_course_setting",
				bson.M{
					"user_id": v,
				},
				nil,
				&current,
			)

			if current != nil && current.CurrentCourse != nil {
				if current.CurrentCourse.Code == room.LiveInfo.BasicCourseCode {
					noticeUids = append(noticeUids, v)
				} else {
					silentUids = append(silentUids, v)
				}
			} else {
				silentUids = append(silentUids, v)
			}

		}

		var inUserIds *struct {
			UserIds []string `bson:"userIds"`
		}
		st := toClient.Database(KUYU).Collection("liveCourseInUser").FindOne(
			ctx,
			bson.M{
				"courseCode": msg.CourseCode,
			},
		)
		st.Decode(&inUserIds)
		if inUserIds != nil {
			for _, uid := range inUserIds.UserIds {
				if _, ok := userExists[uid]; !ok {
					silentUids = append(silentUids, uid)
				}
			}
		}
	}
	fmt.Println("noticeUids", noticeUids)
	fmt.Println("silentUids", silentUids)

	wg := sync.WaitGroup{}
	noticeUids = filterUid("app:LiveRemind", noticeUids)
	go sendNotice(
		noticeUids,
		"app:LiveRemind",
		showTitle,
		showTitle,
		extraData,
		wg,
	)
	silentUids = filterUid("app:LiveRemindSilent", silentUids)
	go sendNotice(
		silentUids,
		"app:LiveRemindSilent",
		showTitle,
		"",
		extraData,
		wg,
	)

	wg.Wait()
}

//直播课程下线程序
func LiveCourseOfflineSub(msg *CourseOfflineMsg) {

	//设置recover，recover只能放在defer后面使用
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Printf("==> %s\n", string(buf[:n]))
		}
	}()

	ctx := context.Background()
	fmt.Println("开始销毁聊天室：", msg.UUID)
	rst := &chatroomProto.DistroyRequest{
		Ids: []string{
			msg.UUID,
		},
	}

	roomSrv.Distroy(ctx, rst)
	fmt.Println("成功销毁聊天室：", msg.UUID)
	noticeUids := []string{}
	silentUids := []string{}
	uidExists := map[string]struct{}{}

	rsts, _ := toClient.Database(KUYU).Collection("liveCourseEndNotice").Find(
		ctx,
		bson.M{
			"course_uuid": msg.UUID,
		},
		options.Find().SetProjection(bson.M{
			"user_id": 1,
			"_id":     0,
		}),
	)

	endNotices := []map[string]string{}
	rsts.All(ctx, &endNotices)
	rsts.Close(ctx)

	for _, sub := range endNotices {
		noticeUids = append(
			noticeUids,
			sub["user_id"],
		)
		uidExists[sub["user_id"]] = struct{}{}
	}

	//订阅课程消息

	subInfos := []map[string]string{}
	srsts, _ := toClient.Database(KUYU).Collection("users_subscribe_course").Find(
		ctx,
		bson.M{
			"course_code": msg.CourseCode,
		},

		options.Find().SetProjection(bson.M{
			"user_id": 1,
			"_id":     0,
		}),
	)

	srsts.All(ctx, &subInfos)
	srsts.Close(ctx)

	for _, sub := range subInfos {
		if _, ok := uidExists[sub["user_id"]]; !ok {
			silentUids = append(
				silentUids,
				sub["user_id"],
			)
			uidExists[sub["user_id"]] = struct{}{}
		}

	}
	var inUserIds *struct {
		UserIds []string `bson:"userIds"`
	}
	st := toClient.Database(KUYU).Collection("liveCourseInUser").FindOne(
		ctx,
		bson.M{
			"courseCode": msg.CourseCode,
		},
	)
	st.Decode(&inUserIds)
	if inUserIds != nil {
		for _, uid := range inUserIds.UserIds {
			if _, ok := uidExists[uid]; !ok {
				silentUids = append(silentUids, uid)
			}
		}
	}

	showTitle := "叮铃~下课啦！坚持学习的你太棒了，记得课后也要好好复习哟！"
	extraData, _ := json.Marshal(map[string]interface{}{
		"courseCode":  msg.CourseCode,
		"uuid":        msg.UUID,
		"remindTime":  time.Now().Unix(),
		"state":       "end",
		"finishTitle": msg.FinishTitle,
		"finishInfo":  msg.FinishInfo,
		"weixinNo":    msg.WeixinNo,
		"showTitle":   showTitle,
	})
	fmt.Println("noticeUids", noticeUids)
	fmt.Println("silentUids", silentUids)

	wg := sync.WaitGroup{}

	noticeUids = filterUid("app:LiveRemind", noticeUids)
	go sendNotice(
		noticeUids,
		"app:LiveRemind",
		showTitle,
		showTitle,
		extraData,
		wg,
	)

	silentUids = filterUid("app:LiveRemindSilent", silentUids)

	go sendNotice(
		silentUids,
		"app:LiveRemindSilent",
		showTitle,
		"",
		extraData,
		wg,
	)
	wg.Wait()
}

func sendNotice(
	uids []string,
	msgType,
	content string,
	pushContent string,
	extraData []byte,
	wg sync.WaitGroup,
) {
	fmt.Println("msgType:", msgType)
	fmt.Println("content:", content)
	fmt.Println("pushContent:", pushContent)
	fmt.Println("extraData:", string(extraData))

	wg.Add(1)
	defer wg.Done()

	if len(uids) == 0 {
		return
	}
	imSrv := imProto.NewMsgService(
		"go.micro.srv.talkmateSrv",
		client.NewClient(
			client.Registry(consul.NewRegistry(func(op *registry.Options) {
				op.Addrs = []string{
					"127.0.0.1:8500",
				}
			})),
		),
	)
	limit := 50
	start := 0
	end := limit
	for {
		tmpUids := []string{}
		if end > len(uids) {
			end = len(uids)
			tmpUids = uids[start:]
		} else {
			tmpUids = uids[start:end]
		}

		rsp, _ := imSrv.System(context.Background(), &imProto.Request{
			MsgType:     msgType,
			ToUids:      tmpUids,
			Content:     content,
			PushContent: pushContent,
			ExtraData:   string(extraData),
			PushData:    string(extraData),
		})
		if len(uids) == end {
			break
		}
		start = end
		end += limit

		if rsp == nil || rsp.Code != 200 {
			if rsp != nil && rsp.Code != 200 {
				fmt.Println("错误消息：", rsp.Msg)
			}

			fmt.Println("发送消息失败：", tmpUids)
		}

	}
}

func filterUid(msgType string, userIds []string) []string {
	//过滤掉不能发送消息的用户ID
	fmt.Println("filterUid start")
	msgToClient := map[string]map[string]string{
		"app:LiveRemind": map[string]string{
			"iOS":     "6.0.1",
			"ANDROID": "6.0.0",
		},
		"app:LiveRemindSilent": map[string]string{
			"iOS":     "6.0.6",
			"ANDROID": "6.0.4",
		},
	}

	versions, ok := msgToClient[msgType]
	if !ok {
		return nil
	}
	returnUids := []string{}
	agents := []*struct {
		Agent  string `bson:"agent"`
		UserId string `bson:"user_id"`
	}{}

	mgdb.Find(
		mgdb.EnvOnline,
		mgdb.DbKuyu,
		"app_open_state",
		bson.M{
			"user_id": bson.M{
				"$in": userIds,
			},
		},
		nil,
		nil,
		0,
		0,
		&agents,
	)
	uidExists := map[string]struct{}{}
	for _, a := range agents {
		if _, ok := uidExists[a.UserId]; !ok {
			tmpAgent := ""
			tmpVersion := ""

			if strings.Contains(a.Agent, "iOS") {
				tmpAgent = "iOS"
				tmpVersion = strings.TrimLeft(a.Agent, tmpAgent)
			} else if strings.Contains(a.Agent, "ANDROID") {
				tmpAgent = "ANDROID"
				tmpVersion = strings.Trim(strings.TrimLeft(a.Agent, tmpAgent), " ")
			}

			if tmpAgent != "" {
				msgVersion := versions[tmpAgent]
				if strings.Compare(tmpVersion, msgVersion) >= 0 {
					fmt.Println(a.UserId, tmpAgent, tmpVersion)
					returnUids = append(returnUids, a.UserId)
					uidExists[a.UserId] = struct{}{}
				}
			}
		}
	}
	fmt.Println("filterUid end")
	fmt.Println(returnUids)
	return returnUids
}
