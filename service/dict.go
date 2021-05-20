package service

import (
	"bufio"
	"context"
	"editorApi/repository"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DictService struct {
}

func (s *DictService) List(ctx *gin.Context, params requests.DictListRequests) (result []responses.DictResponse, err error) {
	var dict repository.Dict
	result, err = dict.DictList(ctx, params)

	if err != nil {
		return nil, err
	}
	return
}

func (s *DictService) FindOne(ctx context.Context, params requests.DictDetailRequests) (result responses.DictResponse, err error) {
	var dict repository.Dict
	result, err = dict.Detail(ctx, params)
	return
}

func (s *DictService) Detail(ctx *gin.Context, params requests.DictDetailRequests) (result responses.DictResponse, err error) {
	var dict repository.Dict
	result, err = dict.Detail(ctx, params)

	var paramsDictTranslateDetailRequests requests.DictTranslateDetailRequests
	paramsDictTranslateDetailRequests.From = params.From
	paramsDictTranslateDetailRequests.To = params.To
	paramsDictTranslateDetailRequests.Parent = params.Uuid
	if err != nil {
		return result, err
	}

	var dictTranslateDetailModel repository.DictTranslate
	translateDetail, err := dictTranslateDetailModel.Detail(ctx, paramsDictTranslateDetailRequests)
	if err != nil {
		return result, err
	}
	result.DictTranslate = translateDetail

	var phraseAllRequests requests.PhraseAllRequests
	phraseAllRequests.From = params.From
	phraseAllRequests.To = params.To
	var dictUuid []string
	dictUuid = append(dictUuid, params.Uuid)
	phraseAllRequests.DictUuid = dictUuid

	var phraseModel repository.Phrase
	phrase, err := phraseModel.FindAll(ctx, phraseAllRequests)
	if err != nil {
		return result, err
	}
	result.Phrase = phrase

	return
}

func (s *DictService) Update(ctx *gin.Context, params requests.DictUpdateRequests) (id interface{}, err error) {
	var dict repository.Dict
	id, err = dict.DictUpdate(ctx, params)
	return
}

func (s *DictService) ReadFile(param *requests.DictHanderParams) {
	if param.FilePath == "" {
		return
	}

	file, err := os.OpenFile(param.FilePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	var dictInfo []string

	// 读取数据
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}
		line = strings.TrimSpace(line)
		dictInfo = append(dictInfo, line)
	}

	for _, words := range dictInfo {
		var dictModel repository.Dict
		if false {
			result, _ := dictModel.DictFindOneByUuid(param.From, helpers.MD5(words))
			if helpers.Empty(result.Uuid) {
				fmt.Println("**************************************单词:" + words + "****************************************")
				s.Collect(words, strings.ToLower(param.From), strings.ToLower(param.To))
			}
		} else {
			fmt.Println("**************************************单词:" + words + "****************************************")
			s.Collect(words, strings.ToLower(param.From), strings.ToLower(param.To))
		}
	}
}

func (s *DictService) AddCarDId(ctx *gin.Context, params requests.DictCardId) {

	var dict repository.Dict

	//card, _ := dict.DictFindOneByCardId(params.From, params.CardId)
	card, _ := dict.DictFindOneByUuid(params.From, helpers.MD5(strings.TrimSpace(params.Content)))

	if helpers.Empty(card.Uuid) {
		dict.DictAddCardId(ctx, params)
	} else {

		if !helpers.Empty(card.SoundInfos) {
			for _, soundInfo := range card.SoundInfos {
				for key, item := range params.SoundInfos {
					if strings.Contains(item.Sound, "ENG-UK-W") && strings.ToLower(soundInfo.Ct) == "en" {
						params.SoundInfos[key].Ct = soundInfo.Ct
						params.SoundInfos[key].Ps = soundInfo.Ps
					}

					if strings.Contains(item.Sound, "ENG-US-M") && strings.ToLower(soundInfo.Ct) == "us" {
						params.SoundInfos[key].Ct = soundInfo.Ct
						params.SoundInfos[key].Ps = soundInfo.Ps
					}
				}
			}
		}

		dict.DictUpdateCardId(ctx, params)
		//result, _ := dict.DictFindOneByUuid(params.From, params.Uuid)
		//if helpers.Empty(result.Uuid) {
		//
		//} else {
		//	commons.Info(ctx, "数据冲突："+result.Uuid)
		//}
	}

}

func (s *DictService) AddDict(ctx context.Context, params requests.Dict)(id interface{}, err error) {
	var model repository.Dict
	return model.AddDict(params)
}

func (s *DictService) Collect(words string, from string, to string) {
	// 单词为空时退出
	if helpers.Empty(words) {
		return
	}
	var lock = sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()
	res, err := http.Get("http://dict.youdao.com/w/" + words)
	if err != nil {
		log.Println(err)
	}
	if res.StatusCode != 200 {
		log.Println("status code error: %d %s", res.StatusCode, res.Status)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
	}

	var dict requests.Dict
	var wordUuid string
	dict.From = from

	fmt.Println("**************************************单词****************************************")
	doc.Find("#results").Each(func(i int, s *goquery.Selection) {
		keyword := s.Find("span[class=keyword]").Text()
		if helpers.Empty(keyword) {
			dict.Content = words
			dict.IsDel = true
			wordUuid = helpers.MD5(words)
		} else {
			dict.Content = keyword
			dict.IsDel = false
			wordUuid = helpers.MD5(keyword)
		}
	})
	dict.Uuid = wordUuid

	fmt.Println("**************************************单词音频****************************************")
	var soundInfos []requests.SoundInfos
	doc.Find("#results [class=pronounce]").Each(func(i int, s *goquery.Selection) {
		keyword := s.Find("span[class=phonetic]").Text()
		var tmp requests.SoundInfos
		if i == 0 {
			tmp.Ct = "en"
		} else {
			tmp.Ct = "us"
		}
		tmp.Ps = keyword
		soundInfos = append(soundInfos, tmp)
	})

	wordsSound := collectPhonetic(words)
	if !helpers.Empty(wordsSound) {
		for _, soundInfo := range soundInfos {
			for key, item := range wordsSound {
				if item.Ct == soundInfo.Ct {
					wordsSound[key].Ps = soundInfo.Ps
				}
			}
		}
		dict.SoundInfos = wordsSound
	} else {
		dict.SoundInfos = soundInfos
	}

	var model repository.Dict
	model.AddDict(dict)

	var dictTranslate requests.DictTranslate
	dictTranslate.Parent = wordUuid
	dictTranslate.From = from
	dictTranslate.To = to
	fmt.Println("**************************************单词拓展****************************************")
	doc.Find("#results").Each(func(i int, s *goquery.Selection) {
		additional := s.Find("p[class=additional]").First().Text()
		// fmt.Println(additional)
		additionalData := strings.ReplaceAll(strings.ReplaceAll(additional, "\n                   ", ""), "\n        ", " ")

		if strings.Contains(additionalData, "[") {
			dictTranslate.Expansion = additionalData
		}
	})

	fmt.Println("**************************************单词词义****************************************")
	var contentTr []requests.ContentTr
	doc.Find("#results div[class=trans-container] ul").First().Each(func(i int, s *goquery.Selection) {
		s.Find("li").Each(func(i int, selection *goquery.Selection) {
			var tmp requests.ContentTr
			tr := strings.Split(selection.First().Text(), ".")
			if len(tr) >= 2 {
				tmp.Cx = tr[0] + "."
				tmp.Content = strings.TrimSpace(tr[1])
				contentTr = append(contentTr, tmp)
			}
		})
	})
	dictTranslate.ContentTr = contentTr

	fmt.Println("**************************************近义词****************************************")
	var synonym []requests.Synonym
	doc.Find("#results #synonyms ul").Each(func(i int, s *goquery.Selection) {
		s.Find("li").Each(func(i int, selection *goquery.Selection) {
			var synonymTmp requests.Synonym
			tmp := selection.First().Text()
			data := strings.Split(tmp, ".")
			if len(data) >= 2 {
				synonymTmp.Cx = data[0] + "."
				synonymTmp.Content = strings.TrimSpace(data[1])
				s.Find("p").Each(func(j int, selection *goquery.Selection) {
					if j == i {
						synonymTmp.ContentTr = strings.ReplaceAll(helpers.TrimHtml(selection.Find("span").Text()), "\n", "")
					}
				})
				synonym = append(synonym, synonymTmp)
			}
		})
	})
	dictTranslate.Synonym = synonym

	fmt.Println("**************************************同根词****************************************")
	var homonyms []requests.Homonyms
	doc.Find("#results #relWordTab").Each(func(i int, s *goquery.Selection) {
		arr := strings.Split(helpers.TrimHtml(s.First().Text()), "\n")
		var tmpWordAttr requests.WordAttr
		var data []requests.WordAttr
		var head string
		var heads []string
		var next int
		for key, item := range arr {
			if key > 1 {
				if strings.ContainsAny(item, ".") {
					head = item
					heads = append(heads, head)
					next = 1
				} else {
					if next%2 == 0 {
						tmpWordAttr.ContentTr = item
						tmpWordAttr.Cx = head
						data = append(data, tmpWordAttr)
					} else {
						tmpWordAttr.Content = item
					}
					next = next + 1
				}
			}
		}

		for _, item := range heads {
			var tmpHomonyms requests.Homonyms
			tmpHomonyms.Cx = item
			var tmp []requests.WordAttr
			for _, value := range data {
				var wordAttr requests.WordAttr
				if value.Cx == item {
					wordAttr.ContentTr = value.ContentTr
					wordAttr.Content = value.Content
					tmp = append(tmp, wordAttr)
				}
			}
			tmpHomonyms.Attrs = tmp
			homonyms = append(homonyms, tmpHomonyms)
		}
	})
	dictTranslate.Homonyms = homonyms

	var dictTranslateModel repository.DictTranslate
	dictTranslateModel.AddDictTranslate(dictTranslate)

	fmt.Println("**************************************单词一短语****************************************")
	var modelPhrase repository.Phrase
	var modelPhraseTranslate repository.PhraseTranslate
	doc.Find("#results #webPhrase").Each(func(i int, s *goquery.Selection) {
		s.Find("p").Each(func(i int, selection *goquery.Selection) {
			arr := strings.Split(helpers.TrimHtml(selection.First().Text()), "\n")
			if len(arr) >= 2 {
				// 添加短语
				phraseContent := arr[0]
				phraseUuid := helpers.MD5(phraseContent)
				var phrase requests.Phrase
				phrase.Content = phraseContent
				phrase.From = from
				phrase.Uuid = phraseUuid
				phrase.DictUuid = []string{wordUuid}

				modelPhrase.AddPhrase(phrase)

				// 添加短语翻译
				var phraseTranslate requests.PhraseTranslate
				phraseTranslate.From = from
				phraseTranslate.To = to
				phraseTranslate.Parent = phraseUuid
				phraseTranslate.ContentTr = arr[1]

				modelPhraseTranslate.AddPhraseTranslate(phraseTranslate)
			}
		})
	})

	fmt.Println("**************************************单词二短语****************************************")
	doc.Find("#results #wordGroup").Each(func(i int, s *goquery.Selection) {
		s.Find("p").Each(func(i int, selection *goquery.Selection) {
			arr := strings.Split(helpers.TrimHtml(selection.First().Text()), "\n")
			if len(arr) >= 2 {
				// 添加短语
				phraseContent := arr[0]
				phraseUuid := helpers.MD5(phraseContent)
				var phrase requests.Phrase
				phrase.Content = phraseContent
				phrase.From = from
				phrase.Uuid = phraseUuid
				phrase.DictUuid = []string{wordUuid}

				modelPhrase.AddPhrase(phrase)

				// 添加短语翻译
				var phraseTranslate requests.PhraseTranslate
				phraseTranslate.From = from
				phraseTranslate.To = to
				phraseTranslate.Parent = phraseUuid
				phraseTranslate.ContentTr = arr[1]

				modelPhraseTranslate.AddPhraseTranslate(phraseTranslate)
			}
		})
	})

	var modelSentence repository.Sentence
	var modelSentenceTranslate repository.SentenceTranslate
	// NAMING1
	fmt.Println("**************************************单词例句****************************************")
	doc.Find("#results #NAMING1").Each(func(i int, s *goquery.Selection) {
		s.Find("li").Each(func(i int, selection *goquery.Selection) {
			arr := strings.Split(helpers.TrimHtml(selection.Find("div[class=examples]").Text()), "\n")
			var sentenceUuid string
			for key, item := range arr {
				if key%2 == 0 {
					sentenceUuid = helpers.MD5(item)
					// 添加句子
					sentenceContent := item
					var sentence requests.Sentence
					sentence.Sentence = sentenceContent
					sentence.From = from
					sentence.Mold = 1
					sentence.Uuid = sentenceUuid
					sentence.Source = "《柯林斯英汉双解大词典》"

					modelSentence.AddSentence(sentence)
				} else {
					// 添加句子翻译
					var sentenceTranslate requests.SentenceTranslate
					sentenceTranslate.From = from
					sentenceTranslate.To = to
					sentenceTranslate.Parent = sentenceUuid
					sentenceTranslate.ContentTr = item

					modelSentenceTranslate.AddSentenceTranslate(sentenceTranslate)
				}
			}
		})
	})

	// 口语
	kouyu_res, err := http.Get("http://dict.youdao.com/example/oral/" + words + "/#keyfrom=dict.sentence.details.kouyu")
	if err != nil {
		log.Println(err)
	}
	if kouyu_res.StatusCode != 200 {
		log.Println("status code error: %d %s", kouyu_res.StatusCode, kouyu_res.Status)
	}
	defer kouyu_res.Body.Close()
	kouyu_doc, err := goquery.NewDocumentFromReader(kouyu_res.Body)
	fmt.Println("**************************************单词口语****************************************")

	kouyu_doc.Find("#results #results-contents #bilingual ul").Each(func(i int, s *goquery.Selection) {
		s.Find("li").Each(func(i int, selection *goquery.Selection) {
			arr := strings.Split(helpers.TrimHtml(selection.First().Text()), "\n")
			dataRel, _ := selection.Find("a").Attr("data-rel")
			arr = append(arr, "http://dict.youdao.com/dictvoice?audio="+dataRel)
			//fmt.Println(arr)
			if len(arr) >= 3 {
				// 添加句子
				sentenceContent := arr[0]
				sentenceUuid := helpers.MD5(sentenceContent)
				var sentence requests.Sentence
				sentence.Sentence = sentenceContent
				sentence.From = from
				sentence.Mold = 1
				sentence.Uuid = sentenceUuid
				sentence.Source = arr[2]

				modelSentence.AddSentence(sentence)

				// 添加句子翻译
				var sentenceTranslate requests.SentenceTranslate
				sentenceTranslate.From = from
				sentenceTranslate.To = to
				sentenceTranslate.Parent = sentenceUuid
				sentenceTranslate.ContentTr = arr[1]

				modelSentenceTranslate.AddSentenceTranslate(sentenceTranslate)
			}
		})
	})

	// 书面语
	shumian_res, err := http.Get("http://dict.youdao.com/example/written/" + words + "/#keyfrom=dict.sentence.details.shumian")
	if err != nil {
		log.Println(err)
	}
	if shumian_res.StatusCode != 200 {
		log.Println("status code error: %d %s", shumian_res.StatusCode, shumian_res.Status)
	}
	defer kouyu_res.Body.Close()
	shumian_doc, err := goquery.NewDocumentFromReader(shumian_res.Body)
	fmt.Println("**************************************单词书面语****************************************")
	shumian_doc.Find("#results #results-contents #bilingual ul").Each(func(i int, s *goquery.Selection) {
		s.Find("li").Each(func(i int, selection *goquery.Selection) {
			arr := strings.Split(helpers.TrimHtml(selection.First().Text()), "\n")
			dataRel, _ := selection.Find("a").Attr("data-rel")
			arr = append(arr, "http://dict.youdao.com/dictvoice?audio="+dataRel)
			//fmt.Println(arr)
			if len(arr) >= 3 {
				// 添加句子
				sentenceContent := arr[0]
				sentenceUuid := helpers.MD5(sentenceContent)
				var sentence requests.Sentence
				sentence.Sentence = sentenceContent
				sentence.From = from
				sentence.Mold = 2
				sentence.Uuid = sentenceUuid
				sentence.Source = arr[2]

				modelSentence.AddSentence(sentence)

				// 添加句子翻译
				var sentenceTranslate requests.SentenceTranslate
				sentenceTranslate.From = from
				sentenceTranslate.To = to
				sentenceTranslate.Parent = sentenceUuid
				sentenceTranslate.ContentTr = arr[1]

				modelSentenceTranslate.AddSentenceTranslate(sentenceTranslate)
			}
		})
	})
	// 每采集一次休息一秒
	time.Sleep(time.Second * 1)
}

func collectPhonetic(words string) (soundInfos []requests.SoundInfos) {
	res, err := http.Get("http://dict.cn/" + words)
	if err != nil {
		log.Println(err)
	}
	if res.StatusCode != 200 {
		log.Println("status code error: %d %s", res.StatusCode, res.Status)
		return
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("**************************************单词音频****************************************")
	var lang, sex string
	wordDir := "eng"

	doc.Find("#content [class=main]").Each(func(i int, s *goquery.Selection) {
		s.Find("div[class=phonetic]").Each(func(i int, selection *goquery.Selection) {
			selection.Find("span").Each(func(index int, span *goquery.Selection) {
				bdo := span.Find("bdo").Text()
				span.Find("i").Each(func(k int, a *goquery.Selection) {
					if index == 0 {
						lang = "en"
					} else {
						lang = "us"
					}
					title, _ := a.Attr("title")
					if title == "女生版发音" {
						sex = "female"
					} else {
						sex = "male"
					}
					naudio, _ := a.Attr("naudio")
					url := "http://audio.dict.cn/" + naudio
					filename := helpers.MD5(words)
					dir := "dict/" + wordDir + "/" + lang + "/" + sex + "/"

					var tmp requests.SoundInfos
					tmp.Ps = bdo
					tmp.Ct = lang
					tmp.Sound = dir + filename + ".mp3?time=" + strconv.FormatInt(time.Now().Unix(), 10)
					tmp.Gender = sex
					soundInfos = append(soundInfos, tmp)
					Download(url, "/opt/data/goPro/editorAPILinux/data/"+dir, filename+".mp3")
				})
			})
		})
	})
	return
}

func Download(wordsUrl string, filePath string, filename string) {
	client := http.DefaultClient
	client.Timeout = time.Second * 60 //设置超时时间
	resp, err := client.Get(wordsUrl)
	if err != nil {
		panic(err)
	}
	if resp.ContentLength <= 0 {
		log.Println("[*] Destination server does not support breakpoint download.")
	}
	raw := resp.Body
	defer raw.Close()
	reader := bufio.NewReaderSize(raw, 1024*32)

	os.MkdirAll(filePath, os.ModePerm)
	file, err := os.Create(filePath + filename)
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(file)

	buff := make([]byte, 32*1024)
	written := 0
	go func() {
		for {
			nr, er := reader.Read(buff)
			if nr > 0 {
				nw, ew := writer.Write(buff[0:nr])
				if nw > 0 {
					written += nw
				}
				if ew != nil {
					err = ew
					break
				}
				if nr != nw {
					err = io.ErrShortWrite
					break
				}
			}
			if er != nil {
				if er != io.EOF {
					err = er
				}
				break
			}
		}
		if err != nil {
			// panic(err)
		}
	}()

	spaceTime := time.Second * 1
	ticker := time.NewTicker(spaceTime)
	lastWtn := 0
	stop := false

	for {
		select {
		case <-ticker.C:
			speed := written - lastWtn
			log.Printf("[*] Speed %s / %s \n", bytesToSize(speed), spaceTime.String())
			if written-lastWtn == 0 {
				ticker.Stop()
				stop = true
				break
			}
			lastWtn = written
		}
		if stop {
			break
		}
	}
}

func bytesToSize(length int) string {
	var k = 1024 // or 1024
	var sizes = []string{"Bytes", "KB", "MB", "GB", "TB"}
	if length == 0 {
		return "0 Bytes"
	}
	i := math.Floor(math.Log(float64(length)) / math.Log(float64(k)))
	r := float64(length) / math.Pow(float64(k), i)
	return strconv.FormatFloat(r, 'f', 3, 64) + " " + sizes[int(i)]
}

func (s *DictService) DictAddTag(ctx *gin.Context, params requests.DictAddTag) {
	xlFile, err := xlsx.OpenFile(params.FilePath)
	if err != nil {
		panic(err)
	}
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
			}
			lock.Unlock()
		}
	}
}

func (s *DictService) DictSearch(ctx *gin.Context, params requests.DictSearch) {

}