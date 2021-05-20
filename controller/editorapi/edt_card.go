package editorapi

import (
	"bufio"
	"editorApi/commons"
	"editorApi/repository"
	"editorApi/requests"
	"editorApi/service"
	"editorApi/tools/helpers"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/tealeg/xlsx"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Langs struct {
	Key    int
	From   string
	Male   string
	Female string
}

var dataLangs []Langs

// 多语卡
func ImportCard(ctx *gin.Context) {

	/**
	dataLangs = append(dataLangs, Langs{1, "ALB", "", ""})
	dataLangs = append(dataLangs, Langs{2, "ARA", "", ""})
	dataLangs = append(dataLangs, Langs{4, "AZE", "", ""})
	dataLangs = append(dataLangs, Langs{5, "GLE", "", ""})
	dataLangs = append(dataLangs, Langs{6, "EST", "", ""})
	dataLangs = append(dataLangs, Langs{9, "BUL", "", ""})
	dataLangs = append(dataLangs, Langs{10, "ICE", "", ""})
	dataLangs = append(dataLangs, Langs{12, "BOS", "", ""})
	dataLangs = append(dataLangs, Langs{13, "FAR", "", ""})
	dataLangs = append(dataLangs, Langs{19, "FIL", "", ""})
	dataLangs = append(dataLangs, Langs{20, "FIN", "", ""})
	dataLangs = append(dataLangs, Langs{22, "KHM", "", ""})
	dataLangs = append(dataLangs, Langs{23, "GLG", "", ""})
	dataLangs = append(dataLangs, Langs{24, "GUJ", "", ""})
	dataLangs = append(dataLangs, Langs{25, "KAZ", "", ""})
	dataLangs = append(dataLangs, Langs{28, "HAU", "", ""})
	dataLangs = append(dataLangs, Langs{33, "CZE", "", ""})
	dataLangs = append(dataLangs, Langs{34, "KAN", "", ""})
	dataLangs = append(dataLangs, Langs{36, "HRV", "", ""})
	dataLangs = append(dataLangs, Langs{38, "LAT", "", ""})
	dataLangs = append(dataLangs, Langs{39, "LAV", "", ""})
	dataLangs = append(dataLangs, Langs{40, "LAO", "", ""})
	dataLangs = append(dataLangs, Langs{41, "LIT", "", ""})
	dataLangs = append(dataLangs, Langs{45, "MLT", "", ""})
	dataLangs = append(dataLangs, Langs{46, "MAR", "", ""})
	dataLangs = append(dataLangs, Langs{47, "MAL", "", ""})
	dataLangs = append(dataLangs, Langs{48, "MAY", "", ""})
	dataLangs = append(dataLangs, Langs{49, "MAC", "", ""})
	dataLangs = append(dataLangs, Langs{51, "MOG", "", ""})
	dataLangs = append(dataLangs, Langs{52, "BEN", "", ""})
	dataLangs = append(dataLangs, Langs{53, "BUR", "", ""})
	dataLangs = append(dataLangs, Langs{57, "NEP", "", ""})
	dataLangs = append(dataLangs, Langs{58, "NOR", "", ""})
	dataLangs = append(dataLangs, Langs{59, "PUN", "", ""})
	dataLangs = append(dataLangs, Langs{61, "PUS", "", ""})
	dataLangs = append(dataLangs, Langs{64, "SWE", "", ""})
	dataLangs = append(dataLangs, Langs{66, "SER", "", ""})
	dataLangs = append(dataLangs, Langs{68, "SIN", "", ""})
	dataLangs = append(dataLangs, Langs{69, "ESP", "", ""})
	dataLangs = append(dataLangs, Langs{70, "SLO", "", ""})
	dataLangs = append(dataLangs, Langs{71, "SLV", "", ""})
	dataLangs = append(dataLangs, Langs{76, "TGK", "", ""})
	dataLangs = append(dataLangs, Langs{77, "TEL", "", ""})
	dataLangs = append(dataLangs, Langs{78, "TAM", "", ""})
	dataLangs = append(dataLangs, Langs{80, "TUR", "", ""})
	dataLangs = append(dataLangs, Langs{82, "URD", "", ""})
	dataLangs = append(dataLangs, Langs{83, "UKR", "", ""})
	dataLangs = append(dataLangs, Langs{84, "UZB", "", ""})
	dataLangs = append(dataLangs, Langs{86, "HEB", "", ""})
	dataLangs = append(dataLangs, Langs{87, "GRE", "", ""})
	dataLangs = append(dataLangs, Langs{89, "SND", "", ""})
	dataLangs = append(dataLangs, Langs{90, "HUN", "", ""})
	dataLangs = append(dataLangs, Langs{92, "HYE", "", ""})
	dataLangs = append(dataLangs, Langs{96, "HIN", "", ""})
	dataLangs = append(dataLangs, Langs{100, "YOR", "", ""})
	dataLangs = append(dataLangs, Langs{101, "VIE", "", ""})

	dataLangs = append(dataLangs, Langs{11, "POL", "PL-M", "PL-W"})
	dataLangs = append(dataLangs, Langs{15, "DAN", "", "DA-W"})
	dataLangs = append(dataLangs, Langs{16, "GER", "DE-M", "DE-W"})
	dataLangs = append(dataLangs, Langs{17, "RUS", "", "RU-W"})
	dataLangs = append(dataLangs, Langs{18, "FRE", "FR-M", "FR-W"})

	dataLangs = append(dataLangs, Langs{29, "DUT", "NL-M", "NL-W"})
	dataLangs = append(dataLangs, Langs{43, "RUM", "", "RO-W"})
	dataLangs = append(dataLangs, Langs{60, "POR", "PT-M", "PT-W"})
	dataLangs = append(dataLangs, Langs{85, "SPA", "ES-M", "ES-W"})
	dataLangs = append(dataLangs, Langs{94, "ITA", "IT-M", "IT-W"})

	dataLangs = append(dataLangs, Langs{63, "JPN","JP-M-","JP-W-"})
	dataLangs = append(dataLangs, Langs{27, "KOR","KO-M-","KO-W-"})
	dataLangs = append(dataLangs, Langs{102, "CHS","CN-M","CN-W-"})
	dataLangs = append(dataLangs, Langs{103, "CHI","CN-M","CN-W-"})
	**/

	/**

	**/

	dataLangs = append(dataLangs, Langs{0, "ENG", "", ""})
	dataLangs = append(dataLangs, Langs{3, "AMH", "", ""})
	dataLangs = append(dataLangs, Langs{7, "BAS", "", ""})
	dataLangs = append(dataLangs, Langs{8, "BEL", "", ""})
	dataLangs = append(dataLangs, Langs{14, "AFK", "", ""})
	dataLangs = append(dataLangs, Langs{21, "FRS", "", ""})
	dataLangs = append(dataLangs, Langs{26, "HCR", "", ""})
	dataLangs = append(dataLangs, Langs{30, "KGH", "", ""})
	dataLangs = append(dataLangs, Langs{31, "GAL", "", ""})
	dataLangs = append(dataLangs, Langs{32, "CAT", "", ""})
	dataLangs = append(dataLangs, Langs{35, "CRS", "", ""})
	dataLangs = append(dataLangs, Langs{37, "KRD", "", ""})
	dataLangs = append(dataLangs, Langs{42, "LXB", "", ""})
	dataLangs = append(dataLangs, Langs{44, "MLG", "", ""})
	dataLangs = append(dataLangs, Langs{50, "MRI", "", ""})
	dataLangs = append(dataLangs, Langs{54, "CML", "", ""})
	dataLangs = append(dataLangs, Langs{55, "HAS", "", ""})
	dataLangs = append(dataLangs, Langs{56, "ZUL", "", ""})
	dataLangs = append(dataLangs, Langs{62, "CHW", "", ""})
	dataLangs = append(dataLangs, Langs{65, "SMA", "", ""})
	dataLangs = append(dataLangs, Langs{67, "SOT", "", ""})
	dataLangs = append(dataLangs, Langs{72, "SAM", "", ""})
	dataLangs = append(dataLangs, Langs{73, "GLI", "", ""})
	dataLangs = append(dataLangs, Langs{74, "CBU", "", ""})
	dataLangs = append(dataLangs, Langs{79, "THA", "", ""})
	dataLangs = append(dataLangs, Langs{81, "WSH", "", ""})
	dataLangs = append(dataLangs, Langs{88, "HWI", "", ""})
	dataLangs = append(dataLangs, Langs{91, "SNA", "", ""})
	dataLangs = append(dataLangs, Langs{93, "IBO", "", ""})
	dataLangs = append(dataLangs, Langs{95, "YDS", "", ""})
	dataLangs = append(dataLangs, Langs{97, "IXT", "", ""})
	dataLangs = append(dataLangs, Langs{98, "IND", "", ""})
	dataLangs = append(dataLangs, Langs{99, "JAV", "", ""})

	commons.Info(ctx, "处理中.........")
	excelFileName := "data/cart/多语卡数据.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		panic(err)
	}
	for _, sheet := range xlFile.Sheets {
		for key, row := range sheet.Rows {
			if key == 0 {

			} else {
				var words string
				for j, cell := range row.Cells {
					if j == 0 {
						words = strings.TrimSpace(cell.String())
					}

					if strings.Contains(words, " ") || strings.Contains(words, "-") {
						continue
					}

					for _, itme := range dataLangs {
						if itme.Key == j {
							var requestsDictCartId requests.DictCardId

							var content = strings.TrimSpace(cell.String())
							requestsDictCartId.CardId = helpers.MD5(words)
							requestsDictCartId.Uuid = helpers.MD5(content)
							requestsDictCartId.From = itme.From
							requestsDictCartId.Content = content

							var SoundInfos []requests.SoundInfos

							if !helpers.Empty(itme.Male) {
								maleSound := "dict/" + itme.Male + "/" + itme.Male + "_" + RandomCompletion(key) + ".mp3?time=" + strconv.FormatInt(time.Now().Unix(), 10)
								SoundInfos = append(SoundInfos, requests.SoundInfos{"us", "", maleSound, "male"})
							}

							if !helpers.Empty(itme.Female) {
								femaleSound := "dict/" + itme.Female + "/" + itme.Female + "_" + RandomCompletion(key) + ".mp3?time=" + strconv.FormatInt(time.Now().Unix(), 10)
								SoundInfos = append(SoundInfos, requests.SoundInfos{"en", "", femaleSound, "female"})
							}

							requestsDictCartId.SoundInfos = SoundInfos
							obj := service.AppService()
							obj.AppService.DictService.AddCarDId(ctx, requestsDictCartId)
						}
					}
				}
			}
		}
	}

	commons.Info(ctx, "处理完成.........")

	commons.Success(ctx, nil, "提交成功！", nil)
}

func RandomCompletion(num int) (result string) {
	cat := len(strconv.Itoa(num))
	switch cat {
	case 1:
		result = "000" + strconv.Itoa(num)
	case 2:
		result = "00" + strconv.Itoa(num)
	case 3:
		result = "0" + strconv.Itoa(num)
	case 4:
		result = strconv.Itoa(num)
	}
	return result
}

func DeleteSentence(ctx *gin.Context) {
	commons.Info(ctx, "处理中.........")
	excelFileName := "data/repetition/repetition.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		panic(err)
	}
	obj := service.AppService()
	var lock = sync.Mutex{}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			lock.Lock()
			var request requests.SentenceDelete
			request.From = "eng"
			for j, cell := range row.Cells {
				if j == 0 {
					request.Uuid = cell.String()
				}
				fmt.Println(request)
				if j == 1 {
					num, _ := strconv.Atoi(cell.String())
					for i := num; i > 1; i-- {
						obj.AppService.SentenceService.DeleteOne(ctx, request)
					}
				}
			}
			lock.Unlock()
		}
	}
	commons.Info(ctx, "处理完成.........")
	commons.Success(ctx, nil, "提交成功！", nil)
}

func ReadFiles(ctx *gin.Context) {

	files, _ := ioutil.ReadDir("/opt/data/goPro/editorAPILinux/data/dict/eng/en/female")
	obj := service.AppService()
	for _, onefile := range files {
		if onefile.IsDir() {
			fmt.Println(onefile.Name(), "目录:")
		} else {
			if onefile.Size() <= 0 {
				f, err := os.OpenFile("data/exports/test3.txt", os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {

				}
				defer func() {
					if err = f.Close(); err != nil {

					}
				}()

				var uuid = strings.Replace(onefile.Name(), ".mp3", "", 1)
				var params requests.DictDetailRequests
				params.From = "eng"
				params.Uuid = uuid
				result, err := obj.AppService.DictService.FindOne(ctx, params)
				if !helpers.Empty(result.Uuid) {
					_, err = fmt.Fprintln(f, result.Content)
				}

				if err != nil {

				}
			}
		}
	}
}

func ReadSoundInfos(ctx *gin.Context) {
	var dict repository.Dict
	var params requests.DictDetailRequests
	params.From = "eng"
	result, _ := dict.SoundInfosSize(ctx, params)

	for _, onefile := range result {
		f, err := os.OpenFile("data/exports/words.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {

		}
		defer func() {
			if err = f.Close(); err != nil {

			}
		}()

		if !helpers.Empty(onefile.Uuid) {
			_, err = fmt.Fprintln(f, onefile.Content)
		}

		if err != nil {

		}
	}
}

func ReadSentence(ctx *gin.Context) {
	commons.Info(ctx, "Sentence 处理中.........")
	var params requests.SentenceFindAll
	params.From = "eng"
	var sentence repository.Sentence

	var page_lock = sync.Mutex{}
	createdOn := time.Now()
	time.Sleep(1000 * time.Millisecond) //休眠1秒
	for page := 1; page <= 263; page++ {
		page_lock.Lock()
		params.CreatedOn = createdOn
		params.Page.PageIndex = 1
		params.Page.PageSize = 5000
		result, err := sentence.FindAll(ctx, params)
		if err != nil {
			commons.Info(ctx, err.Error())
		}

		var lock = sync.Mutex{}
		if !FileExist("data/sentence/sentence_" + strconv.Itoa(page) + ".txt") {
			os.Create("data/sentence/sentence_" + strconv.Itoa(page) + ".txt")
		}
		f, err := os.OpenFile("data/sentence/sentence_"+strconv.Itoa(page)+".txt", os.O_APPEND|os.O_WRONLY, 0777)
		for _, onefile := range result {
			lock.Lock()
			if err != nil {
				commons.Info(ctx, err.Error())
			}

			if !helpers.Empty(onefile.Uuid) {
				_, err = fmt.Fprintln(f, onefile.Sentence+"\r\n")
				createdOn = onefile.CreatedOn
			}

			if err != nil {
				commons.Info(ctx, err.Error())
			}
			lock.Unlock()
		}
		defer func() {
			if err = f.Close(); err != nil {
				commons.Info(ctx, err.Error())
			}
		}()
		page_lock.Unlock()
	}

	commons.Info(ctx, "Sentence 处理完成.........")
	commons.Success(ctx, nil, "Sentence 成功！", nil)
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func ReadWords(ctx *gin.Context) {
	commons.Info(ctx, "ReadWords 处理中.........")
	var params requests.DictFindAll
	params.From = "eng"
	var dict repository.Dict

	var page_lock = sync.Mutex{}
	createdOn := time.Now()
	time.Sleep(1000 * time.Millisecond) //休眠1秒
	for page := 1; page <= 263; page++ {
		page_lock.Lock()
		params.CreatedOn = createdOn
		params.Page.PageIndex = 1
		params.Page.PageSize = 5000
		result, err := dict.FindAll(ctx, params)
		if err != nil {
			commons.Info(ctx, err.Error())
		}

		var lock = sync.Mutex{}
		if !FileExist("data/words/words_" + strconv.Itoa(page) + ".txt") {
			os.Create("data/words/words_" + strconv.Itoa(page) + ".txt")
		}
		f, err := os.OpenFile("data/words/words_"+strconv.Itoa(page)+".txt", os.O_APPEND|os.O_WRONLY, 0777)
		for _, onefile := range result {
			lock.Lock()
			if err != nil {
				commons.Info(ctx, err.Error())
			}

			if !helpers.Empty(onefile.Uuid) {
				_, err = fmt.Fprintln(f, onefile.Content+"\r\n")
				createdOn = onefile.CreatedOn
			}

			if err != nil {
				commons.Info(ctx, err.Error())
			}
			lock.Unlock()
		}
		defer func() {
			if err = f.Close(); err != nil {
				commons.Info(ctx, err.Error())
			}
		}()
		page_lock.Unlock()
	}

	commons.Info(ctx, "ReadWords 处理完成.........")
	commons.Success(ctx, nil, "ReadWords 成功！", nil)
}

// 添加tag
func AddTag(ctx *gin.Context) {
	commons.Info(ctx, "处理中.........")
	excelFileName := "data/repetition/ENG-SEPOxford.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		panic(err)
	}
	var lock = sync.Mutex{}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			lock.Lock()
			var request requests.AddTags
			var dictTranslate repository.DictTranslate
			request.From = "eng"
			for j, cell := range row.Cells {
				if j == 0 {
					var dict repository.Dict
					var dictParams requests.Dict
					dictParams.ListOrder = strings.ToLower(strings.TrimSpace(cell.String()))
					dictParams.From = "eng"
					dictData, _ := dict.DictFindOne(dictParams)
					if !helpers.Empty(dictData.Uuid) {
						request.Parent = dictData.Uuid
						var Tags []requests.Tag
						var tag requests.Tag
						tag.Name = "小学英语（沪教版）"
						tag.Key = "ENG-SEPOxford"
						Tags = append(Tags, tag)
						request.Tags = Tags
						request.From = "eng"
						request.To = "chi"
						dictTranslate.AddTags(ctx, request)
					}
				}
				fmt.Println(request)
			}
			lock.Unlock()
		}
	}
	commons.Info(ctx, "处理完成.........")
	commons.Success(ctx, nil, "提交成功！", nil)
}

func AddSound(ctx *gin.Context) {
	txtFile := ctx.PostForm("name")
	if helpers.Empty(txtFile) {
		commons.Error(ctx, 500, nil, "name 不能为空")
	}
	fileName := "data/words/" + txtFile + ".txt"

	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)

	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	var size = stat.Size()
	fmt.Println("file size=", size)

	buf := bufio.NewReader(file)
	var index = 1
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if line != "" {
			fmt.Println(index)
			fmt.Println(line)

			lock.Lock()
			var request requests.DictSound
			request.From = "eng"
			var dict repository.Dict
			var dictParams requests.Dict
			dictParams.ListOrder = strings.ToLower(strings.TrimSpace(line))
			dictParams.From = "eng"
			dictData, _ := dict.DictFindOne(dictParams)
			if !helpers.Empty(dictData.Uuid) {
				request.Uuid = dictData.Uuid
				var soundInfos []requests.SoundInfos
				for _, item := range dictData.SoundInfos {
					var tmp requests.SoundInfos

					tmp.Sound = item.Sound

					if strings.ToLower(strings.TrimSpace(item.Ct)) == "en" && strings.ToLower(strings.TrimSpace(item.Gender)) == "male" {
						tmp.Sound = "dict/" + txtFile + "/eng/male/" + RandomCompletion(index) + ".mp3?time=1602725870"
					}

					if strings.ToLower(strings.TrimSpace(item.Ct)) == "us" && strings.ToLower(strings.TrimSpace(item.Gender)) == "female" {
						tmp.Sound = "dict/" + txtFile + "/usa/female/" + RandomCompletion(index) + ".mp3?time=1602725870"
					}

					tmp.Gender = item.Gender
					tmp.Ct = item.Ct
					tmp.Ps = item.Ps
					soundInfos = append(soundInfos, tmp)
				}
				request.SoundInfos = soundInfos
				dict.DictUpdateSound(ctx, request)
			}
			index++
			lock.Unlock()
		}

		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}
	}
}

type AeneasParams struct {
	Language string `json:"language"` // language = "eng" https://www.readbeyond.it/aeneas/docs/language.html
	Content  string `json:"content"`
	Plain    string `json:"plain"` //var plain = "mplain" // mplain 每个单词    plain 每个句子 unparsed
}


// @Tags 音频数据打点
// @Summary 音频数据打点
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body AeneasParams true "目录列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"目录列表"}"
// @Router /editor/card/aeneas_job [post]
func AeneasJob(ctx *gin.Context) {
	var param AeneasParams

	param.Language = ctx.PostForm("language")
	param.Content = ctx.PostForm("content")
	param.Plain = ctx.PostForm("plain")

	filename, err := ctx.FormFile("filename")
	if err != nil {
		commons.Error(ctx, 500, err, "文件名不能为空！")
	}
	uuidName := uuid.NewV4().String()
	mp3File := "data/aeneas/" + uuidName + ".mp3"

	if err := ctx.SaveUploadedFile(filename, mp3File); err != nil {
		commons.Error(ctx, 500, err, "文件保存失败！")
	}

	dir, _ := os.Getwd()
	os.MkdirAll("data/aeneas/", os.ModePerm)
	
	resultFile := "data/aeneas/" + uuidName + ".json"
	os.RemoveAll(resultFile)
	os.Create(resultFile)

	contentFile := "data/aeneas/" + uuidName + ".txt"
	os.RemoveAll(contentFile)
	os.Create(contentFile)

	contents := strings.Split(param.Content, "\n")

	file, err := os.Create(contentFile)
	if err != nil {
		fmt.Printf("create map file error: %v\n", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, content := range contents {
		if !helpers.Empty(content) {
			lineStr := fmt.Sprintf("%s", content)
			fmt.Fprintln(writer, lineStr)
		}
	}
	writer.Flush()

	chmod := exec.Command("chmod", "-R", "755", resultFile)
	chmod.Run()

	python := exec.Command("/usr/local/python3.7.5/bin/python3", "-m", "aeneas.tools.execute_task", dir+"/"+mp3File, dir+"/"+contentFile, "task_language="+param.Language+"|os_task_file_format=json|is_text_type="+param.Plain, dir+"/"+resultFile)
	//python := exec.Command("/usr/local/bin/python3", "-m", "aeneas.tools.execute_task", dir+"/"+mp3File, dir+"/"+contentFile, "task_language="+param.Language+"|os_task_file_format=json|is_text_type="+param.Plain, dir+"/"+resultFile)

	err = python.Run()

	// 删除缓存文件
	os.RemoveAll(mp3File)
	os.RemoveAll(contentFile)

	if err != nil {
		fmt.Println(python)
		commons.Success(ctx, python, "提交成功！", nil)
	} else {
		f, err := os.Open(dir + "/" + resultFile)
		if err != nil {
			fmt.Println("read file fail", err)
		}
		defer f.Close()

		fd, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println("read to fd fail", err)
		}
		// 删除结果文件
		os.RemoveAll(resultFile)
		commons.Success(ctx, string(fd), "提交成功！", nil)
	}
}
