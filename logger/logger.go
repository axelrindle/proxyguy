package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func make() *logrus.Logger {
	logger := logrus.New()

	logger.SetOutput(os.Stderr)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableSorting:         true,
		DisableLevelTruncation: true,
		SortingFunc:            func(s []string) {},
	})

	return logger
}

var Logger = make()

func WithField(key string, value interface{}) *logrus.Entry {
	return Logger.WithField(key, value)
}

func ForModule(name string) *logrus.Entry {
	return WithField("module", name)
}

func ForComponent(name string) *logrus.Entry {
	return WithField("component", name)
}
