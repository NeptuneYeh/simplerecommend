package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var MyLog *zerolog.Logger

type Module struct {
	ZeroLogger zerolog.Logger
}

func NewModule() *Module {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// global pretty logging 效果
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// 如果只想要個別創建的實例有 pretty logging 效果
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: zerolog.TimeFormatUnix}
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	MyLog = &logger

	return &Module{
		ZeroLogger: logger,
	}
}
