package xlsxtool

import (
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

//数据导出excel并下载
func ExportToExcel(c *gin.Context, titleList []string, data []interface{}, fileName string) {
	// 生成一个新的文件
	file := xlsx.NewFile()
	// 添加sheet页
	sheet, _ := file.AddSheet("Sheet1")
	// 插入表头
	titleRow := sheet.AddRow()
	for _, v := range titleList {
		cell := titleRow.AddCell()
		cell.Value = v
		//表头字体颜色
		cell.GetStyle().Font.Color = "00FF0000"
		//居中显示
		cell.GetStyle().Alignment.Horizontal = "center"
		cell.GetStyle().Alignment.Vertical = "center"
	}
	// 插入内容
	for _, v := range data {
		row := sheet.AddRow()
		row.WriteStruct(v, -1)
	}
	//c.Writer.Header().Set("Content-Type", "application/octet-stream")
	//disposition := fmt.Sprintf("attachment; filename=\"%s-%s.xlsx\"", fileName, time.Now().Format("2006-01-02 15:04:05"))
	//c.Writer.Header().Set("Content-Disposition", disposition)
	//_ = file.Write(c.Writer)

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+"Workbook.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")

	//回写到web 流媒体 形成下载
	_ = file.Write(c.Writer)
}
