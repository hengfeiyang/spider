package log

import (
	"log"
	"os"
)

const (
	LOGOPEN  = true
	LOGFILE  = "os.Stderr"
	LOGLEVEL = LOG_DEBUG
)

// DefaultLogger 日志接口
var DefaultLogger Logger

func Level() LogLevel {
	return DefaultLogger.Level()
}

func SetLevel(l LogLevel) {
	DefaultLogger.SetLevel(l)
}

func Flush() {
	DefaultLogger.Flush()
}

func Print(v ...interface{}) {
	DefaultLogger.Print(v...)
}

func Printf(format string, v ...interface{}) {
	DefaultLogger.Printf(format, v...)
}

func Println(v ...interface{}) {
	DefaultLogger.Println(v...)
}

func Fatal(v ...interface{}) {
	DefaultLogger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	DefaultLogger.Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	DefaultLogger.Fatalln(v...)
}

func Panic(v ...interface{}) {
	DefaultLogger.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	DefaultLogger.Panicf(format, v...)
}

func Panicln(v ...interface{}) {
	DefaultLogger.Panicln(v...)
}

func Error(v ...interface{}) {
	DefaultLogger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	DefaultLogger.Errorf(format, v...)
}

func Errorln(v ...interface{}) {
	DefaultLogger.Errorln(v...)
}

func Warn(v ...interface{}) {
	DefaultLogger.Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	DefaultLogger.Warnf(format, v...)
}

func Warnln(v ...interface{}) {
	DefaultLogger.Warnln(v...)
}

func Info(v ...interface{}) {
	DefaultLogger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	DefaultLogger.Infof(format, v...)
}

func Infoln(v ...interface{}) {
	DefaultLogger.Infoln(v...)
}

func Debug(v ...interface{}) {
	DefaultLogger.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	DefaultLogger.Debugf(format, v...)
}

func Debugln(v ...interface{}) {
	DefaultLogger.Debugln(v...)
}

func init() {
	var f *os.File
	var err error
	if LOGOPEN == false {
		if f, err = os.Open(os.DevNull); err != nil {
			panic("open log file error: " + err.Error())
		}
	} else {
		if LOGFILE == "" || LOGFILE == "os.Stderr" {
			f = os.Stderr
		} else if LOGFILE == "os.Stdout" {
			f = os.Stdout
		} else {
			if f, err = os.OpenFile(LOGFILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
				panic("open log file error: " + err.Error())
			}
		}
	}
	DefaultLogger = New(f, "[spider]", log.LstdFlags)
	DefaultLogger.SetLevel(LogLevel(LOGLEVEL))
	DefaultLogger.(*tLogger).SetCallerLevel(3)
}
