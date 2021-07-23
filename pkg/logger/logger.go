package logger

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"myst/pkg/config"
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
	defaultLogger = NewLogger("SERVER", DefaultColor)
}

func Colorize(s string, c Color) string {
	return aurora.Colorize(s, aurora.Color(c)).String()
}

func NewLogger(prefix string, color Color) *Logger {
	prefix = strings.ToUpper(prefix)
	prefix = fmt.Sprintf("[%s] ", prefix)
	prefix = Colorize(prefix, color)
	return &Logger{
		stdout:   log.New(StdoutWriter, prefix, 0),
		stderr:   log.New(StderrWriter, prefix, log.Lshortfile),
		debugLog: log.New(DebugLogWriter, prefix, log.Ldate|log.Ltime|log.Lmicroseconds),
		errorLog: log.New(ErrorLogWriter, prefix, log.Ldate|log.Ltime|log.Lmicroseconds),
	}
}

func (l *Logger) print(s string) {
	s = strings.TrimRight(s, "\n")
	_ = l.stdout.Output(3, s)
}

func (l *Logger) debugPrint(s string) {
	s = strings.TrimRight(s, "\n")
	_ = l.debugLog.Output(3, s)
	if config.Debug {
		_ = l.stdout.Output(3, s)
	}
}

func (l *Logger) errorPrint(s string) {
	s = strings.TrimRight(s, "\n")
	s = Colorize(s, RedFg)
	_ = l.errorLog.Output(3, s)
	if config.Debug {
		_ = l.stderr.Output(3, s)
	}
}

func (l *Logger) tracePrint() {
	stack := debug.Stack()
	s := Colorize(string(stack), RedFg)
	_ = l.errorLog.Output(3, s)
	if config.Debug {
		_ = l.stderr.Output(3, s)
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
