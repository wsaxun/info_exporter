package tools

import (
	"github.com/sirupsen/logrus"
)

// LogsInit 设置格式
func LogsInit() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func logWithFields(level logrus.Level, info logrus.Fields, args ...any) {
	entry := logrus.WithFields(info)
	switch level {
	case logrus.InfoLevel:
		entry.Infoln(args...)
	case logrus.ErrorLevel:
		entry.Errorln(args...)
	case logrus.FatalLevel:
		entry.Fatalln(args...)
	default:
		entry.Warn("Unrecognized logging level")
	}
}

func LogInfo(info logrus.Fields, args ...any) {
	logWithFields(logrus.InfoLevel, info, args...)
}

func LogErr(info logrus.Fields, args ...any) {
	logWithFields(logrus.ErrorLevel, info, args...)
}

func LogFatal(info logrus.Fields, args ...any) {
	logWithFields(logrus.FatalLevel, info, args...)
}
