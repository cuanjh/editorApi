package commons

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	// 代码
	Code int `json:"code" example:"200"`
	// 数据集
	Data interface{} `json:"data"`
	// 消息
	Msg string `json:"msg"`
	// success
	Success bool `json:"success"`
	// request_id
	RequestId string `json:"request_id"`
}

type Page struct {
	List      interface{} `json:"list"`
	Count     int64       `json:"count"`
	PageIndex int64       `json:"pageIndex"`
	PageSize  int64       `json:"pageSize"`
}

func (res *Response) ReturnOK() *Response {
	res.Code = 200
	return res
}

func (res *Response) ReturnError(code int) *Response {
	res.Code = code
	return res
}

// 失败数据处理
func Error(c *gin.Context, code int, err error, msg string) {
	var res Response
	Errors(c, err.Error())
	res.RequestId = c.Request.Header.Get("X-Request-Id")
	res.Msg = err.Error()
	res.Success = false
	if msg != "" {
		res.Msg = msg
	} else {
		res.Msg = "error"
	}
	c.JSON(http.StatusOK, res.ReturnError(code))
	c.Abort()
	panic(nil)
}

// 通常成功数据处理
func Success(c *gin.Context, data interface{}, msg string, request interface{}) {
	var res Response
	ResponseSuccess(c, data, request)
	res.RequestId = c.Request.Header.Get("X-Request-Id")
	res.Data = data
	res.Success = true
	if msg != "" {
		res.Msg = msg
	} else {
		res.Msg = "success"
	}
	c.JSON(http.StatusOK, res.ReturnOK())
	c.Abort()
	panic(nil)
}

func SuccessJsonp(c *gin.Context, data interface{}, msg string, request interface{}) {
	var res Response
	ResponseSuccess(c, data, request)
	res.RequestId = c.Request.Header.Get("X-Request-Id")
	res.Data = data
	res.Success = true
	if msg != "" {
		res.Msg = msg
	} else {
		res.Msg = "success"
	}
	c.JSONP(http.StatusOK, res.ReturnOK())
	c.Abort()
	panic(nil)
}

// 分页数据处理
func PageOK(c *gin.Context, result interface{}, count int64, pageIndex int64, pageSize int64, msg string, request interface{}) {
	var res PageResponse
	ResponseSuccess(c, result, request)
	res.Data.List = result
	res.Data.Count = count
	res.Data.PageIndex = pageIndex
	res.Data.PageSize = pageSize
	if msg != "" {
		res.Msg = msg
	}
	c.JSON(http.StatusOK, res.ReturnOK())
	c.Abort()
	panic(nil)
}

// 兼容函数
func Custum(c *gin.Context, data gin.H) {
	c.JSON(http.StatusOK, data)
}

type PageResponse struct {
	// 代码
	Code int `json:"code" example:"200"`
	// 数据集
	Data Page `json:"data"`
	// 消息
	Msg string `json:"msg"`
}

func (res *PageResponse) ReturnOK() *PageResponse {
	res.Code = 200
	return res
}

type LayPage struct {
	Limit int64 `json:"limit"`
	Skip  int64 `json:"skip"`
}

func DefaultPage() LayPage {
	return LayPage{
		Limit: 50,
		Skip:  0,
	}
}
