package commons

import (
	"editorApi/tools/helpers"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var logLevel = zap.NewAtomicLevel()

var logger *zap.Logger

type Logger struct {
	*zap.Logger
}

const (
	// logFormat
	LOGFORMAT_JSON    = "json"
	LOGFORMAT_CONSOLE = "console"

	// EncoderConfig
	TIME_KEY       = "time"
	LEVLE_KEY      = "level"
	NAME_KEY       = "logger"
	CALLER_KEY     = "caller"
	MESSAGE_KEY    = "msg"
	STACKTRACE_KEY = "stacktrace"

	// 日志归档配置项
	// 每个日志文件保存的最大尺寸 单位：M
	MAX_SIZE = 1
	// 文件最多保存多少天
	MAX_BACKUPS = 3
	// 日志文件最多保存多少个备份
	MAX_AGE = 7
)

func ResponseSuccess(c *gin.Context, response, request interface{}) {
	logger := zap.L()
	start := time.Now()
	// some evil middlewares modify this values
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	c.Next()

	request_id := c.Request.Header.Get("X-Request-Id")

	end := time.Now().UTC()
	latency := end.Sub(start)

	if len(c.Errors) > 0 {
		// Append error field if this is an erroneous request.
		for _, e := range c.Errors.Errors() {
			logger.Error(e)
		}
	} else {
		logger.Info(
			"Success",
			zap.String("request_id", request_id),
			zap.String("env", helpers.GetENV()),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("time", end.Format(time.RFC3339)),
			zap.Duration("latency", latency),
			zap.Any("request", request),
			zap.Any("response", response),
		)
	}
}

func Info(c *gin.Context, message string) {
	logger := zap.L()
	request_id := c.Request.Header.Get("X-Request-Id")
	if len(c.Errors) > 0 {
		// Append error field if this is an erroneous request.
		for _, e := range c.Errors.Errors() {
			logger.Error(e)
		}
	} else {
		logger.Info(
			message,
			zap.String("request_id", request_id),
			zap.String("env", helpers.GetENV()),
		)
	}
}

func Errors(c *gin.Context, message string) {
	logger := zap.L()
	start := time.Now()
	// some evil middlewares modify this values
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	c.Next()

	request_id := c.Request.Header.Get("X-Request-Id")

	end := time.Now().UTC()
	latency := end.Sub(start)

	if len(c.Errors) > 0 {
		// Append error field if this is an erroneous request.
		for _, e := range c.Errors.Errors() {
			logger.Error(e)
		}
	} else {
		logger.Error(
			message,
			zap.String("request_id", request_id),
			zap.String("env", helpers.GetENV()),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("time", end.Format(time.RFC3339)),
			zap.Duration("latency", latency),
		)
	}
}

func Panic(c *gin.Context, message string) {
	logger := zap.L()
	start := time.Now()
	// some evil middlewares modify this values
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	c.Next()

	request_id := c.Request.Header.Get("X-Request-Id")

	end := time.Now().UTC()
	latency := end.Sub(start)

	if len(c.Errors) > 0 {
		// Append error field if this is an erroneous request.
		for _, e := range c.Errors.Errors() {
			logger.Error(e)
		}
	} else {
		logger.Panic(
			message,
			zap.String("request_id", request_id),
			zap.String("env", helpers.GetENV()),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("time", end.Format(time.RFC3339)),
			zap.Duration("latency", latency),
		)
	}
}

func Warn(c *gin.Context, message string) {
	logger := zap.L()
	start := time.Now()
	// some evil middlewares modify this values
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	c.Next()

	request_id := c.Request.Header.Get("X-Request-Id")

	end := time.Now().UTC()
	latency := end.Sub(start)

	if len(c.Errors) > 0 {
		// Append error field if this is an erroneous request.
		for _, e := range c.Errors.Errors() {
			logger.Error(e)
		}
	} else {
		logger.Warn(
			message,
			zap.String("request_id", request_id),
			zap.String("env", helpers.GetENV()),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("time", end.Format(time.RFC3339)),
			zap.Duration("latency", latency),
		)
	}
}

// 设置日志级别、输出格式和日志文件的路径
func SetLogs(logLevel zapcore.Level, logFormat, fileName string) {

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        TIME_KEY,
		LevelKey:       LEVLE_KEY,
		NameKey:        NAME_KEY,
		CallerKey:      CALLER_KEY,
		MessageKey:     MESSAGE_KEY,
		StacktraceKey:  STACKTRACE_KEY,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 大写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器(相对路径+行号)
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志输出格式
	var encoder zapcore.Encoder
	switch logFormat {
	case LOGFORMAT_JSON:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 设置默认存放
	if helpers.Empty(fileName) {
		fileName = getFilePath()
	}

	// 添加日志切割归档功能
	hook := lumberjack.Logger{
		Filename:   fileName,    // 日志文件路径
		MaxSize:    MAX_SIZE,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: MAX_BACKUPS, // 日志文件最多保存多少个备份
		MaxAge:     MAX_AGE,     // 文件最多保存多少天
		Compress:   true,        // 是否压缩
		LocalTime:  true,
	}

	core := zapcore.NewCore(
		encoder,                                                                         // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), zapcore.AddSync(&hook)), // 打印到控制台和文件
		zap.NewAtomicLevelAt(logLevel),                                                  // 日志级别
	)

	// 开启文件及行号
	caller := zap.AddCaller()
	// 开启开发模式，堆栈跟踪
	//development := zap.Development()

	// 构造日志
	logger := zap.New(core, caller, zap.AddCallerSkip(1))

	// 将自定义的logger替换为全局的logger
	zap.ReplaceGlobals(logger)
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		//log.Info(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func getFilePath() string {
	logfile := getCurrentDirectory() + "/log/" + getAppname() + ".log"
	return logfile
}

func getAppname() string {
	full := os.Args[0]
	full = strings.Replace(full, "\\", "/", -1)
	splits := strings.Split(full, "/")
	if len(splits) >= 1 {
		name := splits[len(splits)-1]
		name = strings.TrimSuffix(name, ".exe")
		return name
	}
	return ""
}
