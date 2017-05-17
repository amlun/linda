package linda

import (
	log "github.com/sirupsen/logrus"
	"os"
)

var logger = log.New()

func init() {
	log.SetLevel(log.DebugLevel)
	logger.Out = os.Stdout
	logger.Debug("init logger")
}
