package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

type tLogger struct {
	logger      *log.Logger
	level       LogLevel
	callerLevel int
	buf         chan string
	close       chan struct{}
	mu          sync.Mutex
}

// New 创建一个新的日志记录器
func New(out io.Writer, prefix string, flag int) Logger {
	l := new(tLogger)
	l.logger = log.New(out, prefix, flag)
	l.level = LOG_DEBUG
	l.callerLevel = 2

	// buffer flush
	l.buf = make(chan string, 128)
	l.close = make(chan struct{})
	go func(l *tLogger) {
		for s := range l.buf {
			l.logger.Print(s)
		}
		l.close <- struct{}{}
	}(l)

	return l
}

func (t *tLogger) Level() LogLevel {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.level
}

func (t *tLogger) SetLevel(level LogLevel) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.level = level
}

func (t *tLogger) LevelName(level LogLevel) string {
	if n, ok := LevelNames[level]; ok {
		return n
	}
	return "^?^"
}

func (t *tLogger) SetCallerLevel(v int) {
	t.callerLevel = v
}

func (t *tLogger) Flush() {
	close(t.buf)
	<-t.close
}

func (t *tLogger) Log(v ...interface{}) {
	t.logger.Print(v...)
}

func (t *tLogger) Logf(format string, v ...interface{}) {
	t.logger.Printf(format, v...)
}

func (t *tLogger) Logln(v ...interface{}) {
	t.logger.Println(v...)
}

func (t *tLogger) Print(v ...interface{}) {
	t.logger.Print(v...)
}

func (t *tLogger) Printf(format string, v ...interface{}) {
	t.logger.Printf(format, v...)
}

func (t *tLogger) Println(v ...interface{}) {
	t.logger.Println(v...)
}

func (t *tLogger) Fatal(v ...interface{}) {
	t.output(LOG_FATAL, "", false, v...)
	os.Exit(1)
}

func (t *tLogger) Fatalf(format string, v ...interface{}) {
	t.output(LOG_FATAL, format, false, v...)
	os.Exit(1)
}

func (t *tLogger) Fatalln(v ...interface{}) {
	t.output(LOG_FATAL, "", true, v...)
	os.Exit(1)
}

func (t *tLogger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	t.output(LOG_PANIC, "", false, s)
	panic(s)
}

func (t *tLogger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	t.output(LOG_PANIC, "", false, s)
	panic(s)
}

func (t *tLogger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	t.output(LOG_PANIC, "", false, s)
	panic(s)
}

func (t *tLogger) Error(v ...interface{}) {
	t.output(LOG_ERROR, "", false, v...)
}

func (t *tLogger) Errorf(format string, v ...interface{}) {
	t.output(LOG_ERROR, format, false, v...)
}

func (t *tLogger) Errorln(v ...interface{}) {
	t.output(LOG_ERROR, "", true, v...)
}

func (t *tLogger) Warn(v ...interface{}) {
	t.output(LOG_WARN, "", false, v...)
}

func (t *tLogger) Warnf(format string, v ...interface{}) {
	t.output(LOG_WARN, format, false, v...)
}

func (t *tLogger) Warnln(v ...interface{}) {
	t.output(LOG_WARN, "", true, v...)
}

func (t *tLogger) Info(v ...interface{}) {
	t.output(LOG_INFO, "", false, v...)
}

func (t *tLogger) Infof(format string, v ...interface{}) {
	t.output(LOG_INFO, format, false, v...)
}

func (t *tLogger) Infoln(v ...interface{}) {
	t.output(LOG_INFO, "", true, v...)
}

func (t *tLogger) Debug(v ...interface{}) {
	t.output(LOG_DEBUG, "", false, v...)
}

func (t *tLogger) Debugf(format string, v ...interface{}) {
	t.output(LOG_DEBUG, format, false, v...)
}

func (t *tLogger) Debugln(v ...interface{}) {
	t.output(LOG_DEBUG, "", true, v...)
}

func (t *tLogger) output(level LogLevel, format string, newline bool, v ...interface{}) {
	if t.level == LOG_OFF {
		return
	}
	if t.level < level {
		return
	}

	// 记录调用者位置
	_, file, line, ok := runtime.Caller(t.callerLevel)
	if !ok {
		file = "???"
		line = 0
	}

	// short file
	// pos := strings.LastIndex(file, "/") + 1
	// file = string(file[pos:])

	s := fmt.Sprintf("%s:%d [%s] ", file, line, t.LevelName(level))

	var formatText string
	if format == "" {
		formatText = fmt.Sprint(v...)
	} else {
		formatText = fmt.Sprintf(format, v...)
	}

	s += formatText

	if newline {
		s += "\n"
	}

	// 错误上报
	if level > LOG_FATAL && level < LOG_WARN {
		// 上报
	}

	// PANIC, FATAL immediately output
	if level <= LOG_PANIC {
		t.logger.Print(s)
	} else {
		// in buffer
		t.buf <- s
	}
}
