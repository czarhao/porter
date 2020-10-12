package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
	Logger.Formatter = new(logrus.TextFormatter)
	Logger.Formatter.(*logrus.TextFormatter).DisableTimestamp = true
	Logger.Level = logrus.InfoLevel
	Logger.Out = os.Stdout
}
