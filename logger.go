package artnet

import (
	"os"

	"github.com/Sirupsen/logrus"
)

type Fields map[string]interface{}

type Logger interface {
	logrus.StdLogger
	With(fields Fields) *logger
}

type logger struct {
	*logrus.Entry
}

type logRus struct {
	*logrus.Logger
}

func NewLogger() Logger {
	log := &logRus{
		Logger: logrus.New(),
	}

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05.0000"
	customFormatter.DisableColors = true
	customFormatter.FullTimestamp = true

	log.Formatter = customFormatter

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.Out = os.Stdout

	// Only log the debug severity or above.
	log.Level = logrus.DebugLevel

	// Disable concurrency mutex as we use Stdout
	log.SetNoLock()
	return &logger{Entry: log.WithFields(nil)}
}

/*
func (l *logger) With(fields Fields) *Entry {
	return &Entry{Entry: l.Logger.WithFields(logrus.Fields(fields))}
}
*/

func (l *logger) With(fields Fields) *logger {
	return &logger{Entry: l.WithFields(logrus.Fields(fields))}
}
