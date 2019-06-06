package gorabbitgo

import "fmt"

// Logger is used for logging gorabbitgo related messages
type Logger interface {
	Print(v ...interface{})
}

type NullLogger struct {}

func (l *NullLogger) Print(v ...interface{}) {}

type StdOutLogger struct {}

func (l *StdOutLogger) Print(v ...interface{}) {
	fmt.Println(v)
}