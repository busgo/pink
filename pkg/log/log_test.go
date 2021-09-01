package log

import (
	"github.com/labstack/gommon/log"
	"testing"
)

func TestNewLoggerSugar(t *testing.T) {

	NewLoggerSugar("log_test.log", "1.log", -1)

	log.Debugf("admin:%s", "-----admin")
	log.Infof("admin:%s", "-----admin")
	log.Warnf("admin:%s", "-----admin")
	log.Errorf("admin:%s", "-----admin")

}
