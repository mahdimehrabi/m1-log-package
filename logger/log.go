package logger

import (
	logv1 "github.com/mahdimehrabi/m1-log-proto/gen/go/log/v1"
	"google.golang.org/grpc"
)

type Logger struct {
	sender
}

func NewLogger(cc grpc.ClientConnInterface) *Logger {
	lsc := logv1.NewLogServiceClient(cc)
	return &Logger{
		sender: sender{
			lsc:   lsc,
			queue: make(chan message, queueMaxLength),
		},
	}
}

func (l *Logger) Error(err error) {
	l.send(message{
		level: errLevel, message: err.Error(),
	})
}

func (l *Logger) Info(msg string) {
	l.send(message{
		level: infoLevel, message: msg,
	})
}

func (l *Logger) Warning(msg string) {
	l.send(message{
		level: warningLevel, message: msg,
	})
}

func (l *Logger) Fatal(msg string) {
	l.send(message{
		level: fatalLevel, message: msg,
	})
}
