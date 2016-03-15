package log

import (
	"io/ioutil"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/rifflock/lfshook"
	"github.com/spf13/viper"
)

var termLogger = logrus.New()
var fileLogger = logrus.New()

func InitLog() {
	// Initialize logger
	termLogger.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	fileLogger.Formatter = &logrus.JSONFormatter{}
	fileLogger.Out = ioutil.Discard

	logLevel := viper.GetString("log.level")
	logFile := viper.GetString("log.file")

	strings.ToLower(logLevel)
	switch logLevel {
	case "panic":
		termLogger.Level = logrus.PanicLevel
		fileLogger.Level = logrus.PanicLevel
	case "fatal":
		termLogger.Level = logrus.FatalLevel
		fileLogger.Level = logrus.FatalLevel
	case "error":
		termLogger.Level = logrus.ErrorLevel
		fileLogger.Level = logrus.ErrorLevel
	case "warn":
		termLogger.Level = logrus.WarnLevel
		fileLogger.Level = logrus.WarnLevel
	case "info":
		termLogger.Level = logrus.InfoLevel
		fileLogger.Level = logrus.InfoLevel
	case "debug":
		termLogger.Level = logrus.DebugLevel
		fileLogger.Level = logrus.DebugLevel
	default:
		termLogger.Level = logrus.ErrorLevel
		fileLogger.Level = logrus.ErrorLevel
	}

	fileLogger.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
		logrus.PanicLevel: logFile,
		logrus.FatalLevel: logFile,
		logrus.ErrorLevel: logFile,
		logrus.WarnLevel:  logFile,
		logrus.InfoLevel:  logFile,
		logrus.DebugLevel: logFile,
	}))
}

// wrapper for panic
func Panic(fields map[string]interface{}, msg string) {

	if fields == nil {
		fileLogger.Panic(msg)
		if logrus.IsTerminal() {
			termLogger.Panic(msg)
		}
	} else {
		fileLogger.WithFields(fields).Panic(msg)
		if logrus.IsTerminal() {
			termLogger.WithFields(fields).Panic(msg)
		}
	}
}

// wrapper for fatal
func Fatal(fields map[string]interface{}, msg string) {

	if fields == nil {
		fileLogger.Fatal(msg)
		if logrus.IsTerminal() {
			termLogger.Fatal(msg)
		}
	} else {
		fileLogger.WithFields(fields).Fatal(msg)
		if logrus.IsTerminal() {
			termLogger.WithFields(fields).Fatal(msg)
		}
	}
}

// wrapper for error
func Error(fields map[string]interface{}, msg string) {

	if fields == nil {
		fileLogger.Error(msg)
		if logrus.IsTerminal() {
			termLogger.Error(msg)
		}
	} else {
		fileLogger.WithFields(fields).Error(msg)
		if logrus.IsTerminal() {
			termLogger.WithFields(fields).Error(msg)
		}
	}
}

// wrapper for warn
func Warn(fields map[string]interface{}, msg string) {

	if fields == nil {
		fileLogger.Warn(msg)
		if logrus.IsTerminal() {
			termLogger.Warn(msg)
		}
	} else {
		fileLogger.WithFields(fields).Warn(msg)
		if logrus.IsTerminal() {
			termLogger.WithFields(fields).Warn(msg)
		}
	}
}

// wrapper for info
func Info(fields map[string]interface{}, msg string) {

	if fields == nil {
		fileLogger.Info(msg)
		if logrus.IsTerminal() {
			termLogger.Info(msg)
		}
	} else {
		fileLogger.WithFields(fields).Info(msg)
		if logrus.IsTerminal() {
			termLogger.WithFields(fields).Info(msg)
		}
	}
}

// wrapper for debug
func Debug(fields map[string]interface{}, msg string) {

	if fields == nil {
		fileLogger.Debug(msg)
		if logrus.IsTerminal() {
			termLogger.Debug(msg)
		}
	} else {
		fileLogger.WithFields(fields).Debug(msg)
		if logrus.IsTerminal() {
			termLogger.WithFields(fields).Debug(msg)
		}
	}
}
