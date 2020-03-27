package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"github.com/sht/myst/go/config"
	"log"
	"os"
	"strings"
	"time"
)

var routesPrefix = "github.com/sht/myst/go/routes."

var accessLogWriter *os.File
var errorLogWriter *os.File

var StdoutLogger logger
var StderrLogger logger
var AccessLogger logger
var ErrorLogger logger

type logger struct {
	logger *log.Logger
}

func (l *logger) Write(p []byte) (n int, err error) {
	str := color.Sprint(string(p))
	l.logger.Print(str)
	// color.Lprint(l.logger, string(p))
	return 0, nil
}

func Init() error {
	// Create logs folder if it doesn't already exist
	_, err := os.Stat("logs")
	if os.IsNotExist(err) {
		err = os.Mkdir("logs", 0666)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	// Assign access and error log writers to their respective files
	accessLogWriter, err = os.OpenFile("logs/access.log",
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	errorLogWriter, err = os.OpenFile("logs/error.log",
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	flag := log.Ldate | log.Ltime | log.Lmicroseconds

	// Create goroutine-safe loggers
	StdoutLogger.logger = log.New(os.Stdout, "", 0)
	StderrLogger.logger = log.New(os.Stderr, "", 0)
	AccessLogger.logger = log.New(accessLogWriter, "", flag)
	ErrorLogger.logger = log.New(errorLogWriter, "", flag)

	return nil
}

func Printf(domain, format string, a ...interface{}) {
	a = append([]interface{}{domain}, a...)
	StdoutLogger.logger.Printf("[%s] "+format, a...)
}

func Logf(domain, format interface{}, a ...interface{}) {
	switch format.(type) {
	case string:
	case error:
		format = format.(error).Error()
	default:
		format = ""
	}
	a = append([]interface{}{domain}, a...)
	AccessLogger.logger.Printf("[%s] "+format.(string)+"\n", a...)
	if config.Debug {
		// Also log to console if debugging
		StdoutLogger.logger.Printf("[%s] "+format.(string), a...)
	}
}

func Errorf(domain, format interface{}, a ...interface{}) {
	switch format.(type) {
	case string:
	case error:
		format = format.(error).Error()
	default:
		format = ""
	}
	a = append([]interface{}{domain}, a...)
	ErrorLogger.logger.Printf("[%s] "+format.(string)+"\n", a...)
	if config.Debug {
		// Also log to console if debugging
		StderrLogger.logger.Printf("[%s] "+format.(string)+"\n", a...)
	}
}

func LogFormatter(param gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("[GIN] %3d | %13v | %15s | %-7s %s\n%s",
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
}

func ConsoleFormatter(param gin.LogFormatterParams) string {
	statusColor := param.StatusCodeColor()
	methodColor := param.MethodColor()
	resetColor := param.ResetColor()

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("[GIN] %s %3d %s | %13v | %15s | %s %-7s %s %s\n%s",
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

func PrintRoutes(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	handlerName = strings.TrimPrefix(handlerName, routesPrefix)
	Printf("ROUTES", "%-7s %-25s --> %6s", httpMethod, absolutePath, handlerName)
}
