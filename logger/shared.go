package logger

import (
	"fmt"
	"os"
)

const (
	bold   = "\x1b[1m"
	red    = "\x1b[31m"
	yellow = "\x1b[33m"
	endSeq = "\x1b[0m"
)

var (
	warningPrefix = fmt.Sprintf("%s%sWARNING%s", bold, yellow, endSeq)
	errorPrefix   = fmt.Sprintf("%s%sERROR%s", bold, red, endSeq)
)

func ErrorLogln(args ...interface{}) {
	t := fmt.Sprintf("%s %s", errorPrefix, fmt.Sprint(args...))
	fmt.Println(t)
}

func ErrorLoglnFatal(args ...interface{}) {
	t := fmt.Sprintf("%s %s", errorPrefix, fmt.Sprint(args...))
	fmt.Println(t)
	os.Exit(1)
}

func WarningLogln(args ...interface{}) {
	t := fmt.Sprintf("%s %s", warningPrefix, fmt.Sprint(args...))
	fmt.Println(t)
}

func Log(args ...interface{}) {
	fmt.Print(args...)
}

func Logln(args ...interface{}) {
	fmt.Println(args...)
}
