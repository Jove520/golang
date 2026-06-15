package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var sugarLogger *zap.SugaredLogger

func GetLogger() *zap.SugaredLogger {
	if sugarLogger != nil {
		return sugarLogger
	}

	level := zap.NewAtomicLevelAt(getLogLevel(os.Getenv("GEEKAI_LOG_LEVEL")))
	encoder := getEncoder()

	fileCore := zapcore.NewCore(encoder, getLogWriter(), level)
	consoleCore := zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), level)

	core := zapcore.NewTee(fileCore, consoleCore)

	zapLogger := zap.New(core, zap.AddCaller())
	sugarLogger = zapLogger.Sugar()

	return sugarLogger
}

func getEncoder() zapcore.Encoder {
	config := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",

		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	return zapcore.NewConsoleEncoder(config)
}

func getLogWriter() zapcore.WriteSyncer {
	_ = os.MkdirAll("logs", 0755)

	writer := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}

	return zapcore.AddSync(writer)
}

func getLogLevel(level string) zapcore.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return zapcore.DebugLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
