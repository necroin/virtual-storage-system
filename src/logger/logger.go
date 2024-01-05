package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
)

const (
	errorLevel = iota
	infoLevel
	verboseLevel
	debugLevel
)

var (
	enable      = false
	logLevel    = infoLevel
	mutex       sync.Mutex
	logLevelMap = map[string]int{
		"error":   errorLevel,
		"info":    infoLevel,
		"verbose": verboseLevel,
		"debug":   debugLevel,
	}
)

func Configure(logsEnable bool, logsPath string, logsLevel string) error {
	enable = logsEnable
	if enable {
		log.SetOutput(os.Stdout)
		if logsPath != "" {
			if err := os.MkdirAll(path.Dir(logsPath), os.ModePerm); err != nil {
				return fmt.Errorf("[Logger] [Error] failed create logs directory: %s", err)
			}

			var logsFile *os.File
			logsFile, err := os.OpenFile(logsPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
			if err != nil {
				return fmt.Errorf("[Logger] [Error] failed to open/create logs file: %s", err)
			}
			log.SetOutput(logsFile)
		}

		configLogLevel, ok := logLevelMap[logsLevel]
		if ok {
			logLevel = configLogLevel
		}
	}
	return nil
}

func print(message string) {
	if enable {
		mutex.Lock()
		defer mutex.Unlock()
		log.Println(message)
	}
}

func Error(message string, args ...any) {
	if enable && logLevel >= errorLevel {
		print("ERROR: " + fmt.Sprintf(message, args...))
	}
}

func Info(message string, args ...any) {
	if logLevel >= infoLevel {
		print("INFO: " + fmt.Sprintf(message, args...))
	}
}

func Verbose(message string, args ...any) {
	if logLevel >= verboseLevel {
		print("VERBOSE: " + fmt.Sprintf(message, args...))
	}
}

func Debug(message string, args ...any) {
	if logLevel >= debugLevel {
		print("DEBUG: " + fmt.Sprintf(message, args...))
	}
}
