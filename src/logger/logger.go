package logger

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"
)

var prefix = "Go"

func SetPrefix(newPrefix string) {
	prefix = newPrefix
}

type Logger struct {
	level   Level
	context string
}

func (l *Logger) printLog(level Level, message string, color string, args ...interface{}) {
	if !l.isLevelEnabled(level) {
		return
	}
	timestamp := time.Now()

	var b bytes.Buffer

	b.WriteString(color)
	b.WriteString("[" + prefix + "]")
	b.WriteByte(' ')
	b.WriteString(strconv.Itoa(os.Getpid()))
	b.WriteString(" - ")
	b.WriteString(colors["reset"])
	b.WriteString(timestamp.Format(time.DateTime))
	b.WriteString(color)
	b.WriteByte(' ')
	b.WriteString(levelName[level])
	b.WriteByte(' ')
	b.WriteString(colors["yellow"])
	b.WriteString("[" + l.context + "]")
	b.WriteByte(' ')
	b.WriteString(color)
	b.WriteString(message)
	b.WriteString(colors["reset"])
	l.printArgs(&b, args...)

	fmt.Println(string(b.Bytes()))
}

func (l *Logger) printArgs(buff *bytes.Buffer, args ...interface{}) {
	if args == nil || len(args) == 0 {
		return
	}

	for _, arg := range args {
		switch value := arg.(type) {
		case error:
			buff.WriteByte('\n')
			buff.WriteString(value.Error())
		case struct{}:
		case interface{}:
			buff.WriteByte('\n')
			buff.WriteString(fmt.Sprintf("%s", value))
		}
	}
}

var levelName = map[Level]string{
	ErrorLevel: "ERROR",
	WarnLevel:  " WARN",
	InfoLevel:  " INFO",
	DebugLevel: "DEBUG",
}

func (l *Logger) isLevelEnabled(level Level) bool {
	return level <= l.level
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.printLog(InfoLevel, message, colors["green"], args...)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.printLog(WarnLevel, message, colors["yellow"], args...)
}

func (l *Logger) Error(message string, err error, args ...interface{}) {
	l.printLog(ErrorLevel, message+": "+err.Error(), colors["red"], append([]interface{}{err}, args...)...)
}

func (l *Logger) Fatal(message string, err error, args ...interface{}) {
	l.printLog(ErrorLevel, message+": "+err.Error(), colors["red"], append([]interface{}{err}, args...)...)
	os.Exit(1)
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.printLog(DebugLevel, message, colors["cyan"], args...)
}

func New(context string) *Logger {
	return &Logger{
		level:   defaultLogger.level,
		context: context,
	}
}
