package log

type Level int

const (
	PANIC Level = iota
	FATAL
	ERROR
	WARNING
	INFO
	DEBUG
)

var levelNames = []string{"PANIC", "FATAL", "ERROR", "WARNI", "INFOR", "DEBUG"}

func (l Level) String() string {
	return levelNames[l]
}

func (l Level) Prepend(s string) string {
	return l.String() + " " + s
}

func (l Level) PrependV(v ...interface{}) []interface{} {
	return append([]interface{}{l}, v...)
}
