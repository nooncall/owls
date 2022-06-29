package logger

import (
	"os"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zlog *zap.SugaredLogger

const loggerTimeLayout = "2006-01-02 15:04:05.000"

func InitLog(logDir, fileName string, level string) {
	logFile := logDir + "/" + fileName

	//1024, 0, 0, true, false
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    1024,
		MaxBackups: 0,
		MaxAge:     0,
		LocalTime:  true,
		Compress:   false,
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	// 格式化时间
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(loggerTimeLayout))
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	writeSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	core := zapcore.NewCore(encoder, writeSyncer, ConvertLevel(level))

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zlog = logger.Sugar()
}

func Info(args ...interface{}) {
	zlog.Info(args...)
}

func Infof(template string, args ...interface{}) {
	zlog.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	zlog.Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	zlog.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	zlog.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	zlog.Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	zlog.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	zlog.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	zlog.Errorw(msg, keysAndValues...)
}

// ConvertLevel convert string level to Level
func ConvertLevel(level string) zapcore.Level {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	case "panic":
		return zap.PanicLevel
	default:
		return zap.InfoLevel
	}
}
