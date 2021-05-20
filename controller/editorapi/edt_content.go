package editorapi

import (
	"context"
	"editorApi/controller/servers"
	"editorApi/init/initNats"
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"editorApi/tools/helpers"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	uuid "github.com/satori/go.uuid"
	"github.com/tealeg/xlsx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"gopkg.in/go-playground/validator.v9"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
)

type contentGetParam struct {
	ParentUuid   string `json:"parent_uuid"`
	ContentModel string `json:"content_model"`
}

// @Tags EditorContentAPI（内容接口）
// @Summary 获取文件内容
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.contentGetParam true "获取文件内容"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/content [post]
func Content(ctx *gin.Context) {

	var param contentGetParam
	ctx.BindJSON(&param)
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(param.ContentModel)
	cusor, e := collection.Find(
		ctx,
		bson.M{
			"parent_uuid": param.ParentUuid,
			"has_del":     false,
		},
		options.Find().SetSort(bson.M{
			"list_order": 1,
		}).SetProjection(bson.M{
			"_id": 0,
		}),
	)
	defer cusor.Close(ctx)

	if e != nil {
		checkErr(ctx, e)
		return
	}
	var contents []bson.M
	cusor.All(ctx, &contents)
	// if param.ContentModel == "content_model_pro_sound" {
	// 	for _, c := range contents {
	// 		if ms, ok := c["motherSound"].(string); !ok || ms == "" {

	// 			if code, ok := c["code"].(string); ok && code != "" {
	// 				motherCode := code

	// 				if strings.Contains(code, "CHI-Basic") {
	// 					motherCode = strings.Replace(code, "CHI", "ENG", 1)
	// 				} else {
	// 					codeSlice := strings.Split(code, "-")
	// 					codeSlice[0] = "CHI"
	// 					motherCode = strings.Join(codeSlice, "-")
	// 				}
	// 				var oldC *struct {
	// 					Sentence string `bson:"sentence"`
	// 					Sound    string `bson:"sound"`
	// 				}

	// 				mgdb.FindOne(
	// 					mgdb.EnvEditor,
	// 					mgdb.DbEditor,
	// 					param.ContentModel,
	// 					bson.M{
	// 						"code": motherCode,
	// 					},
	// 					nil,
	// 					&oldC,
	// 				)

	// 				if oldC != nil {
	// 					c["motherSound"] = oldC.Sound
	// 					set := bson.M{
	// 						"motherSound": oldC.Sound,
	// 					}

	// 					if s, ok := c["sentence_temp"].(string); !ok || s == "" {
	// 						c["sentence_temp"] = oldC.Sentence
	// 						set["sentence_temp"] = oldC.Sentence
	// 					}

	// 					mgdb.UpdateOne(
	// 						mgdb.EnvEditor,
	// 						mgdb.DbEditor,
	// 						param.ContentModel,
	// 						bson.M{
	// 							"uuid": c["uuid"],
	// 						},
	// 						bson.M{
	// 							"$set": set,
	// 						},
	// 						false,
	// 					)
	// 				}
	// 			}
	// 		}

	// 	}
	// }
	servers.ReportFormat(ctx, true, "内容列表", gin.H{
		"contents": contents,
	})
}

type csp struct {
	SearchType   uint8  `json:"searchType"` // 0模糊搜索，1精确搜索
	Words        string `json:"words"`
	ContentModel string `json:"content_model"`
}

// @Tags EditorContentAPI（内容接口）
// @Summary 内容搜索
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.csp true "内容搜索"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/content/search [post]
func ContentSearch(ctx *gin.Context) {

	var param csp
	ctx.BindJSON(&param)
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(param.ContentModel)

	q := bson.M{
		"sentence": bsonx.Regex(param.Words, "i"),
	}

	if param.SearchType == 1 {
		q = bson.M{
			"sentence": param.Words,
		}
	}
	cusor, e := collection.Find(
		ctx,
		q,
		options.Find().SetSort(bson.M{
			"list_order": 1,
		}).SetProjection(bson.M{
			"_id":      0,
			"sentence": 1,
			"image":    1,
			"sound":    1,
		}).SetLimit(200),
	)
	defer cusor.Close(ctx)

	if e != nil {
		checkErr(ctx, e)
		return
	}
	var contents []bson.M
	cusor.All(ctx, &contents)
	servers.ReportFormat(ctx, true, "内容列表", gin.H{
		"contents": contents,
	})
}

// @Tags EditorContentAPI（内容接口）
// @Summary 获取内容题型列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/content/showTypes [get]
func ContentTypes(ctx *gin.Context) {

	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection("content_types")
	cusor, e := collection.Find(
		ctx,
		bson.M{
			"has_del": false,
		},
		options.Find().SetProjection(bson.M{
			"_id":     0,
			"has_del": 0,
		}),
	)
	defer cusor.Close(ctx)

	if e != nil {
		checkErr(ctx, e)
		return
	}
	var types []bson.M
	cusor.All(ctx, &types)
	servers.ReportFormat(ctx, true, "内容类型列表", gin.H{
		"showTypes": types,
	})
}

type contentEditParam struct {
	ParentUuid   string   `json:"parent_uuid"`   //文件目录的UUID
	ContentModel string   `json:"content_model"` //内容模型名
	Contents     []bson.M `json:"contents"`      //多个元素的文件内容
}

// @Tags EditorContentAPI（内容接口）
// @Summary 编辑文件内容
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.contentEditParam true "编辑文件内容，Contents是内容数组，内容的字段，根据内容模型设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/content/edit [post]
func ContentEdit(ctx *gin.Context) {
	var param contentEditParam
	ctx.BindJSON(&param)
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(param.ContentModel)

	//记录操作日志
	var (
		operation              OperateLogServer
		operateLogCreateParams = OperateLogParams{Ctx: ctx, Model: param.ContentModel, ParentUuid: param.ParentUuid, Pattern: PATTERN_CONTENT, Mold: MOLD_CONTENT_EDIT}
		operateLogEditParams   = OperateLogParams{Ctx: ctx, Model: param.ContentModel, ParentUuid: param.ParentUuid, Pattern: PATTERN_CONTENT, Mold: MOLD_CONTENT_EDIT}
	)
	operation.Start(operateLogCreateParams)
	defer operation.Finish(operateLogEditParams)

	for _, c := range param.Contents {

		if uid, ok := c["uuid"]; !ok || uid == "" {
			c["uuid"] = uuid.NewV4().String()
		}
		c["parent_uuid"] = param.ParentUuid
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
	updateParentCatalogInfo(ctx, []string{
		param.ParentUuid,
	})
	servers.ReportFormat(ctx, true, "编辑成功", gin.H{
		"contents": param.Contents,
	})
}

type contentDelParam struct {
	ParentUuid   string   `json:"parent_uuid"`
	ContentModel string   `json:"content_model"`
	DelUuids     []string `json:"del_uuids"`
}

// @Tags EditorContentAPI（内容接口）
// @Summary 删除文件内容
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.contentDelParam true "删除文件内容"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/content/del [post]
func ContentDel(ctx *gin.Context) {

	var param contentDelParam
	ctx.BindJSON(&param)
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(param.ContentModel)

	//记录操作日志
	var (
		operation              OperateLogServer
		operateLogCreateParams = OperateLogParams{Ctx: ctx, Model: param.ContentModel, ParentUuid: param.ParentUuid, Pattern: PATTERN_CONTENT, Mold: MOLD_CONTENT_CREATE}
		operateLogDelParams    = OperateLogParams{Ctx: ctx, Model: param.ContentModel, ParentUuid: param.ParentUuid, Pattern: PATTERN_CONTENT, Mold: MOLD_CONTENT_DELETE}
	)
	operation.Start(operateLogCreateParams)
	defer operation.Finish(operateLogDelParams)

	r, e := collection.UpdateMany(
		ctx,
		bson.M{
			"uuid": bson.M{
				"$in": param.DelUuids,
			},
		},
		bson.M{
			"$set": bson.M{
				"has_del":     true,
				"has_changed": true,
				"update_time": time.Now().Unix(),
			},
		},
	)
	if e != nil {
		checkErr(ctx, e)
		return
	}

	updateParentCatalogInfo(ctx, []string{
		param.ParentUuid,
	})

	servers.ReportFormat(ctx, true, "删除成功", gin.H{
		"modifiedCount": r.ModifiedCount,
	})
}

// @Tags EditorContentAPI（内容接口）
// @Summary 导入文件内容
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param file query string true "文件名称"
// @Param parent_uuid query string true "parent_uuid"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/content/import [post]
func ContentImport(ctx *gin.Context) {
	//获取表单数据 参数为name值
	file, err := ctx.FormFile("file")
	parentUuid := ctx.PostForm("parent_uuid")
	//错误处理
	if err != nil {
		checkErr(ctx, err)
		return
	} else {
		// 打开文件
		src, err := file.Open()
		defer src.Close()

		// 读取文件内容
		xlFile, err := xlsx.OpenReaderAt(src, file.Size)
		if err != nil {
			checkErr(ctx, err)
			return
		}
		// 文档数据校验
		err = checkSheets(parentUuid, xlFile.Sheets, false)
		if err != nil {
			checkErr(ctx, err)
			return
		}

		// 文件保存
		fielPath := "static/excel/" + uuid.NewV4().String() + ".xlsx"
		err = ctx.SaveUploadedFile(file, fielPath)
		if err != nil {
			checkErr(ctx, err)
			return
		}

		// 异步操作
		initNats.NatsConn.Publish("ImportCourseContent",
			&ImportContentMsg{
				ParentUuid: parentUuid,
				FilePath:   fielPath,
			},
		)

		//保存成功返回正确的Json数据
		servers.ReportFormat(ctx, true, "上传成功", nil)
	}
}

type ImportContentMsg struct {
	ParentUuid string
	FilePath   string
}

// 处理文件
func HanderImportContent(msg *ImportContentMsg) {
	// 读取文件内容
	xlFile, _ := xlsx.OpenFile(msg.FilePath)

	// 文档数据校验
	checkSheets(msg.ParentUuid, xlFile.Sheets, true)

	// 删除临时文件
	os.Remove(msg.FilePath)
}

/**
数据校验
*/
func checkSheets(parent_uuid string, sheets []*xlsx.Sheet, flag bool) (err error) {
	var catalog, slide string
	var validateErr error
	var contents []ContentModelMroMound
	var levelGrade []ContentModelLevelGrade
	// 遍历sheet页读取
	for _, sheet := range sheets {
		sheetName := sheet.Name
		//数据校验
		for index, row := range sheet.Rows {
			if index > 0 {
				var parentUuid, pattern, chineseSentence, sentence, image, sound, uniqueCode string
				var options, optionsPhoneticize []string
				for key, cell := range row.Cells {
					cellValue := strings.TrimSpace(cell.String())
					switch key {
					case 0: //课程类别
						if !helpers.Empty(cellValue) {
							catalog = strings.TrimSuffix(strings.TrimSpace(cellValue), ">")
						}
						if helpers.Empty(catalog) {
							err = errors.New(sheetName + ":课程类别不能为空！")
						}
					case 1: //Slide
						if !helpers.Empty(cellValue) {
							slide = cellValue
						}
						if helpers.Empty(slide) {
							err = errors.New(sheetName + ":Slide不能为空！")
						}
						parentUuid = getLastCatalog(parent_uuid, catalog, slide)
					case 2: //题型
						if helpers.Empty(cellValue) {
							err = errors.New(sheetName + ":题型数据不能为空！")
						}
						pattern = getType(cellValue)
					case 3: //中文语句
						if !helpers.Empty(cellValue) {
							chineseSentence = cellValue
						}
					case 4: //语句
						if !helpers.Empty(cellValue) {
							sentence = cellValue
						}
					case 5: //图片地址
						if !helpers.Empty(cellValue) {
							image = cellValue
						}
					case 6: //声音地址(不做处理)
						if !helpers.Empty(cellValue) {
							sound = cellValue
						}
					case 7: //选项
						if !helpers.Empty(cellValue) {
							optionsTmp := strings.Split(cellValue, "|")
							for _, item := range optionsTmp {
								options = append(options, strings.TrimSpace(item))
							}
						}
					case 8: //选项拼音
						if !helpers.Empty(cellValue) {
							optionsPhoneticizeTmp := strings.Split(cellValue, "|")
							for _, item := range optionsPhoneticizeTmp {
								optionsPhoneticize = append(optionsPhoneticize, strings.TrimSpace(item))
							}
						}
					case 9: //标记
					case 10: //本次导入的唯一标识
						if helpers.Empty(cellValue) {
							err = errors.New(sheetName + ":本次导入的唯一标识不能为空！")
						}
						uniqueCode = cellValue
					}
				}

				switch pattern {
				case "listenForSentence", "speakToImg", "writeWords", "makeSentence", "fillGap", "countDown", "makeWord", "makePhrase":
					if helpers.Empty(optionsPhoneticize) {
						optionsPhoneticize = []string{"1"}
					}
				}

				// 生成uuid
				switch pattern {
				case "listenForSentence":
					var contentModelLevelGrade ContentModelLevelGrade
					contentModelLevelGrade.Uuid = uuid.NewV4().String()
					contentModelLevelGrade.HasDel = false
					contentModelLevelGrade.SentenceShow = false
					contentModelLevelGrade.ListOrder = 10 * index
					contentModelLevelGrade.ParentUuid = parentUuid
					contentModelLevelGrade.ContentType = pattern
					contentModelLevelGrade.Sentence = sentence
					contentModelLevelGrade.SentenceAudio = sound
					contentModelLevelGrade.WordsChoice = options
					contentModelLevelGrade.ImgsChoice = []string{}
					contentModelLevelGrade.Radar = []string{}
					contentModelLevelGrade.UniqueCode = uniqueCode
					//中文翻译器
					zh_ch := zh.New()
					uni := ut.New(zh_ch)
					trans, _ := uni.GetTranslator("zh")
					//验证器
					validate := validator.New()
					//验证器注册翻译器
					zh_translations.RegisterDefaultTranslations(validate, trans)
					validateErr = validate.Struct(contentModelLevelGrade)
					if validateErr != nil {
						//for _, err := range err.(validator.ValidationErrors) {
						//	fmt.Println(err.Translate(trans))
						//}
						err = errors.New(sheetName + ":文件校验失败！")
					}

					if flag {
						levelGrade = append(levelGrade, contentModelLevelGrade)
					}
				default:
					var ContentModel ContentModelMroMound
					ContentModel.Uuid = uuid.NewV4().String()
					ContentModel.HasChanged = false
					ContentModel.HasDel = false
					ContentModel.IsShow = true
					ContentModel.ListOrder = 10 * index
					ContentModel.UpdateTime = time.Now().Unix()
					ContentModel.ParentUuid = parentUuid
					ContentModel.Type = pattern
					ContentModel.ChineseSentence = chineseSentence
					ContentModel.Sentence = sentence
					ContentModel.Image = image
					ContentModel.Sound = sound
					ContentModel.Options = options
					ContentModel.OptionsPhoneticize = optionsPhoneticize
					ContentModel.UniqueCode = uniqueCode
					//中文翻译器
					zh_ch := zh.New()
					uni := ut.New(zh_ch)
					trans, _ := uni.GetTranslator("zh")
					//验证器
					validate := validator.New()
					//验证器注册翻译器
					zh_translations.RegisterDefaultTranslations(validate, trans)
					validateErr = validate.Struct(ContentModel)
					if validateErr != nil {
						//for _, err := range err.(validator.ValidationErrors) {
						//	fmt.Println(err.Translate(trans))
						//}
						err = errors.New(sheetName + ":文件校验失败！")
					}
					if flag {
						contents = append(contents, ContentModel)
					}
				}
			}
		}
	}

	if flag {
		// 添加内容
		if !helpers.Empty(contents) {
			addContent(contents)
		}

		if !helpers.Empty(levelGrade) {
			addContentLevelGrade(levelGrade)
		}
	}

	return
}

func addContent(contents []ContentModelMroMound) {
	var mutex sync.Mutex
	// 遍历sheet页读取
	for _, content := range contents {
		mutex.Lock()
		// 1，查询数据是否存在  content_model_pro_sound
		tmp := findContent(content.ParentUuid, content.UniqueCode)

		collection := mgdb.MongoClient.Database(EDITOR_DB).Collection("content_model_pro_sound")
		var ctx context.Context
		if helpers.Empty(tmp) {
			// 2，添加、更新数据
			insertResult, _ := collection.InsertOne(ctx, content)
			fmt.Println(insertResult)
		} else {
			collection.UpdateOne(
				ctx,
				bson.M{
					"uuid": tmp.Uuid,
				},
				bson.M{
					"$set": bson.M{
						"options":             content.Options,
						"options_phoneticize": content.OptionsPhoneticize,
						"image":               content.Image,
						"type":                content.Type,
						"has_changed":         true,
						"sentence":            content.Sentence,
						"sound":               content.Sound,
						"chinese_sentence":    content.ChineseSentence,
						"update_time":         time.Now().Unix(),
					},
				},
				options.Update().SetUpsert(true),
			)
		}
		mutex.Unlock()
	}
}

func addContentLevelGrade(contents []ContentModelLevelGrade) {
	var mutex sync.Mutex
	// 遍历sheet页读取
	for _, content := range contents {
		mutex.Lock()
		// 1，查询数据是否存在  content_model_pro_sound
		tmp := findContentLevelGrade(content.ParentUuid, content.UniqueCode)

		collection := mgdb.MongoClient.Database(EDITOR_DB).Collection("content_model_level_grade")
		var ctx context.Context
		if helpers.Empty(tmp) {
			// 2，添加、更新数据
			insertResult, _ := collection.InsertOne(ctx, content)
			fmt.Println(insertResult)
		} else {
			collection.UpdateOne(
				ctx,
				bson.M{
					"uuid": tmp.Uuid,
				},
				bson.M{
					"$set": bson.M{
						"list_order":     content.ListOrder,
						"content_type":   content.ContentType,
						"words_choice":   content.WordsChoice,
						"imgs_choice":    content.ImgsChoice,
						"radar":          content.Radar,
						"sentence":       content.Sentence,
						"sentence_audio": content.SentenceAudio,
					},
				},
				options.Update().SetUpsert(true),
			)
		}

		collectionCatalogs := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)
		collectionCatalogs.UpdateOne(
			ctx,
			bson.M{
				"uuid": content.ParentUuid,
			},
			bson.M{
				"$set": bson.M{
					"content_model": "content_model_level_grade",
				},
			},
			options.Update().SetUpsert(true),
		)

		mutex.Unlock()
	}
}

func findContent(parent_uuid, unique_code string) (content *ContentModelMroMound) {
	mgdb.FindOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		"content_model_pro_sound",
		bson.M{
			"parent_uuid": parent_uuid,
			"unique_code": unique_code,
			"has_del":     false,
		},
		nil,
		&content,
	)
	return
}

func findContentLevelGrade(parent_uuid, unique_code string) (content *ContentModelLevelGrade) {
	mgdb.FindOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		"content_model_level_grade",
		bson.M{
			"parent_uuid": parent_uuid,
			"unique_code": unique_code,
			"has_del":     false,
		},
		nil,
		&content,
	)
	return
}

func getLastCatalog(parent_uuid, catalog, slide string) string {
	var mutex sync.Mutex
	var ctx context.Context
	catalogData := strings.Split(catalog, ">")
	for key, item := range catalogData {
		mutex.Lock()
		name := strings.TrimSpace(item)
		catalogData, _ := findCatalog(parent_uuid, name, "catalog")
		catalogCount, _ := countCatalog(parent_uuid)
		if helpers.Empty(catalogData) {
			catalogUuid := uuid.NewV4().String()
			desc := map[string]string{
				"zh-CN": name,
				"en":    "",
			}
			title := map[string]string{
				"zh-CN": name,
				"en":    "",
			}
			var catalog catalogsInfo
			catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)
			catalog.Type = "catalog"
			catalog.ListOrder = 10 * (key + 1 + int(catalogCount))
			catalog.Name = name
			catalog.ParentUUID = parent_uuid
			catalog.UUID = catalogUuid
			catalog.HasChanged = true
			catalog.IsShow = true
			catalog.HasDel = false
			catalog.Flag = []string{}
			catalog.Cover = []string{}
			catalog.Desc = desc
			catalog.Title = title
			catalog.UpdateTime = time.Now().Unix()
			catalogsCollection.InsertOne(ctx, catalog)
			parent_uuid = catalogUuid
		} else {
			parent_uuid = catalogData.Uuid
		}

		mutex.Unlock()
	}

	slideData, _ := findCatalog(parent_uuid, slide, "content")
	slideCount, _ := countCatalog(parent_uuid)
	if helpers.Empty(slideData) {
		slideUuid := uuid.NewV4().String()
		var catalog catalogsInfo
		desc := map[string]string{
			"zh-CN": slide,
			"en":    "",
		}
		title := map[string]string{
			"zh-CN": slide,
			"en":    "",
		}
		catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)
		catalog.Type = "content"
		catalog.ListOrder = 10 * (int(slideCount) + 1)
		catalog.Name = slide
		catalog.ParentUUID = parent_uuid
		catalog.UUID = slideUuid
		catalog.HasChanged = true
		catalog.IsShow = true
		catalog.HasDel = false
		catalog.Flag = []string{}
		catalog.Cover = []string{}
		catalog.Desc = desc
		catalog.Title = title
		catalog.ContentModel = "content_model_pro_sound"
		catalog.UpdateTime = time.Now().Unix()
		catalogsCollection.InsertOne(ctx, catalog)
		return slideUuid
	} else {
		return slideData.Uuid
	}
}

func countCatalog(parent_uuid string) (count int64, err error) {
	count, err = mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs).CountDocuments(context.Background(), bson.M{
		"parent_uuid": parent_uuid,
		"has_del":     false,
	})
	return
}

func findCatalog(parent_uuid, name, mold string) (catalog *editor.Catalogs, err error) {
	mgdb.FindOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tblCatalogs,
		bson.M{
			"parent_uuid": parent_uuid,
			"name":        name,
			"type":        mold,
			"has_del":     false,
		},
		nil,
		&catalog,
	)
	if helpers.Empty(catalog) {
		err = errors.New("数据不存在")
	}
	return
}

func getType(name string) string {
	var value string
	switch strings.TrimSpace(name) {
	case "听声音选句子（单词）":
		value = "listenForSentence"
	case "自动读":
		value = "autoSpeak"
	case "跟读":
		value = "repeatSpeak"
	case "句子选图片":
		value = "sentenceToImg"
	case "由图片选句子":
		value = "imgToSentence"
	case "根据图片跟读":
		value = "speakToImg"
	case "写单词":
		value = "writeWords"
	case "组句子":
		value = "makeSentence"
	case "选词填空":
		value = "fillGap"
	case "倒计时":
		value = "countDown"
	case "组单词":
		value = "makeWord"
	case "组短语":
		value = "makePhrase"
	}
	return value
}


func getTypeContent(name interface{}) string {
	var value string
	switch name {
	case "listenForSentence":
		value = "听声音选句子（单词）"
	case "autoSpeak":
		value = "自动读"
	case "repeatSpeak":
		value = "跟读"
	case "sentenceToImg":
		value = "句子选图片"
	case "imgToSentence":
		value = "由图片选句子"
	case "speakToImg":
		value = "根据图片跟读"
	case "writeWords":
		value = "写单词"
	case "makeSentence":
		value = "组句子"
	case "fillGap":
		value = "选词填空"
	case "countDown":
		value = "倒计时"
	case "makeWord":
		value = "组单词"
	case "makePhrase":
		value = "组短语"
	}
	return value
}

type ContentModelMroMound struct {
	Uuid                string   `bson:"uuid" json:"uuid"`
	ParentUuid          string   `bson:"parent_uuid" json:"parent_uuid" validate:"required"`
	ListOrder           int      `bson:"list_order" json:"list_order"`
	Code                string   `bson:"code" json:"code"`
	SentenceTemp        string   `bson:"sentence_temp" json:"sentence_temp"`
	Options             []string `bson:"options" json:"options"`
	Sound               string   `bson:"sound" json:"sound"`
	Type                string   `bson:"type" json:"type"`
	OptionsPhoneticize  []string `bson:"options_phoneticize" json:"options_phoneticize"`
	SentencePhoneticize string   `bson:"sentence_phoneticize" json:"sentence_phoneticize"`
	Image               string   `bson:"image" json:"image"`
	Sentence            string   `bson:"sentence" json:"sentence"`
	ChineseSentence     string   `bson:"chinese_sentence" json:"chinese_sentence"`
	UpdateTime          int64    `bson:"update_time" json:"update_time"`
	HasDel              bool     `bson:"has_del" json:"has_del"`
	HasChanged          bool     `bson:"has_changed" json:"has_changed"`
	IsShow              bool     `bson:"is_show" json:"is_show"`
	UniqueCode          string   `bson:"unique_code" json:"unique_code"`
}

type ContentModelLevelGrade struct {
	Uuid              string   `bson:"uuid" json:"uuid"`
	ParentUuid        string   `bson:"parent_uuid" json:"parent_uuid" validate:"required"`
	ListOrder         int      `bson:"list_order" json:"list_order"`
	ContentType       string   `bson:"content_type" json:"content_type"`
	SentenceShow      bool     `bson:"sentence_show" json:"sentence_show"`
	Sentence          string   `bson:"sentence" json:"sentence"`
	SentenceAudio     string   `bson:"sentence_audio" json:"sentence_audio"`
	SentenceAudioTime int64    `bson:"sentence_audio_time" json:"sentence_audio_time"`
	Text              string   `bson:"text" json:"text"`
	TextAudio         string   `bson:"text_audio" json:"text_audio"`
	TextAudioTime     int      `bson:"text_audio_time" json:"text_audio_time"`
	Radar             []string `bson:"radar" json:"radar"`
	UseTime           int      `bson:"use_time" json:"use_time"`
	WordsChoice       []string `bson:"words_choice" json:"words_choice"`
	ImgsChoice        []string `bson:"imgs_choice" json:"imgs_choice"`
	Score             int      `bson:"score" json:"score"`
	HasDel            bool     `bson:"has_del" json:"has_del"`
	UniqueCode        string   `bson:"unique_code" json:"unique_code"`
}
