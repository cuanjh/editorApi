package editorapi

import (
	"bytes"
	"context"
	"editorApi/commons"
	"editorApi/init/initNats"
	"editorApi/init/mgdb"
	"editorApi/requests"
	"editorApi/service"
	"editorApi/tools/helpers"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/tealeg/xlsx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

/**
导入课程数据
*/
func ExportData(ctx *gin.Context) {
	dir, _ := os.Getwd()
	commons.Info(ctx, "处理中.........")
	excelFileName, err := ctx.FormFile("filename")
	sheetName := ctx.PostForm("sheet")
	code := ctx.PostForm("code")
	if err != nil || helpers.Empty(sheetName) {
		commons.Error(ctx, 500, err, "文件名不能为空！")
	}

	fmt.Println(sheetName)
	os.MkdirAll("data/dict/", os.ModePerm)
	dst := "data/dict/" + uuid.NewV4().String() + ".xlsx"
	// gin 简单做了封装,拷贝了文件流
	if err := ctx.SaveUploadedFile(excelFileName, dst); err != nil {
		commons.Error(ctx, 500, err, "文件保存失败！")
	}

	xlFile, err := xlsx.OpenFile(dst)
	if err != nil {
		panic(err)
	}
	var lock = sync.Mutex{}
	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			if i > 0 {
				lock.Lock()
				var tmp courseDataParam
				for key, cell := range row.Cells {
					if key == 3 {
						tmp.Uuid = cell.String()
					}
					if key == 4 {
						tmp.Model = cell.String()
					}
				}

				var linkName = sheetName
				var mp3File = linkName + "_" + RandomCompletion(i) + ".mp3"

				soundsAaronFile, _ := Download("https://course-assets1.talkmate.com/actors/sounds/" + code + "/Aaron/" + linkName + "/" + mp3File)
				if err != nil {
					fmt.Println(err)
					return
				}
				AaronSoundTime, err := GenerateLength("ffmpeg", dir+"/"+soundsAaronFile)
				if err != nil {
					fmt.Println(err)
					return
				}
				os.Remove(soundsAaronFile)
				var Aaron = ActorsSound{
					Gender:    1,
					Role:      "b9a24556-98ca-4880-a650-216804af11ea",
					Sound:     "actors/sounds/" + code + "/Aaron/" + linkName + "/" + mp3File,
					SoundTime: strconv.FormatFloat(AaronSoundTime, 'g', -1, 64),
				}
				tmp.ActorsSounds = append(tmp.ActorsSounds, Aaron)

				soundsMilaFile, _ := Download("https://course-assets1.talkmate.com/actors/sounds/" + code + "/Mila/" + linkName + "/" + mp3File)
				MilaSoundTime, err := GenerateLength("ffmpeg", dir+"/"+soundsMilaFile)
				os.Remove(soundsMilaFile)
				var Mila = ActorsSound{
					Gender:    0,
					Role:      "4f154cda-be4b-48a8-a7c9-1533c66922e1",
					Sound:     "actors/sounds/" + code + "/Mila/" + linkName + "/" + mp3File,
					SoundTime: strconv.FormatFloat(MilaSoundTime, 'g', -1, 64),
				}
				tmp.ActorsSounds = append(tmp.ActorsSounds, Mila)

				soundsTomFile, _ := Download("https://course-assets1.talkmate.com/actors/sounds/" + code + "/Tom/" + linkName + "/" + mp3File)
				TomSoundTime, err := GenerateLength("ffmpeg", dir+"/"+soundsTomFile)
				os.Remove(soundsTomFile)
				var Tom = ActorsSound{
					Gender:    1,
					Role:      "cf08465e-711e-4b45-bdfc-bc130ef88644",
					Sound:     "actors/sounds/" + code + "/Tom/" + linkName + "/" + mp3File,
					SoundTime: strconv.FormatFloat(TomSoundTime, 'g', -1, 64),
				}
				tmp.ActorsSounds = append(tmp.ActorsSounds, Tom)

				soundsJohnFile, _ := Download("https://course-assets1.talkmate.com/actors/sounds/" + code + "/John/" + linkName + "/" + mp3File)
				JohnSoundTime, err := GenerateLength("ffmpeg", dir+"/"+soundsJohnFile)
				os.Remove(soundsJohnFile)
				var John = ActorsSound{
					Gender:    1,
					Role:      "11d70675-323b-4da8-befa-485e660cf2fa",
					Sound:     "actors/sounds/" + code + "/John/" + linkName + "/" + mp3File,
					SoundTime: strconv.FormatFloat(JohnSoundTime, 'g', -1, 64),
				}
				tmp.ActorsSounds = append(tmp.ActorsSounds, John)

				soundsKaylaFile, _ := Download("https://course-assets1.talkmate.com/actors/sounds/" + code + "/Kayla/" + linkName + "/" + mp3File)
				KaylaSoundTime, err := GenerateLength("ffmpeg", dir+"/"+soundsKaylaFile)
				os.Remove(soundsKaylaFile)
				var Kayla = ActorsSound{
					Gender:    0,
					Role:      "ddc825da-7b2e-4a87-97da-bf9a48483380",
					Sound:     "actors/sounds/" + code + "/Kayla/" + linkName + "/" + mp3File,
					SoundTime: strconv.FormatFloat(KaylaSoundTime, 'g', -1, 64),
				}
				tmp.ActorsSounds = append(tmp.ActorsSounds, Kayla)

				collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tmp.Model)
				_, err = collection.UpdateOne(ctx, bson.M{
					"uuid": tmp.Uuid,
				}, bson.M{
					"$set": bson.M{
						"actors_sound": tmp.ActorsSounds,
					},
				})
				if err != nil {
					checkErr(ctx, err)
					return
				}

				var result DataResult
				var filter = bson.D{}
				filter = append(filter, bson.E{"uuid", tmp.Uuid})
				dataResult := collection.FindOne(ctx, filter)
				dataResult.Decode(&result)
				var parentUuids []string
				parentUuids = append(parentUuids, result.ParentUuid)
				if !helpers.Empty(parentUuids) {
					updateParentCatalogInfo(ctx, parentUuids)
				}
				lock.Unlock()
			}
		}
	}

	commons.Info(ctx, "处理完成.........")
	commons.Success(ctx, nil, "提交成功！", nil)
}

type courseDataParam struct {
	Model        string        `json:"model"`
	Uuid         string        `json:"uuid" bson:"uuid"`
	ActorsSounds []ActorsSound `json:"actors_sound" bson:"actors_sound"`
}

type ActorsSound struct {
	Gender    int    `json:"gender" bson:"gender"`
	Role      string `json:"role" bson:"role"`
	Sound     string `json:"sound" bson:"sound"`
	SoundTime string `json:"sound_time" bson:"sound_time"`
}

// 通过ffmpeg 获取时长
func GenerateLength(ffmpegPath string, url string) (float64, error) {
	var length float64
	//视频处理使用，延长超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	cmd := exec.CommandContext(ctx, ffmpegPath, "-i", url)
	defer cancel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()
	str := stderr.String()
	if helpers.Empty(str) {
		return 0, nil
	}
	arr1 := strings.Split(str, "Duration:")
	if helpers.Empty(arr1) {
		return 0, nil
	}
	arr2 := strings.Split(arr1[1], ", start:")
	if helpers.Empty(arr2) {
		return 0, nil
	}
	str = "2006-01-02" + strings.TrimPrefix(arr2[0], "Duration:")
	start, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 00:00:00")
	videotime, _ := time.Parse("2006-01-02 15:04:05", str)
	length, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", float64(videotime.UnixNano()-start.UnixNano())/1e9), 64)
	return length, nil
}

func Download(url string) (path string, err error) {
	res, err := http.Get(url)
	path = "data/dict/" + uuid.NewV4().String() + ".mp3"
	f, err := os.Create(path)
	if err != nil {
		return
	}
	io.Copy(f, res.Body)
	return
}

type DataResult struct {
	ParentUuid string `json:"parent_uuid" bson:"parent_uuid"`
}

/*********************************************************************************************************************************************/

type CourseDataModle struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	LanCode      string             `json:"lan_code" bson:"lan_code"`
	Level        string             `json:"level" bson:"level"`
	Unit         string             `json:"unit" bson:"unit"`
	FormPicture  Picture            `json:"form_picture" bson:"form_picture"`
	FormSound    string             `json:"form_sound" bson:"form_sound"`
	Record       string             `json:"record" bson:"record"`
	Module       string             `json:"module" bson:"module"`
	SlideType    string             `json:"slide_type" bson:"slide_type"`
	Sentence     string             `json:"sentence" bson:"sentence"`
	FormShowType string             `json:"form_show_type" bson:"form_show_type"`
}

type Picture struct {
	Url string `json:"url" bson:"url"`
}

func CourseData(ctx *gin.Context) {
	for skip := 0; skip <= 40477; skip++ {
		lock.Lock()
		var result CourseDataModle
		model := mgdb.OnlineClient.Database(KUYU).Collection("course_eng").FindOne(ctx, bson.M{
			"is_handled": false,
		})
		err := model.Decode(&result)
		if err != nil {
			fmt.Println(err)
		}

		if !helpers.Empty(result.Sentence) {
			initNats.NatsConn.Publish("HandleCourseData", &result)
			_, err = mgdb.OnlineClient.Database(KUYU).Collection("course_eng").UpdateOne(ctx, bson.M{
				"_id": result.ID,
			}, bson.M{
				"$set": bson.M{
					"is_handled": true,
				},
			})
		}
		lock.Unlock()
	}

	//var courses []CourseDataModle
	//for skip := 0; skip <= 410; skip++ {
	//	model, err := toClient.Database(KUYU).Collection("course_eng").Find(
	//		ctx,
	//		bson.M{
	//			"unchange": true,
	//		},
	//		options.Find().SetSort(bson.M{"chapter": -1}),
	//		options.Find().SetLimit(1),
	//		options.Find().SetSkip(int64(0)),
	//	)
	//	defer model.Close(ctx)
	//	err = model.All(ctx, &courses)
	//	if err != nil {
	//		fmt.Println(err)
	//		continue
	//	}
	//	for _, course := range courses {
	//		var lock = sync.Mutex{}
	//		lock.Lock()
	//		if !helpers.Empty(course.Sentence) {
	//			// 异步操作
	//			initNats.NatsConn.Publish("HandleCourseData", &course)
	//		}
	//		lock.Unlock()
	//	}
	//}

	commons.Success(ctx, nil, "提交成功！", nil)
}

func HandleCourseData(course *CourseDataModle) {
	var ctx context.Context
	langs := "doz,hak,lit,taj,tam,guj,mar,swa,pas,pus,som,tel,tkm,tuk,uzb,ben,fil,lav,mal,nep,pun,snd,sun,tgk,esp,cze,lao,kin,lug,mac,slo,amb,aze,bos,bul,bur,bye,est,fan,fin,fra,glg,gre,heb,hye,jav,kan,kas,khm,kyr,lat,mon,mrd,run,sag,tet,yor,mog,far,tib,hin,urd,uyg,sin,alb,may,ukr,por,ice,ind,vie,hrv,ara,kaz,jpn,kor,chi,can,nor,tha,hun,ksp,spa,ser,dut,fre,kfr,pol,dan,ger,ita,ken,rum,rus,swe,tur"
	for _, FROM := range strings.Split(langs, ",") {
		//commons.Info(ctx, FROM)
		var lock = sync.Mutex{}
		lock.Lock()
		var result CourseDataModle
		model := mgdb.OnlineClient.Database(KUYU).Collection("course_"+strings.ToLower(FROM)).FindOne(ctx, bson.M{
			//"data_area": "cn",
			"level":          course.Level,
			"unit":           course.Unit,
			"module":         course.Module,
			"slide_type":     course.SlideType,
			"form_show_type": course.FormShowType,
		})
		err := model.Decode(&result)

		if err != nil {
			//commons.Info(ctx, "Error:"+err.Error())
			continue
		}

		if helpers.Empty(result.Sentence){
			model := mgdb.OnlineClient.Database(KUYU).Collection("course_"+strings.ToLower(FROM)).FindOne(ctx, bson.M{
				//"data_area": "cn",
				//"level":          course.Level,
				//"unit":           course.Unit,
				//"module":         course.Module,
				//"slide_type":     course.SlideType,
				//"form_show_type": course.FormShowType,
				"eng_content": course.Sentence,
			})
			err := model.Decode(&result)

			if err != nil {
				//commons.Info(ctx, "Error:"+err.Error())
				continue
			}
		}

		if !helpers.Empty(result.Sentence) {
			var sex = "male"
			if course.Record == "f" {
				sex = "female"
			}

			if IsWords(course.Sentence) && false {
				var request requests.DictDetailRequests
				UUID := helpers.MD5(strings.TrimRight(course.Sentence, "."))
				cardId := helpers.MD5(strings.TrimRight(course.Sentence, "."))
				request.Uuid = UUID
				request.From = "eng"
				obj := service.AppService()
				dict, err := obj.AppService.DictService.FindOne(ctx, request)
				if err != nil {
					//commons.Info(ctx, "Error:"+err.Error())
					continue
				}

				if helpers.Empty(dict) {
					var engSoundInfos []requests.SoundInfos
					if !helpers.Empty(course.FormSound) {
						engSoundInfos = append(engSoundInfos, requests.SoundInfos{
							Sound:  course.FormSound,
							Gender: sex,
						})
					}

					var engImages []requests.Image
					if !helpers.Empty(course.FormPicture) {
						engImages = append(engImages, requests.Image{
							Url: course.FormPicture.Url,
						})
					}

					var engParams requests.Dict
					engParams.From = "eng"
					engParams.CardId = cardId
					engParams.Content = result.Sentence
					engParams.IsDel = false
					engParams.SoundInfos = engSoundInfos
					engParams.Images = engImages
					_, err = obj.AppService.DictService.AddDict(ctx, engParams)
					if err != nil {
						//commons.Info(ctx, "Error:"+err.Error())
						continue
					}
				}

				var soundInfos []requests.SoundInfos
				if !helpers.Empty(result.FormSound) {
					soundInfos = append(soundInfos, requests.SoundInfos{
						Sound:  result.FormSound,
						Gender: sex,
					})
				}

				var images []requests.Image
				if !helpers.Empty(result.FormPicture) {
					images = append(images, requests.Image{
						Url: result.FormPicture.Url,
					})
				}

				var params requests.Dict
				params.From = FROM
				params.CardId = cardId
				params.Content = result.Sentence
				params.IsDel = false
				params.SoundInfos = soundInfos
				params.Images = images
				params.Uuid = UUID
				_, err = obj.AppService.DictService.AddDict(ctx, params)
				if err != nil {
					//commons.Info(ctx, "Error:"+err.Error())
					continue
				}
				// 添加翻译
			} else {
				var request requests.Sentence
				UUID := helpers.MD5(course.Sentence)
				cardId := helpers.MD5(course.Sentence)
				request.Uuid = UUID
				request.From = "eng"
				obj := service.AppService()
				sentence, err := obj.AppService.SentenceService.FindOne(ctx, request)
				if err != nil {
					//commons.Info(ctx, "Error:"+err.Error())
					continue
				}
				// 处理句子数据
				if helpers.Empty(sentence) {
					var engSoundInfos []requests.SentenceSoundInfos
					if !helpers.Empty(course.FormSound) {
						engSoundInfos = append(engSoundInfos, requests.SentenceSoundInfos{
							Sound:  course.FormSound,
							Gender: sex,
						})
					}

					var engImages []string
					if !helpers.Empty(course.FormPicture) {
						engImages = append(engImages, course.FormPicture.Url)
					}

					var engParams requests.Sentence
					engParams.From = "eng"
					engParams.Uuid = UUID
					engParams.CardId = cardId
					engParams.Sentence = course.Sentence
					engParams.IsDel = false
					engParams.SoundInfos = engSoundInfos
					engParams.Image = engImages
					_, err = obj.AppService.SentenceService.AddSentence(ctx, engParams)
					if err != nil {
						//commons.Info(ctx, "Error:"+err.Error())
						continue
					}
				} else {
					// 更新eng句子中的CardId
					var paramsSentenceCardId requests.SentenceCardId
					paramsSentenceCardId.From = "eng"
					paramsSentenceCardId.Uuid = UUID
					paramsSentenceCardId.CardId = cardId
					_, err = obj.AppService.SentenceService.SentenceAddCardId(ctx, paramsSentenceCardId)
					if err != nil {
						//commons.Info(ctx, "Error:"+err.Error())
						continue
					}
				}

				//翻译数据
				var soundInfos []requests.SentenceSoundInfos
				if !helpers.Empty(result.FormSound) {
					soundInfos = append(soundInfos, requests.SentenceSoundInfos{
						Sound:  result.FormSound,
						Gender: sex,
					})
				}

				var images []string
				if !helpers.Empty(result.FormPicture) {
					images = append(images, result.FormPicture.Url)
				}

				var sentenceUuid = helpers.MD5(result.Sentence)
				var params requests.Sentence
				params.From = FROM
				params.Uuid = sentenceUuid
				params.CardId = cardId
				params.Sentence = result.Sentence
				params.IsDel = false
				params.SoundInfos = soundInfos
				params.Image = images
				_, err = obj.AppService.SentenceService.AddSentence(ctx, params)
				if err != nil {
					//commons.Info(ctx, "Error:"+err.Error())
					continue
				}

				if FROM == "chi" {
					continue
				}
				//翻译
				var paramsSentenceTranslate requests.SentenceTranslate
				paramsSentenceTranslate.From = "eng"
				paramsSentenceTranslate.To = "chi"
				paramsSentenceTranslate.Parent = UUID
				sentenceTranslate, err := obj.AppService.SentenceTranslateService.FindOne(ctx, paramsSentenceTranslate)
				if err != nil {
					//commons.Info(ctx, "Error:"+err.Error())
					continue
				}
				var contentTr string
				if helpers.Empty(sentenceTranslate) {
					var resultChi CourseDataModle
					model := mgdb.OnlineClient.Database(KUYU).Collection("course_"+strings.ToLower("chi")).FindOne(ctx, bson.M{
						//"data_area": "cn",
						"level":          course.Level,
						"unit":           course.Unit,
						"module":         course.Module,
						"slide_type":     course.SlideType,
						"form_show_type": course.FormShowType,
						//"eng_content": course.Sentence,
					})
					err = model.Decode(&resultChi)
					if err != nil {
						//commons.Info(ctx, "Error:"+err.Error())
						continue
					}

					if helpers.Empty(resultChi.Sentence){
						model := mgdb.OnlineClient.Database(KUYU).Collection("course_"+strings.ToLower("chi")).FindOne(ctx, bson.M{
							//"data_area": "cn",
							//"level":          course.Level,
							//"unit":           course.Unit,
							//"module":         course.Module,
							//"slide_type":     course.SlideType,
							//"form_show_type": course.FormShowType,
							"eng_content": course.Sentence,
						})
						err = model.Decode(&resultChi)
						if err != nil {
							//commons.Info(ctx, "Error:"+err.Error())
							continue
						}
					}

					// 添加新的eng翻译数据
					var addChiSentenceTranslate requests.SentenceTranslate
					addChiSentenceTranslate.From = "eng"
					addChiSentenceTranslate.To = "chi"
					addChiSentenceTranslate.Parent = UUID
					addChiSentenceTranslate.ContentTr = resultChi.Sentence

					_, err = obj.AppService.SentenceTranslateService.AddSentenceTranslate(ctx, addChiSentenceTranslate)
					if err != nil {
						//commons.Info(ctx, "Error:"+err.Error())
						continue
					}
					contentTr = resultChi.Sentence
				} else {
					contentTr = sentenceTranslate.ContentTr
				}

				// 添加翻译
				var addParamsSentenceTranslate requests.SentenceTranslate
				addParamsSentenceTranslate.From = FROM
				addParamsSentenceTranslate.To = "chi"
				addParamsSentenceTranslate.Parent = sentenceUuid
				addParamsSentenceTranslate.ContentTr = contentTr
				_, err = obj.AppService.SentenceTranslateService.AddSentenceTranslate(ctx, addParamsSentenceTranslate)
				if err != nil {
					//commons.Info(ctx, "Error:"+err.Error())
					continue
				}
			}
		}

		lock.Unlock()
	}
}

func IsWords(str string) bool {
	str = strings.TrimSpace(str)
	start := len(str)
	// 去除空格
	str = strings.Replace(str, " ", "", -1)
	// 去除换行符
	str = strings.Replace(str, "\n", "", -1)
	end := len(str)
	if start == end {
		return true
	} else {
		return false
	}
}
