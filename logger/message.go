package logger

type level string

const (
	errLevel     level = "error"
	warningLevel level = "warning"
	fatalLevel   level = "fatal"
	infoLevel    level = "info"
)

type message struct {
	level     level
	message   string
	sentLocal bool
	sentGrpc  bool
}
