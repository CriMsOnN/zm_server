package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	loggers = make(map[string]*logrus.Logger)
	mu      sync.Mutex
)

func NewComponentLogger(component string) *logrus.Entry {
	mu.Lock()
	defer mu.Unlock()
	if log, exists := loggers[component]; exists {
		return log.WithField("component", component)
	}
	log := logrus.New()
	if err := os.MkdirAll("logs", 0o755); err != nil {
		log.Errorf("failed to create logs directory: %v", err)
		return nil
	}

	fileWriter := &lumberjack.Logger{
		Filename:   filepath.Join("logs", fmt.Sprintf("%s_%s.log", component, time.Now().Format("2006-01-02"))),
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	mw := io.MultiWriter(os.Stdout, fileWriter)
	log.SetOutput(mw)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	loggers[component] = log
	return log.WithField("component", component)
}
