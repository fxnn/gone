package log

import "log"

type StandardLogger struct {
	backend *log.Logger
}

func NewStandardLogger(backend *log.Logger) *StandardLogger {
	return &StandardLogger{backend}
}

func (l *StandardLogger) Prefix() string {
	return l.backend.Prefix()
}
func (l *StandardLogger) SetPrefix(prefix string) {
	l.backend.SetPrefix(prefix)
}

func (l *StandardLogger) Fatal(v ...interface{}) {
	l.backend.Fatal(FATAL.PrependV(v)...)
}
func (l *StandardLogger) Fatalf(format string, v ...interface{}) {
	l.backend.Fatalf(FATAL.Prepend(format), v...)
}
func (l *StandardLogger) Fatalln(v ...interface{}) {
	l.backend.Fatalln(FATAL.PrependV(v)...)
}

func (l *StandardLogger) Panic(v ...interface{}) {
	l.backend.Panic(PANIC.PrependV(v)...)
}
func (l *StandardLogger) Panicf(format string, v ...interface{}) {
	l.backend.Panicf(PANIC.Prepend(format), v...)
}
func (l *StandardLogger) Panicln(v ...interface{}) {
	l.backend.Panicln(PANIC.PrependV(v)...)
}

func (l *StandardLogger) Warn(v ...interface{}) {
	l.backend.Print(WARNING.PrependV(v)...)
}
func (l *StandardLogger) Warnf(format string, v ...interface{}) {
	l.backend.Printf(WARNING.Prepend(format), v...)
}
func (l *StandardLogger) Warnln(v ...interface{}) {
	l.backend.Println(WARNING.PrependV(v)...)
}

func (l *StandardLogger) Debug(v ...interface{}) {
	l.backend.Print(DEBUG.PrependV(v)...)
}
func (l *StandardLogger) Debugf(format string, v ...interface{}) {
	l.backend.Printf(DEBUG.Prepend(format), v...)
}
func (l *StandardLogger) Debugln(v ...interface{}) {
	l.backend.Println(DEBUG.PrependV(v)...)
}

func (l *StandardLogger) Print(v ...interface{}) {
	l.backend.Print(INFO.PrependV(v)...)
}
func (l *StandardLogger) Printf(format string, v ...interface{}) {
	l.backend.Printf(INFO.Prepend(format), v...)
}
func (l *StandardLogger) Println(v ...interface{}) {
	l.backend.Println(INFO.PrependV(v)...)
}
