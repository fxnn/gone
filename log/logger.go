package log

// Logger is a common interface for all possible logging backends.
type Logger interface {
	Prefix() string
	SetPrefix(prefix string)

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Warnln(v ...interface{})

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})

	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}
