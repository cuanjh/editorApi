package service

import (
	"context"
	"editorApi/repository"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/tealeg/xlsx"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

type StatisticService struct {
}

var mutex sync.Mutex

func (s *StatisticService) StatisticUnlockChapter(ctx *gin.Context, params requests.StatisticUnlockChapter) (result []responses.UnlockInfosChapter, err error) {
	var dict repository.UnlockInfosChapter
	result, err = dict.StatisticUnlockChapter(ctx, params)
	return
}

func (s *StatisticService) HanderStatisticUnlockPart(request requests.StatisticUnlockPart) {
	file := xlsx.NewFile()
	ctx := context.TODO()
	dateData := helpers.GetBetweenDates(request.StartDate, request.EndDate)

	var requestsCourseContentInfos requests.CourseContentInfos
	requestsCourseContentInfos.Course_code = request.CourseCode
	var courseContentInfosModel repository.CourseContentInfos
	courseContentInfos, _ := courseContentInfosModel.FindOne(ctx, requestsCourseContentInfos)

	var catalogsModel repository.Catalogs
	var params requests.Catalogs
	params.Parent_uuid = courseContentInfos.Uuid
	catalogs, _ := catalogsModel.FindAll(ctx, params)
	var titleList []TitleList
	titleList = append(titleList, TitleList{key: 0, name: request.CourseCode + "-L0-C0"})
	for level, catalog := range catalogs {
		level = level + 1
		levelCode := request.CourseCode + "-L" + strconv.Itoa(level) + "-C"
		params.Parent_uuid = catalog.Uuid
		result, _ := catalogsModel.FindAll(ctx, params)
		for chapter, _ := range result {
			chapter = chapter + 1
			chapterCode := levelCode + strconv.Itoa(chapter)
			titleList = append(titleList, TitleList{key: level + chapter, name: chapterCode})
		}
	}
	// 添加sheet页
	sheet, _ := file.AddSheet("统计数据")
	// 插入表头
	titleRow := sheet.AddRow()
	for _, v := range titleList {
		cell := titleRow.AddCell()
		cell.Value = v.name
		//表头字体颜色
		cell.GetStyle().Font.Color = "00FF0000"
		//居中显示
		cell.GetStyle().Alignment.Horizontal = "center"
		cell.GetStyle().Alignment.Vertical = "center"
	}

	var unlockInfosPartModel repository.UnlockInfosPart

	sort.Strings(dateData)
	sort.Sort(titleListSort(titleList))
	for _, dateTime := range dateData {
		row := sheet.AddRow()
		for key, title := range titleList {
			mutex.Lock()
			if key == 0 {
				tmpFormNo := row.AddCell()
				tmpFormNo.SetValue(dateTime)
			} else {
				var tmpParams requests.StatisticUnlockPart
				tmpParams.CourseCode = request.CourseCode
				tmpParams.Chapter = title.name
				tmpParams.StartDate = "2020-09-01 00:00:00"
				tmpParams.EndDate = dateTime + " 23:59:59"

				/**
				var users []string
				statisticUnlockPart, _ := unlockInfosPartModel.StatisticUnlockPart(ctx, tmpParams)
				for _, part := range statisticUnlockPart {
					if part.Chapter == title.name && part.CorrectRate > 0 {
						users = append(users, part.UserId)
					}
				}
				tmpSentence := row.AddCell()
				tmpSentence.SetValue(len(users))
				*/

				//result, _ := unlockInfosPartModel.CountPart(ctx, tmpParams)

				result, _ := unlockInfosPartModel.CountPartUnique(ctx, tmpParams)
				tmpSentence := row.AddCell()
				tmpSentence.SetValue(result)
			}
			time.Sleep(time.Duration(1) * time.Second)
			mutex.Unlock()
		}
	}

	os.MkdirAll("data/statistic/", os.ModePerm)
	filename := fmt.Sprintf("data/statistic/%v", uuid.NewV4().String()+"-"+time.Now().Format("2006-01-02")+".xlsx")
	file.Save(filename)

	var contentExportsModel repository.ContentExports
	var paramsContentExports requests.ContentExports
	paramsContentExports.ID = request.Id
	paramsContentExports.Url = "/editor/" + filename
	contentExportsModel.UpdateContentExports(ctx, paramsContentExports)

}

func (s *StatisticService) HanderStatisticUnlockChapter(request requests.StatisticUnlockChapter) {
	file := xlsx.NewFile()
	ctx := context.TODO()
	dateData := helpers.GetBetweenDates(request.StartDate, request.EndDate)

	var requestsCourseContentInfos requests.CourseContentInfos
	requestsCourseContentInfos.Course_code = request.CourseCode
	var courseContentInfosModel repository.CourseContentInfos
	courseContentInfos, _ := courseContentInfosModel.FindOne(ctx, requestsCourseContentInfos)

	var catalogsModel repository.Catalogs
	var params requests.Catalogs
	params.Parent_uuid = courseContentInfos.Uuid
	catalogs, _ := catalogsModel.FindAll(ctx, params)
	var titleList []TitleList
	titleList = append(titleList, TitleList{key: 0, name: request.CourseCode + "-L0-C0"})
	for level, catalog := range catalogs {
		level = level + 1
		levelCode := request.CourseCode + "-L" + strconv.Itoa(level) + "-C"
		params.Parent_uuid = catalog.Uuid
		result, _ := catalogsModel.FindAll(ctx, params)
		for chapter, _ := range result {
			chapter = chapter + 1
			chapterCode := levelCode + strconv.Itoa(chapter)
			titleList = append(titleList, TitleList{key: level + chapter, name: chapterCode})
		}
	}
	// 添加sheet页
	//sort.Strings(titleList)
	sheet, _ := file.AddSheet("统计数据")
	// 插入表头
	titleRow := sheet.AddRow()
	for _, v := range titleList {
		cell := titleRow.AddCell()
		cell.Value = v.name
		//表头字体颜色
		cell.GetStyle().Font.Color = "00FF0000"
		//居中显示
		cell.GetStyle().Alignment.Horizontal = "center"
		cell.GetStyle().Alignment.Vertical = "center"
	}

	var dict repository.UnlockInfosChapter
	statisticUnlockChapter, _ := dict.StatisticUnlockChapter(ctx, request)

	for _, dateTime := range dateData {
		row := sheet.AddRow()
		sort.Sort(titleListSort(titleList))

		for key, title := range titleList {
			if key == 0 {
				tmpFormNo := row.AddCell()
				tmpFormNo.SetValue(dateTime)
			} else {
				var users []string
				for _, statistic := range statisticUnlockChapter {
					if statistic.Chapter == title.name && statistic.CreatedOn.Format("2006-01-02") == dateTime {
						users = append(users, statistic.UserId)
					}
				}
				tmpSentence := row.AddCell()
				tmpSentence.SetValue(len(users))
			}
		}
	}

	os.MkdirAll("data/statistic/", os.ModePerm)
	filename := fmt.Sprintf("data/statistic/%v", uuid.NewV4().String()+"-"+time.Now().Format("2006-01-02")+".xlsx")
	file.Save(filename)

	//collection := mgdb.MongoClient.Database(mgdb.DbEditor).Collection("content_exports")
	//id := uuid.NewV4().String()
	//var contentExports ContentExports
	//contentExports.ID = id
	//contentExports.CreatedOn = time.Now()
	//contentExports.Status = 2
	//contentExports.Name = request.CourseCode
	//contentExports.Code = request.Code
	//contentExports.UUID = courseContentInfos.Uuid
	//contentExports.Url = "/editor/" + filename
	//collection.InsertOne(ctx, contentExports)
}

// 实现排序
type TitleList struct {
	key  int    `json:"key"`
	name string `json:"name"`
}

type titleListSort []TitleList

func (I titleListSort) Len() int {
	return len(I)
}
func (I titleListSort) Less(i, j int) bool {
	return I[i].key < I[j].key
}
func (I titleListSort) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}
