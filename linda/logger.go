package linda

import (
	log "github.com/sirupsen/logrus"
	"os"
)

var Logger = log.New()

func init() {
	log.SetLevel(log.DebugLevel)
	Logger.Out = os.Stdout
	Logger.Debug("init Logger")
}
