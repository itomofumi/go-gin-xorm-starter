package util

import (
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var globalLogger *logrus.Logger

// GetLogger returns global logger
func GetLogger() *logrus.Logger {

	if globalLogger != nil {
		return globalLogger
	}

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{}
	// Output to stdout instead of the default stderr
	// log.SetOutput(os.Stdout)
	rollingLogger := &lumberjack.Logger{
		Filename:   path.Join(os.Getenv("LOG_DIR"), "debug.log"),
		MaxSize:    10,    // megabytes
		MaxBackups: 100,   // default: not to remove old logs
		Compress:   false, // disabled by default
	}
	logger.Out = io.MultiWriter(os.Stdout, rollingLogger)

	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logger.Warningln(err)
		level = logrus.DebugLevel
	}
	logger.Level = level
	logger.Infof("Global Log Level is %v", level.String())

	globalLogger = logger

	return logger
}
