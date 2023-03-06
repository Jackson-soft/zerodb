package zerolog

import (
	"errors"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	//默认日志保存时间
	defaultMaxLogDay = 7
	//时间格式
	defaultLogTimeFormat = "2006-01-02 15:04:05.000"
)

const (
	JSON    = "json"
	CONSOLE = "console"
)

//Zlog 是日志的封装
type Zlog struct {
	logger *zap.SugaredLogger
}

/* NewLog 创建日志
format 日志格式化方式 json, console
logFile 日志文件路径,至少要是目录加上文件名，例如： zlog/log
lvl 日志的级别包含 debug info warning error fatal panic
module 日志模块（前缀）
maxAge 日志文件最长保存时间
*/
func NewLog(format, logFile, lvl string, maxAge int) (*Zlog, error) {
	if len(logFile) == 0 || len(lvl) == 0 {
		return nil, errors.New("log file is nil")
	}

	if maxAge == 0 {
		maxAge = defaultMaxLogDay
	}

	hook := lumberjack.Logger{
		Filename:   logFile, // 日志文件路径
		MaxSize:    1024,    // megabytes
		MaxAge:     maxAge,  //days
		MaxBackups: 3,       // 最多保留3个备份
		LocalTime:  true,
		Compress:   true, // 是否压缩 disabled by default
	}

	var level zapcore.Level
	switch lvl {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "line",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(defaultLogTimeFormat))
		},
		EncodeDuration: zapcore.StringDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,    // 短路径编码器
	}

	var core zapcore.Core
	switch format {
	case "json":
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(&hook),
			zap.NewAtomicLevelAt(level),
		)
	default:
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			os.Stdout,
			zap.NewAtomicLevelAt(level),
		)
	}

	return &Zlog{
		logger: zap.New(core, zap.AddCaller(), zap.Development(), zap.AddCallerSkip(1)).Sugar(),
	}, nil
}

func (z *Zlog) Sync() error {
	return z.logger.Sync()
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func (z *Zlog) WithField(key string, value interface{}) *Zlog {
	return &Zlog{
		logger: z.logger.With(key, value),
	}
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.

func (z *Zlog) WithFields(args ...interface{}) *Zlog {
	return &Zlog{
		logger: z.logger.With(args...),
	}
}

// Infof logs a message at level Info on the standard logger.
func (z *Zlog) Infof(format string, args ...interface{}) {
	z.logger.Infof(format, args)
}

// Warnf logs a message at level Warn on the standard logger.
func (z *Zlog) Warnf(format string, args ...interface{}) {
	z.logger.Warnf(format, args)
}

// Errorf logs a message at level Error on the standard logger.
func (z *Zlog) Errorf(format string, args ...interface{}) {
	z.logger.Errorf(format, args)
}

// Panicf logs a message at level Panic on the standard logger.
func (z *Zlog) Panicf(format string, args ...interface{}) {
	z.logger.Panicf(format, args)
}

// Fatalf logs a message at level Fatal on the standard logger.
func (z *Zlog) Fatalf(format string, args ...interface{}) {
	z.logger.Fatalf(format, args)
}

// Infoln logs a message at level Info on the standard logger.
func (z *Zlog) Infoln(args ...interface{}) {
	z.logger.Info(args)
}

// Warnln logs a message at level Warn on the standard logger.
func (z *Zlog) Warnln(args ...interface{}) {
	z.logger.Warn(args)
}

// Errorln logs a message at level Error on the standard logger.
func (z *Zlog) Errorln(args ...interface{}) {
	z.logger.Error(args)
}

// Panicln logs a message at level Panic on the standard logger.
func (z *Zlog) Panicln(args ...interface{}) {
	z.logger.Panic(args)
}

// Fatalln logs a message at level Fatal on the standard logger.
func (z *Zlog) Fatalln(args ...interface{}) {
	z.logger.Fatal(args)
}
