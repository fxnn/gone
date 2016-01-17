package log

import (
	"log"
	"os"
)

// DefaultLogger is the Logger instance used by the global funcs in this package
var DefaultLogger Logger = NewStandardLogger(log.New(os.Stdout, "", log.LstdFlags))

func Prefix() string {
	return DefaultLogger.Prefix()
}
func SetPrefix(prefix string) {
	DefaultLogger.SetPrefix(prefix)
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

func Warn(v ...interface{}) {
	DefaultLogger.Warn(v...)
}
func Warnf(format string, v ...interface{}) {
	DefaultLogger.Warnf(format, v...)
}
func Warnln(v ...interface{}) {
	DefaultLogger.Warnln(v...)
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

func Print(v ...interface{}) {
	DefaultLogger.Print(v...)
}
func Printf(format string, v ...interface{}) {
	DefaultLogger.Printf(format, v...)
}
func Println(v ...interface{}) {
	DefaultLogger.Println(v...)
}
