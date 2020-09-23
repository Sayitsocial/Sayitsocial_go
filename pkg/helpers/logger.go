package helpers

import (
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *logrus.Logger

// LoggerInit initializes logger
func LoggerInit() {
	logger = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "02-01-2006 15:04:05",
			LogFormat:       "[%lvl%] [%component%] %time% - %msg%\n",
		},
	}

	l := &lumberjack.Logger{
		Filename:   LogsPath + "/latest.log",
		MaxSize:    1, // MegaBytes
		MaxBackups: 8, // Max Files
		MaxAge:     7, // Days
		Compress:   true,
	}
	mWriter := io.MultiWriter(os.Stdout, l)
	logger.SetOutput(mWriter)
}

// LogError Log errors to console with formatting
func LogError(message interface{}) {
	logger.WithFields(logrus.Fields{
		"component": func() string {
			pc, _, _, ok := runtime.Caller(2)
			return getCalledInfo(pc, ok)
		}(),
	}).Error(message)
}

// LogInfo Log info to console with formatting
func LogInfo(message interface{}) {
	logger.WithFields(logrus.Fields{
		"component": func() string {
			pc, _, _, ok := runtime.Caller(2)
			return getCalledInfo(pc, ok)
		}(),
	}).Info(message)
}

// LogWarning Log warning to console with formatting
func LogWarning(message interface{}) {
	logger.WithFields(logrus.Fields{
		"component": func() string {
			pc, _, _, ok := runtime.Caller(2)
			return getCalledInfo(pc, ok)
		}(),
	}).Warningln(message)
}

func getCalledInfo(pc uintptr, ok bool) string {
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		split := strings.Split(details.Name(), "/")
		return split[len(split)-1]
	}
	return "Unknown"
}
