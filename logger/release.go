// +build !debug

package logger

func InitDebug(path string)          {}
func Close()                         {}
func Fdebugln(v ...interface{})      {}
func FdebuglnFatal(v ...interface{}) {}
