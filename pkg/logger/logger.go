package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/kataras/golog"
)

var logger *golog.Logger

func Init() {
	logger = golog.New()
}

// 设置日志级别
func SetLevel(levelName string) {
	logger.SetLevel(levelName)
}

func Infof(format string, args ...interface{}) {
	Logf(golog.InfoLevel, format, args...)
}

func Debugf(format string, args ...interface{}) {
	Logf(golog.DebugLevel, format, args...)
}

func Errorf(format string, args ...interface{}) {
	Logf(golog.ErrorLevel, format, args...)

	// send message when log
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelError)
		sentry.CaptureMessage(fmt.Sprintf(format, args...))
	})
}

func Warnf(format string, args ...interface{}) {
	Logf(golog.WarnLevel, format, args...)

	// send message when log
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelWarning)
		sentry.CaptureMessage(fmt.Sprintf(format, args...))
	})
}

func Fatalf(format string, args ...interface{}) {
	Logf(golog.FatalLevel, format, args...)

	// send message when log
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelFatal)
		sentry.CaptureMessage(fmt.Sprintf(format, args...))
	})
}

func Logf(level golog.Level, format string, arg ...interface{}) {
	shouldSkip := true
	skipTimes := 0

	var (
		file    string
		line    int
		fun     uintptr
		funName string
	)

	funNames := []string{
		"pkg/logger.Logf",
		"pkg/logger.Errorf",
		"pkg/logger.Debugf",
		"pkg/logger.Warnf",
		"pkg/logger.Infof",
	}

	for shouldSkip {
		fun, file, line, _ = runtime.Caller(skipTimes)
		funName = runtime.FuncForPC(fun).Name()

		find := false
		for _, name := range funNames {
			if strings.Contains(funName, name) {
				find = true
				break
			}
		}
		if find {
			skipTimes++
		} else {
			shouldSkip = false
		}
	}

	funNamePart := strings.Split(funName, ".")
	funName = funNamePart[len(funNamePart)-1]

	logger.SetPrefix(fmt.Sprintf("[%s:%d] %s: ", file, line, funName))
	logger.Logf(level, format, arg...)
}
