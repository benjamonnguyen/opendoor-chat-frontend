package devlog

import "log"

var isEnabled bool

func Enable(b bool) {
	isEnabled = b
}

func Println(v ...any) {
	if isEnabled {
		log.Println(v...)
	}
}

func Printf(format string, v ...any) {
	if isEnabled {
		log.Printf(format, v...)
	}
}
