package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/kataras/golog"
	"github.com/samber/do"
)

type LoggerService interface {
	SetLevel(levelName string)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	Logf(level golog.Level, format string, arg ...any)
}

func NewLoggerService(injector *do.Injector) (LoggerService, error) {
	s := loggerServiceImpl{
		logger: golog.New(),
	}
	s.SetLevel("DEBUG")
	return &s, nil
}

type loggerServiceImpl struct {
	logger *golog.Logger
}

// 设置日志级别
func (l *loggerServiceImpl) SetLevel(levelName string) {
	l.logger.SetLevel(levelName)
}

func (l *loggerServiceImpl) Infof(format string, args ...interface{}) {
	l.Logf(golog.InfoLevel, format, args...)
}

func (l *loggerServiceImpl) Debugf(format string, args ...interface{}) {
	l.Logf(golog.DebugLevel, format, args...)
}

func (l *loggerServiceImpl) Errorf(format string, args ...interface{}) {
	l.Logf(golog.ErrorLevel, format, args...)

	// send message when log
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelError)
		sentry.CaptureMessage(fmt.Sprintf(format, args...))
	})
}

func (l *loggerServiceImpl) Warnf(format string, args ...interface{}) {
	l.Logf(golog.WarnLevel, format, args...)

	// send message when log
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelWarning)
		sentry.CaptureMessage(fmt.Sprintf(format, args...))
	})
}

func (l *loggerServiceImpl) Fatalf(format string, args ...interface{}) {
	l.Logf(golog.FatalLevel, format, args...)

	// send message when log
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelFatal)
		sentry.CaptureMessage(fmt.Sprintf(format, args...))
	})
}

func (l *loggerServiceImpl) Logf(level golog.Level, format string, arg ...interface{}) {
	shouldSkip := true
	skipTimes := 0

	var (
		file    string
		line    int
		fun     uintptr
		funName string
	)

	funNames := []string{
		"pkg/logger.(*loggerServiceImpl).Logf",
		"pkg/logger.(*loggerServiceImpl).Errorf",
		"pkg/logger.(*loggerServiceImpl).Debugf",
		"pkg/logger.(*loggerServiceImpl).Warnf",
		"pkg/logger.(*loggerServiceImpl).Infof",
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

	l.logger.SetPrefix(fmt.Sprintf("[%s:%d] %s: ", file, line, funName))
	l.logger.Logf(level, format, arg...)
}
