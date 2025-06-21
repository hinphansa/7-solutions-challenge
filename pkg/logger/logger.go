package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	logrus.FieldLogger
}

func New(level logrus.Level) Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(level)

	return log
}
