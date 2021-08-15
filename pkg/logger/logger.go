package logger

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/mattn/go-isatty"

	"myst/pkg/config"

	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
)

const (
	DefaultColor = Color(0)

	BlackFg   = Color(aurora.BlackFg)
	RedFg     = Color(aurora.RedFg)
	GreenFg   = Color(aurora.GreenFg)
	YellowFg  = Color(aurora.YellowFg)
	BlueFg    = Color(aurora.BlueFg)
	MagentaFg = Color(aurora.MagentaFg)
	CyanFg    = Color(aurora.CyanFg)
	WhiteFg   = Color(aurora.WhiteFg)

	BlackBg   = Color(aurora.BlackBg)
	RedBg     = Color(aurora.RedBg)
	GreenBg   = Color(aurora.GreenBg)
	YellowBg  = Color(aurora.YellowBg)
	BlueBg    = Color(aurora.BlueBg)
	MagentaBg = Color(aurora.MagentaBg)
	CyanBg    = Color(aurora.CyanBg)
	WhiteBg   = Color(aurora.WhiteBg)
)

var (
	StdoutWriter   *Writer
	StderrWriter   *Writer
	DebugLogWriter *Writer
	ErrorLogWriter *Writer

	defaultLogger *Logger
)

type Color int

type Logger struct {
	name     string
	stdout   *log.Logger
	stderr   *log.Logger
	debugLog *log.Logger
	errorLog *log.Logger
}

func init() {
	// create default stdout and stderr loggers before logger init
	StdoutWriter = NewWriter(colorable.NewColorable(os.Stdout))
	StderrWriter = NewWriter(colorable.NewColorable(os.Stderr))
	// initialize default log writers/loggers
	DebugLogWriter = NewWriter(ioutil.Discard)
	ErrorLogWriter = NewWriter(ioutil.Discard)
	// create default logger
	defaultLogger = New("SERVER", DefaultColor)
}

func Colorize(s string, c Color) string {
	if !(isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())) {
		return s
	}
	return aurora.Colorize(s, aurora.Color(c)).String()
}

func New(name string, color Color) *Logger {
	name = strings.ToUpper(name)
	name = fmt.Sprintf("[%s] ", name)
	name = Colorize(name, color)
	return &Logger{
		name:     name,
		stdout:   log.New(StdoutWriter, "", 0),
		stderr:   log.New(StderrWriter, "", 0),
		debugLog: log.New(DebugLogWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds),
		errorLog: log.New(ErrorLogWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds),
	}
}

//func WithDebugLog(path string) func(*Logger) error {
//	return func(l *Logger) error {
//
//		f, err := os.OpenFile(
//			"logs/debug.log",
//			os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666,
//		)
//		if err != nil {
//			return err
//		}
//
//		DebugLogWriter.SetWriter(colorable.NewNonColorable(f))]
//
//		return nil
//	}
//}

var cwd *string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	wd = filepath.ToSlash(wd + "/")
	cwd = &wd
}

func (l *Logger) prefix(caller bool) string {
	prefix := l.name
	if caller {
		_, file, line, ok := runtime.Caller(3)
		if !ok {
			file = "???"
			line = 0
		}
		if cwd != nil {
			file = strings.Replace(file, *cwd, "", 1)
		}
		prefix = fmt.Sprintf("%s%s:%d ", prefix, file, line)
	}
	return prefix
}

func (l *Logger) print(s string) {
	s = strings.TrimRight(s, "\n")
	_ = l.stdout.Output(3, l.prefix(false)+s)
}

func (l *Logger) debugPrint(s string) {
	s = strings.TrimRight(s, "\n")
	_ = l.debugLog.Output(3, l.prefix(false)+s)
	if config.Debug {
		_ = l.stdout.Output(3, l.prefix(false)+s)
	}
}

func (l *Logger) errorPrint(s string) {
	s = strings.TrimRight(s, "\n")
	s = Colorize(s, RedFg)
	_ = l.errorLog.Output(3, l.prefix(true)+s)
	if config.Debug {
		_ = l.stderr.Output(3, l.prefix(true)+s)
	}
}

func (l *Logger) tracePrint() {
	stack := debug.Stack()
	s := Colorize(string(stack), RedFg)
	_ = l.errorLog.Output(3, l.prefix(true)+s)
	if config.Debug {
		_ = l.stderr.Output(3, l.prefix(true)+s)
	}
}

func Init() error {
	// create logs folder if it doesn't already exist
	_, err := os.Stat("logs")
	if os.IsNotExist(err) {
		err = os.Mkdir("logs", os.ModePerm)
		if err != nil {
			_ = defaultLogger.stderr.Output(0, err.Error())
			return err
		}
	}
	// get debug and error log file
	debugLog, err := os.OpenFile(
		"logs/debug.log",
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666,
	)
	if err != nil {
		_ = defaultLogger.stderr.Output(0, err.Error())
		return err
	}
	errorLog, err := os.OpenFile(
		"logs/error.log",
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666,
	)
	if err != nil {
		_ = defaultLogger.stderr.Output(0, err.Error())
		return err
	}
	// update the writers and loggers
	DebugLogWriter.SetWriter(colorable.NewNonColorable(debugLog))
	ErrorLogWriter.SetWriter(colorable.NewNonColorable(errorLog))
	return nil
}

func Close() {
	err := DebugLogWriter.Sync()
	if err != nil {
		Error(err)
		return
	}
	err = DebugLogWriter.Close()
	if err != nil {
		Error(err)
		return
	}
	err = ErrorLogWriter.Sync()
	if err != nil {
		Error(err)
		return
	}
	err = ErrorLogWriter.Close()
	if err != nil {
		Error(err)
		return
	}
}

func (l *Logger) Debug(v ...interface{}) {
	l.print(fmt.Sprint(v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.print(fmt.Sprintf(format, v...))
}

func (l *Logger) Debugln(v ...interface{}) {
	l.print(fmt.Sprintln(v...))
}

func (l *Logger) Print(v ...interface{}) {
	l.debugPrint(fmt.Sprint(v...))
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.debugPrint(fmt.Sprintf(format, v...))
}

func (l *Logger) Println(v ...interface{}) {
	l.debugPrint(fmt.Sprintln(v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.errorPrint(fmt.Sprint(v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.errorPrint(fmt.Sprintf(format, v...))
}

func (l *Logger) Errorln(v ...interface{}) {
	l.errorPrint(fmt.Sprintln(v...))
}

func (l *Logger) Trace() {
	l.tracePrint()
}

func Debug(v ...interface{}) {
	defaultLogger.print(fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
	defaultLogger.print(fmt.Sprintf(format, v...))
}

func Debugln(v ...interface{}) {
	defaultLogger.print(fmt.Sprintln(v...))
}

func Print(v ...interface{}) {
	defaultLogger.debugPrint(fmt.Sprint(v...))
}

func Printf(format string, v ...interface{}) {
	defaultLogger.debugPrint(fmt.Sprintf(format, v...))
}

func Println(v ...interface{}) {
	defaultLogger.debugPrint(fmt.Sprintln(v...))
}

func Error(v ...interface{}) {
	defaultLogger.errorPrint(fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) {
	defaultLogger.errorPrint(fmt.Sprintf(format, v...))
}

func Errorln(v ...interface{}) {
	defaultLogger.errorPrint(fmt.Sprintln(v...))
}

func Trace() {
	defaultLogger.tracePrint()
}
