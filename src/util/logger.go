package util

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	LevelDebug = 0
	LevelInfo  = 2
	LevelError = 5
)

const (
	LogTypeFile    = "file"
	LogTypeConsole = "console"
	LogTypeConFile = "confile"
)

type Logger struct {
	LogLevel int
	FilePath string
	LogType  string
}

func (p Logger) logToFile(message string) {
	f, err := os.OpenFile(p.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		p.LogType = LogTypeConsole
		p.Log(fmt.Sprintf("Unhandled Logging Error: %s", err.Error()), LevelError)
		p.logToConsole(message)
		return
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			p.LogType = LogTypeConsole
			p.Log(fmt.Sprintf("Unhandled Logging Error: %s", err.Error()), LevelError)
			p.logToConsole(message)
			return
		}
	}(f)

	if _, err = f.WriteString(message + "\n"); err != nil {
		p.LogType = LogTypeConsole
		p.Log(fmt.Sprintf("Unhandled Logging Error: %s", err.Error()), LevelError)
		p.logToConsole(message)
		return
	}
}

func (p Logger) getTime() string {
	dt := time.Now()
	return dt.Format("02.01.2006 15:04:05")
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
		p.logToFile(fmt.Sprintf("%s - %s [%s]: %s", hostname, p.getTime(), level, message))
	case LogTypeConsole:
		p.logToConsole(fmt.Sprintf("%s - %s [%s]: %s", hostname, p.getTime(), level, message))
	case LogTypeConFile:
		p.logToConsole(fmt.Sprintf("%s - %s [%s]: %s", hostname, p.getTime(), level, message))
		p.logToFile(fmt.Sprintf("%s - %s [%s]: %s", hostname, p.getTime(), level, message))
	default:
		p.logToConsole(fmt.Sprintf("%s - %s [%s]: %s", hostname, p.getTime(), level, message))
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
		p.logToFile(fmt.Sprintf("%s - %s [ERROR]: %s", hostname, p.getTime(), message))
	case LogTypeConsole:
		p.logToErrConsole(fmt.Sprintf("%s - %s [ERROR]: %s", hostname, p.getTime(), message))
	case LogTypeConFile:
		p.logToErrConsole(fmt.Sprintf("%s - %s [ERROR]: %s", hostname, p.getTime(), message))
		p.logToFile(fmt.Sprintf("%s - %s [ERROR]: %s", hostname, p.getTime(), message))
	default:
		p.logToErrConsole(fmt.Sprintf("%s - %s [ERROR]: %s", hostname, p.getTime(), message))
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
