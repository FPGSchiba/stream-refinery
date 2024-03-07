package util

import (
	"fmt"
	"log"
	"os"
)

const (
	LevelDebug = iota
	LevelInfo  = iota
	LevelError = iota
)

const (
	LogTypeFile    = iota
	LogTypeConsole = iota
	LogTypeConFile = iota
)

type Logger struct {
	LogLevel     int
	FilePath     string
	LogType      int
	defaultLevel int `default:"1"`
}

func (p Logger) logToFile(message string) {
	panic("File logging not implemented!")
}

func (p Logger) logToConsole(message string) {
	println(message)
}

func (p Logger) logToErrConsole(message string) {
	l := log.New(os.Stderr, "", 0)
	l.Println(message)
}

func (p Logger) logOut(message string, level string) {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(LoggerErrorCode)
	}
	switch p.LogType {
	case LogTypeFile:
		p.logToFile(fmt.Sprintf("%s - %s [%s]: %s", hostname, "test", level, message))
	case LogTypeConsole:
		p.logToConsole(fmt.Sprintf("%s - %s [%s]: %s", hostname, "test", level, message))
	case LogTypeConFile:
		p.logToConsole(fmt.Sprintf("%s - %s [%s]: %s", hostname, "test", level, message))
		p.logToFile(fmt.Sprintf("%s - %s [%s]: %s", hostname, "test", level, message))
	default:
		panic("unhandled default case")
	}
}

func (p Logger) logError(message string) {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(LoggerErrorCode)
	}
	switch p.LogType {
	case LogTypeFile:
		p.logToFile(fmt.Sprintf("%s - %s [ERROR]: %s", hostname, "test", message))
	case LogTypeConsole:
		p.logToErrConsole(fmt.Sprintf("%s - %s [ERROR]: %s", hostname, "test", message))
	case LogTypeConFile:
		p.logToErrConsole(fmt.Sprintf("%s - %s [ERROR]: %s", hostname, "test", message))
		p.logToFile(fmt.Sprintf("%s - %s [ERROR]: %s", hostname, "test", message))
	default:
		panic("unhandled default case")
	}
}

func (p Logger) Log(message string, level int) {
	if message != "" {
		switch level {
		case LevelDebug:
			p.logOut(message, "DEBUG")
		case LevelInfo:
			p.logOut(message, "INFO")
		case LevelError:
			p.logError(message)
		default:
			panic("unhandled default case")
		}
	} else {
		return
	}
}
