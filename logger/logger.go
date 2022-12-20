package logger

import (
	"io"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	Debug         bool
	ConsoleOutput string
	FileOutput    string
	FileConfig    lumberjack.Logger
}

func (l *Logger) StartLoggers() *zerolog.Logger {

	var log zerolog.Logger
	var writers []io.Writer

	// add file writer if required
	if file, _ := strconv.ParseBool(l.FileOutput); file {
		writers = append(writers, &l.FileConfig)
	}
	// add console writer if required
	if console, _ := strconv.ParseBool(l.ConsoleOutput); console {
		cw := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
		writers = append(writers, cw)
	}

	mw := zerolog.MultiLevelWriter(writers...)
	log = zerolog.New(mw).With().Timestamp().Logger()

	// add debug info when debug mode is on
	if l.Debug {
		log = log.Level(zerolog.DebugLevel).With().Caller().Logger()
	}

	return &log
}
