// +build !debug

package logger

func InitDebug(path string) {}
func Close() {}
func Fdebugln(s string) {}
func FdebuglnFatal(s string) {}
