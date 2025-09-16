package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 全局日志实例
var GlobalLogger *zap.Logger

// 初始化日志
func InitLogger() error {
	// 创建日志配置
	config := zap.NewProductionConfig()

	// 根据环境设置日志级别
	if gin.Mode() == gin.DebugMode {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.Development = true
	} else {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// 设置日志输出格式
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 创建日志目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 设置日志文件
	config.OutputPaths = []string{
		"logs/app.log",
		"stdout",
	}
	config.ErrorOutputPaths = []string{
		"logs/error.log",
		"stderr",
	}

	// 构建日志器
	logger, err := config.Build()
	if err != nil {
		return fmt.Errorf("构建日志器失败: %v", err)
	}

	// 替换全局日志器
	GlobalLogger = logger
	zap.ReplaceGlobals(logger)

	log.Println("日志系统初始化成功")
	return nil
}

// 接收gin框架的路由日志
func ZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算请求耗时
		cost := time.Since(start)

		// 记录日志
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("cost", cost),
		}

		// 如果有错误，添加错误信息
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()))
		}

		// 根据状态码选择日志级别
		if c.Writer.Status() >= 400 {
			GlobalLogger.Error("HTTP Request", fields...)
		} else {
			GlobalLogger.Info("HTTP Request", fields...)
		}
	}
}

// 获取全局日志实例
func GetLogger() *zap.Logger {
	if GlobalLogger == nil {
		// 如果日志器未初始化，使用默认配置
		logger, _ := zap.NewProduction()
		return logger
	}
	return GlobalLogger
}

// 同步日志缓冲区
func SyncLogger() {
	if GlobalLogger != nil {
		GlobalLogger.Sync()
	}
}
