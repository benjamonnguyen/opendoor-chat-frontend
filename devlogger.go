package main

import "log"

type DevLogger struct {
	isEnabled bool
}

func NewDevLogger(enable bool) *DevLogger {
	return &DevLogger{
		isEnabled: enable,
	}
}

func (l *DevLogger) Println(v ...any) {
	if l.isEnabled {
		log.Println(v...)
	}
}

func (l *DevLogger) Printf(format string, v ...any) {
	if l.isEnabled {
		log.Printf(format, v...)
	}
}
