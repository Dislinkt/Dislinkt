package logging

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

type Logger struct {
	InfoLogger  *logrus.Logger
	WarnLogger  *logrus.Logger
	ErrorLogger *logrus.Logger
}

func InitLogger(ctx context.Context) *Logger {
	logger := &Logger{}
	logger.InfoLogger = InitLoggerPerLevel("info")
	logger.WarnLogger = InitLoggerPerLevel("warn")
	logger.ErrorLogger = InitLoggerPerLevel("error")

	logger.InfoLogger.SetLevel(logrus.InfoLevel)
	logger.WarnLogger.SetLevel(logrus.WarnLevel)
	logger.ErrorLogger.SetLevel(logrus.ErrorLevel)

	return &Logger{InfoLogger: logger.InfoLogger, WarnLogger: logger.WarnLogger, ErrorLogger: logger.ErrorLogger}
}

func InitLoggerPerLevel(logFile string) *logrus.Logger {
	path := filepath.Join("logs", time.Now().Format("2006-01-02")+logFile+".log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	logger := logrus.New()
	wrt := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(wrt)
	// logger.SetOutput(&lumberjack.Logger{
	// 	Filename:   file.Name(),
	// 	MaxSize:    100, // megabytes
	// 	MaxBackups: 3,
	// 	MaxAge:     30,   // days
	// 	Compress:   true, // disabled by default
	// })

	logger.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 19:10:10",
		LogFormat:       "[%lvl%] %time% | %msg% \n",
	})

	return logger
}
