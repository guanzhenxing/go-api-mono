package logger

import (
	"errors"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// ErrInvalidLogLevel 表示无效的日志级别
	ErrInvalidLogLevel = errors.New("invalid log level")
)

// Logger 封装了zap.Logger
type Logger struct {
	logger *zap.Logger
}

// LogConfig 定义了日志配置选项
type LogConfig struct {
	Level      string `yaml:"level"`      // 日志级别
	Filename   string `yaml:"filename"`   // 日志文件名
	MaxSize    int    `yaml:"maxSize"`    // 单个日志文件最大大小（MB）
	MaxBackups int    `yaml:"maxBackups"` // 最大备份文件数
	MaxAge     int    `yaml:"maxAge"`     // 最大保留天数
	Compress   bool   `yaml:"compress"`   // 是否压缩
}

// New 创建一个新的日志记录器
func New(opts LogConfig) (*Logger, error) {
	// 解析日志级别
	level, err := parseLogLevel(opts.Level)
	if err != nil {
		return nil, err
	}

	// 创建编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder, // 使用彩色日志级别
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建文件写入器
	fileWriter := &lumberjack.Logger{
		Filename:   opts.Filename,
		MaxSize:    opts.MaxSize,
		MaxBackups: opts.MaxBackups,
		MaxAge:     opts.MaxAge,
		Compress:   opts.Compress,
	}

	// 创建多个输出核心
	cores := []zapcore.Core{
		// 文件输出核心（JSON格式）
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(fileWriter),
			level,
		),
		// 控制台输出核心（带颜色的格式化输出）
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		),
	}

	// 使用NewTee将多个核心组合
	core := zapcore.NewTee(cores...)

	// 创建日志记录器
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{logger: logger}, nil
}

// Debug 记录调试级别的日志
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Info 记录信息级别的日志
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Warn 记录警告级别的日志
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Error 记录错误级别的日志
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal 记录致命级别的日志
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// With 创建一个带有额外字段的日志记录器
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{logger: l.logger.With(fields...)}
}

// Sync 同步日志缓冲区
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

// parseLogLevel 解析日志级别字符串
func parseLogLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, ErrInvalidLogLevel
	}
}
