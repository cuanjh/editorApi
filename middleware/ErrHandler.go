package middleware

import (
	"editorApi/commons"
	"errors"
	"github.com/gin-gonic/gin"
)

func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录一个错误的日志
				commons.Error(c, 403, errors.New("系统错误！"), "系统错误！")
			}
		}()
		c.Next()
	}
}
