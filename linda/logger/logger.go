package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var log = logrus.New()

func init() {
	log.Out = os.Stdout
	log.Info("logger init")
}
