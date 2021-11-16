package app

import (
	"io"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/bububa/visiondb/logger"
)

func InitLogger() error {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	zerolog.CallerSkipFrameCount = 2

	consoleWriter := zerolog.NewConsoleWriter()
	consoleWriter.TimeFormat = "01/02-15:04:05"
	log.Logger = log.Output(consoleWriter)
	writers := []io.Writer{consoleWriter}
	if config.LogPath != "" {
		fileWriter, err := logger.File(config.LogPath, logger.FileConfig{
			MaxSize:    100,
			MaxAge:     3,
			MaxBackups: 3,
			LocalTime:  true,
			Compress:   true,
		})
		if err != nil {
			return err
		}
		writers = append(writers, fileWriter)
		multiWriter := zerolog.MultiLevelWriter(writers...)
		log.Logger = zerolog.New(multiWriter).With().Timestamp().Logger()
	}
	logger.SetLogger(log.Logger)
	return nil
}
