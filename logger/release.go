// +build !debug

package logger

func InitDebug(path string) error    { return nil }
func Close()                         {}
func Fdebugln(v ...interface{})      {}
func FdebuglnError(v ...interface{}) {}
func FdebuglnFatal(v ...interface{}) {}
func Debugln(v ...interface{})       {}
func DebuglnError(v ...interface{})  {}
func DebuglnFatal(v ...interface{})  {}
