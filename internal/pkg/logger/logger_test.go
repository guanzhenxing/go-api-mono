package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestLogger(t *testing.T) {
	// 创建测试日志记录器
	logger := &Logger{
		logger: zaptest.NewLogger(t),
	}

	// 测试基本日志级别
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warning message")
	logger.Error("error message")

	// 测试带字段的日志
	logger.With(
		zap.String("key1", "value1"),
		zap.Int("key2", 2),
	).Info("message with fields")

	// 测试错误日志
	logger.Error("error message",
		zap.Error(ErrInvalidLogLevel),
		zap.String("additional", "info"),
	)
}

func TestLoggerOptions(t *testing.T) {
	opts := LogConfig{
		Level:      "debug",
		Filename:   "test.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	logger, err := New(opts)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	logger.Info("test message")
}
