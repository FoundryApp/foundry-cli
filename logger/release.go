// +build !debug

package logger

import io

func Debugf(fmt string, args ...interface{}) {}
// func Fdebugf(w io.Writer, f string, a ...interface) {}
func Fdebugln(s string) {}
func Debugln(s string, args ...interface{}) {}

