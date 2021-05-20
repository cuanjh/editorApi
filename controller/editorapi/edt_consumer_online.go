package editorapi

import (
	"editorApi/init/mgdb"
	"editorApi/init/qmlog"
	"editorApi/mdbmodel/editor"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

var toClient *mongo.Client
var (
	onLineJobsCollection     string
	catalogCollection        string
	contentVersionCollection string
	courseLangsCollection    string
	courseInfoCollection     string
	contentTagsCollection    string
)

func init() {
	toClient = mgdb.OnlineClient
	onLineJobsCollection = "online_jobs"
	catalogCollection = "catalogs"
	contentVersionCollection = "course_content_infos"
	courseLangsCollection = "course_langs"
	courseInfoCollection = "course_infos"
	contentTagsCollection = "content_tags"
}

//用于复制目录以及目录下面的内容
func CopyContentVersion(msg *CatalogCopyParam) {
	//设置recover，recover只能放在defer后面使用
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Printf("==> %s\n", string(buf[:n]))
		}
	}()
	qmlog.QMLog.Info("开始复制版本内容", msg.Uuids[0])

	var catalogs []*editor.Catalogs
	mgdb.Find(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tblCatalogs,
		bson.M{
			"parent_uuid": msg.Uuids[0],
			"has_del":     false,
		},
		nil,
		nil,
		0,
		0,
		&catalogs,
	)

	uuids := []string{}
	for _, c := range catalogs {
		uuids = append(uuids, c.Uuid)
	}
	if len(uuids) > 0 {
		msg.Uuids = uuids
		copyCatalogInfo(msg, false)
	}
}

func PushContent(msg *PushOnlineMsg) {

	if msg.DbEnv == "online" {
		PushContentOnline(msg)
	} else if msg.DbEnv == "test" {
		PushContentTest(msg)
	}

}

//上线课程内容到正式环境
func PushContentOnline(msg *PushOnlineMsg) {
	//设置recover，recover只能放在defer后面使用
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Printf("==> %s\n", string(buf[:n]))
		}
	}()
	qmlog.QMLog.Info("开始上线")

	//更新任务状态

	mgdb.UpdateOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		onLineJobsCollection,
		bson.M{
			"uuid": msg.UUID,
		},
		bson.M{
			"$set": bson.M{
				"state": 1,
			},
		},
		false,
	)
	toConfig := mgdb.EnvOnline
	contentVersionUuid := ""
	if msg.OnlineType == "catalog" {
		//当前目录以及上层目录信息到线上

		catatlogUuid := msg.OnlineUUID
		qmlog.QMLog.Info("开始上线父目录")
		for {
			var ct *editor.Catalogs
			mgdb.FindOne(
				mgdb.EnvEditor,
				EDITOR_DB,
				catalogCollection,
				bson.M{
					"uuid": catatlogUuid,
				},
				nil,
				&ct,
			)

			if ct != nil {
				mgdb.UpdateOne(
					toConfig,
					mgdb.DbContent,
					catalogCollection,
					bson.M{
						"uuid": catatlogUuid,
					},
					bson.M{
						"$set": ct,
					},
					true,
				)

				//父目录的has_changed不更新
				if catatlogUuid == msg.OnlineUUID {
					mgdb.UpdateOne(
						mgdb.EnvEditor,
						mgdb.DbEditor,
						catalogCollection,
						bson.M{
							"uuid": catatlogUuid,
						},
						bson.M{
							"$set": bson.M{
								"has_changed": false,
							},
						},
						false,
					)
				}
				catatlogUuid = ct.Parent_uuid
				qmlog.QMLog.Info("父目录UUID：" + catatlogUuid)
			} else {
				contentVersionUuid = catatlogUuid
				break
			}
		}
		qmlog.QMLog.Info("结束上线父目录")
		//更新内容版本信息

		updateContentVersion(
			contentVersionUuid,
			true,
			toConfig,
		)

		//更新孩子目录以及内容
		updateChildeCatalogs(
			msg.OnlineUUID,
			toConfig,
		)
	} else if msg.OnlineType == "content_version" {
		//更新版本、语言和课程信息
		contentVersionUuid = msg.OnlineUUID
		updateContentVersion(
			msg.OnlineUUID,
			false,
			toConfig,
		)

		//上线子目录以及内容
		updateChildeCatalogs(
			msg.OnlineUUID,
			toConfig,
		)
	}

	//上线课程标签
	qmlog.QMLog.Info("开始上线课程标签")
	var tags []*editor.Content_tags

	mgdb.Find(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		contentTagsCollection,
		nil,
		nil,
		nil,
		0,
		0,
		&tags,
	)

	for _, t := range tags {
		qmlog.QMLog.Info(t.Key)
		mgdb.UpdateOne(
			toConfig,
			mgdb.DbContent,
			contentTagsCollection,
			bson.M{
				"key": t.Key,
			},
			bson.M{
				"$set": t,
			},
			true,
		)
	}
	qmlog.QMLog.Info("结束上线课程标签")
	mgdb.UpdateOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		"online_jobs",
		bson.M{
			"uuid": msg.UUID,
		},
		bson.M{
			"$set": bson.M{
				"state": 2,
			},
		},
		false,
	)

	var version *editor.Course_content_infos
	mgdb.FindOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		contentVersionCollection,
		bson.M{"uuid": contentVersionUuid},
		nil,
		&version,
	)

	mgdb.UpdateOne(
		toConfig,
		mgdb.DbContent,
		contentVersionCollection,
		bson.M{
			"uuid": contentVersionUuid,
		},
		bson.M{
			"$set": version,
		},
		true,
	)
	qmlog.QMLog.Info("结束上线")
	fmt.Println(msg)
}

//上线课程内容到测试环境
func PushContentTest(msg *PushOnlineMsg) {
	//设置recover，recover只能放在defer后面使用
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Printf("==> %s\n", string(buf[:n]))
		}
	}()
	qmlog.QMLog.Info("开始上线")
	contentVersionUuid := ""
	//更新任务状态
	toConfig := mgdb.EnvTest
	if msg.OnlineType == "catalog" {
		//当前目录以及上层目录信息到线上

		catatlogUuid := msg.OnlineUUID
		qmlog.QMLog.Info("开始上线父目录")
		for {
			var ct *editor.Catalogs
			mgdb.FindOne(
				mgdb.EnvEditor,
				EDITOR_DB,
				catalogCollection,
				bson.M{
					"uuid": catatlogUuid,
				},
				nil,
				&ct,
			)

			if ct != nil {
				mgdb.UpdateOne(
					toConfig,
					mgdb.DbContent,
					catalogCollection,
					bson.M{
						"uuid": catatlogUuid,
					},
					bson.M{
						"$set": ct,
					},
					true,
				)

				catatlogUuid = ct.Parent_uuid
				qmlog.QMLog.Info("父目录UUID：" + catatlogUuid)
			} else {
				contentVersionUuid = catatlogUuid
				break
			}
		}
		qmlog.QMLog.Info("结束上线父目录")
		//更新内容版本信息

		updateContentVersion(
			contentVersionUuid,
			true,
			toConfig,
		)

		//更新孩子目录以及内容
		updateChildeCatalogs(
			msg.OnlineUUID,
			toConfig,
		)
	} else if msg.OnlineType == "content_version" {
		//更新版本、语言和课程信息
		contentVersionUuid = msg.OnlineUUID

		updateContentVersion(
			msg.OnlineUUID,
			false,
			toConfig,
		)
		//上线子目录以及内容
		updateChildeCatalogs(
			msg.OnlineUUID,
			toConfig,
		)
	}

	//上线课程标签
	qmlog.QMLog.Info("开始上线课程标签")
	var tags []*editor.Content_tags

	mgdb.Find(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		contentTagsCollection,
		nil,
		nil,
		nil,
		0,
		0,
		&tags,
	)

	for _, t := range tags {
		qmlog.QMLog.Info(t.Key)
		mgdb.UpdateOne(
			toConfig,
			mgdb.DbContent,
			contentTagsCollection,
			bson.M{
				"key": t.Key,
			},
			bson.M{
				"$set": t,
			},
			true,
		)
	}

	var version *editor.Course_content_infos
	mgdb.FindOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		contentVersionCollection,
		bson.M{"uuid": contentVersionUuid},
		nil,
		&version,
	)

	mgdb.UpdateOne(
		toConfig,
		mgdb.DbContent,
		contentVersionCollection,
		bson.M{
			"uuid": contentVersionUuid,
		},
		bson.M{
			"$set": version,
		},
		true,
	)
	qmlog.QMLog.Info("结束上线")
	fmt.Println(msg)
}

func updateContentVersion(
	contentVersionUUID string,
	isParent bool,
	toConfig mgdb.EnvConfig,
) {
	//更新内容版本信息
	var version *editor.Course_content_infos

	mgdb.FindOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		contentVersionCollection,
		bson.M{"uuid": contentVersionUUID},
		nil,
		&version,
	)

	if version != nil {
		qmlog.QMLog.Info("上线内容版本信息:" + version.Code)

		// mgdb.UpdateOne(
		// 	toConfig,
		// 	mgdb.DbContent,
		// 	contentVersionCollection,
		// 	bson.M{
		// 		"uuid": version.Uuid,
		// 	},
		// 	bson.M{
		// 		"$set": version,
		// 	},
		// 	true,
		// )

		if !isParent && toConfig == mgdb.EnvOnline {
			mgdb.UpdateOne(
				mgdb.EnvEditor,
				mgdb.DbEditor,
				contentVersionCollection,
				bson.M{
					"uuid": version.Uuid,
				},
				bson.M{
					"$set": bson.M{
						"has_changed": false,
					},
				},
				false,
			)
		}
		//上线课程信息
		var course bson.M
		mgdb.FindOne(
			mgdb.EnvEditor,
			mgdb.DbEditor,
			courseInfoCollection,
			bson.M{
				"uuid": version.Parent_uuid,
			},
			nil,
			&course,
		)

		if course != nil {
			mgdb.UpdateOne(
				toConfig,
				mgdb.DbContent,
				courseInfoCollection,
				bson.M{
					"uuid": course["uuid"],
				},
				bson.M{
					"$set": course,
				},
				true,
			)

		}

		//更新语言种类
		var lang *editor.Course_langs
		mgdb.FindOne(
			mgdb.EnvEditor,
			mgdb.DbEditor,
			courseLangsCollection,
			bson.M{"lan_code": course["lan_code"]},
			nil,
			&lang,
		)

		if lang != nil {
			qmlog.QMLog.Info("上线语言种类信息：" + lang.Lan_code)
			mgdb.UpdateOne(
				toConfig,
				mgdb.DbContent,
				courseLangsCollection,
				bson.M{
					"lan_code": lang.Lan_code,
				},
				bson.M{
					"$set": lang,
				},
				true,
			)
		}
	}
}

//更新孩子目录以及相关内容
func updateChildeCatalogs(
	catalogUUID string,
	toConfig mgdb.EnvConfig,
) {
	qmlog.QMLog.Info("上线catalogUUID:" + catalogUUID)

	var wg sync.WaitGroup
	//更新孩子目录以及相关内容
	var childCatalogs []*editor.Catalogs

	mgdb.Find(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		catalogCollection,
		bson.M{
			"parent_uuid": catalogUUID,
			"has_changed": true,
		},
		nil,
		nil,
		0,
		0,
		&childCatalogs,
	)

	for _, ct := range childCatalogs {
		qmlog.QMLog.Info("子目录UUID：" + ct.Uuid)

		mgdb.UpdateOne(
			toConfig,
			mgdb.DbContent,
			catalogCollection,
			bson.M{
				"uuid": ct.Uuid,
			},
			bson.M{
				"$set": ct,
			},
			true,
		)
		if toConfig == mgdb.EnvOnline {
			mgdb.UpdateOne(
				mgdb.EnvEditor,
				mgdb.DbEditor,
				catalogCollection,
				bson.M{
					"uuid": ct.Uuid,
				},
				bson.M{
					"$set": bson.M{
						"has_changed": false,
					},
				},
				false,
			)
		}
		if ct.Type == "catalog" {
			go func(tmpUUID string) {
				wg.Add(1)
				defer wg.Done()
				goCatalogs(
					tmpUUID,
					toConfig,
				)
			}(ct.Uuid)
		} else {
			var contents []bson.M
			mgdb.Find(
				mgdb.EnvEditor,
				mgdb.DbEditor,
				ct.Content_model,
				bson.M{
					"parent_uuid": ct.Uuid,
					"has_changed": true,
				},
				nil,
				nil,
				0,
				0,
				&contents,
			)
			for _, cnt := range contents {
				qmlog.QMLog.Info("内容UUID：" + cnt["uuid"].(string))

				mgdb.UpdateOne(
					toConfig,
					mgdb.DbContent,
					ct.Content_model,
					bson.M{
						"uuid": cnt["uuid"],
					},
					bson.M{
						"$set": cnt,
					},
					true,
				)
				if toConfig == mgdb.EnvOnline {
					mgdb.UpdateOne(
						mgdb.EnvEditor,
						mgdb.DbEditor,
						ct.Content_model,
						bson.M{
							"uuid": cnt["uuid"],
						},
						bson.M{
							"$set": bson.M{
								"has_changed": false,
							},
						},
						true,
					)
				}
			}
		}
	}

	wg.Wait()

	qmlog.QMLog.Info("结束上线子目录以及内容")
}

func goCatalogs(
	catalogUUID string,
	toConfig mgdb.EnvConfig,
) {
	qmlog.QMLog.Info("上线catalogUUID:" + catalogUUID)

	//更新孩子目录以及相关内容
	var childCatalogs []*editor.Catalogs

	mgdb.Find(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		catalogCollection,
		bson.M{
			"parent_uuid": catalogUUID,
			"has_changed": true,
		},
		nil,
		nil,
		0,
		0,
		&childCatalogs,
	)

	for _, ct := range childCatalogs {
		qmlog.QMLog.Info("子目录UUID：" + ct.Uuid)

		mgdb.UpdateOne(
			toConfig,
			mgdb.DbContent,
			catalogCollection,
			bson.M{
				"uuid": ct.Uuid,
			},
			bson.M{
				"$set": ct,
			},
			true,
		)
		if toConfig == mgdb.EnvOnline {
			mgdb.UpdateOne(
				mgdb.EnvEditor,
				mgdb.DbEditor,
				catalogCollection,
				bson.M{
					"uuid": ct.Uuid,
				},
				bson.M{
					"$set": bson.M{
						"has_changed": false,
					},
				},
				true,
			)
		}
		if ct.Type == "catalog" {
			goCatalogs(
				ct.Uuid,
				toConfig,
			)
		} else {
			var contents []bson.M
			mgdb.Find(
				mgdb.EnvEditor,
				mgdb.DbEditor,
				ct.Content_model,
				bson.M{
					"parent_uuid": ct.Uuid,
					"has_changed": true,
				},
				nil,
				nil,
				0,
				0,
				&contents,
			)
			for _, cnt := range contents {
				qmlog.QMLog.Info("内容UUID：" + cnt["uuid"].(string))

				mgdb.UpdateOne(
					toConfig,
					mgdb.DbContent,
					ct.Content_model,
					bson.M{
						"uuid": cnt["uuid"],
					},
					bson.M{
						"$set": cnt,
					},
					true,
				)
				if toConfig == mgdb.EnvOnline {
					mgdb.UpdateOne(
						mgdb.EnvEditor,
						mgdb.DbEditor,
						ct.Content_model,
						bson.M{
							"uuid": cnt["uuid"],
						},
						bson.M{
							"$set": bson.M{
								"has_changed": false,
							},
						},
						true,
					)
				}
			}
		}
	}
}

//上线课程信息，包括语种和课程，课程属性标签
func PushOnlineCourseInfos(msg *PushOnlineCourseMsg) {
	//设置recover，recover只能放在defer后面使用
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Printf("==> %s\n", string(buf[:n]))
		}
	}()

	//更新任务状态

	mgdb.UpdateOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		onLineJobsCollection,
		bson.M{
			"uuid": msg.UUID,
		},
		bson.M{
			"$set": bson.M{
				"state": 1,
			},
		},
		true,
	)
	//上线课程信息
	where := bson.M{
		"code": bson.M{
			"$in": msg.CourseCodes,
		},
	}
	if len(msg.CourseCodes) == 0 {
		where = nil
	}
	var courses []bson.M

	mgdb.Find(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		courseInfoCollection,
		where,
		nil,
		bson.M{"_id": 0},
		0,
		0,
		&courses,
	)

	for _, course := range courses {

		mgdb.UpdateOne(
			mgdb.EnvOnline,
			mgdb.DbContent,
			courseInfoCollection,
			bson.M{
				"code": course["code"],
			},
			bson.M{
				"$set": course,
			},
			true,
		)
		mgdb.UpdateOne(
			mgdb.EnvTest,
			mgdb.DbContent,
			courseInfoCollection,
			bson.M{
				"code": course["code"],
			},
			bson.M{
				"$set": course,
			},
			true,
		)
		//上线语言种类
		var lang *editor.Course_langs

		mgdb.FindOne(
			mgdb.EnvEditor,
			mgdb.DbEditor,
			courseLangsCollection,
			bson.M{"lan_code": course["lan_code"]},
			nil,
			&lang,
		)

		if lang != nil {
			qmlog.QMLog.Info("上线语言种类信息：" + lang.Lan_code)
			mgdb.UpdateOne(
				mgdb.EnvOnline,
				mgdb.DbContent,
				courseLangsCollection,
				bson.M{
					"lan_code": lang.Lan_code,
				},
				bson.M{
					"$set": lang,
				},
				true,
			)
			mgdb.UpdateOne(
				mgdb.EnvTest,
				mgdb.DbContent,
				courseLangsCollection,
				bson.M{
					"lan_code": lang.Lan_code,
				},
				bson.M{
					"$set": lang,
				},
				true,
			)
		}
	}

	mgdb.UpdateOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		"online_jobs",
		bson.M{
			"uuid": msg.UUID,
		},
		bson.M{
			"$set": bson.M{
				"state": 2,
			},
		},
		false,
	)
	//上线课程标签
	qmlog.QMLog.Info("开始上线课程标签")
	var tags []*editor.Content_tags

	mgdb.Find(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		contentTagsCollection,
		bson.M{
			"has_changed": true,
		},
		nil, nil, 0, 0, &tags,
	)

	for _, t := range tags {
		var ct *editor.Catalogs
		mgdb.FindOne(
			mgdb.EnvEditor,
			mgdb.DbEditor,
			tblCatalogs,
			bson.M{
				"attr_tag": t.Key,
			},
			nil, &ct,
		)

		if ct != nil {
			uuid := getVersionUUID(ct.Parent_uuid)

			mgdb.UpdateOne(
				mgdb.EnvOnline,
				mgdb.DbContent,
				tblContentInfo,
				bson.M{
					"uuid": uuid,
				},
				bson.M{
					"$set": bson.M{
						"update_time": time.Now().Unix(),
					},
				},
				false,
			)

			mgdb.UpdateOne(
				mgdb.EnvTest,
				mgdb.DbContent,
				tblContentInfo,
				bson.M{
					"uuid": uuid,
				},
				bson.M{
					"$set": bson.M{
						"update_time": time.Now().Unix(),
					},
				},
				false,
			)
		}
		mgdb.UpdateOne(
			mgdb.EnvOnline,
			mgdb.DbContent,
			contentTagsCollection,
			bson.M{
				"key": t.Key,
			},
			bson.M{
				"$set": t,
			},
			true,
		)
		mgdb.UpdateOne(
			mgdb.EnvTest,
			mgdb.DbContent,
			contentTagsCollection,
			bson.M{
				"key": t.Key,
			},
			bson.M{
				"$set": t,
			},
			true,
		)
	}
	qmlog.QMLog.Info("结束上线课程标签")
	qmlog.QMLog.Info("结束上线")
	fmt.Println(msg)
}

//获取内容版本的UUID
func getVersionUUID(parentUUID string) string {
	var ct *editor.Catalogs
	mgdb.FindOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tblCatalogs,
		bson.M{
			"uuid": parentUUID,
		},
		nil,
		&ct,
	)
	if ct == nil {
		return parentUUID
	} else {
		return getVersionUUID(ct.Parent_uuid)
	}

	return ""
}
