// +build debug

package logger

import "log"

func Debugf(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}

func Debugln(args ...interface{}) {
	log.Println(args)
}
