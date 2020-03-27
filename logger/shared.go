package logger

import (
	"fmt"
	"os"
)


func Log(s string, args ...interface{}) {
	fmt.Printf(s, args)
}

func Logln(args ...interface{}) {
	fmt.Println(args...)
}

func LogFatal(args ...interface{}) {
	Logln(args...)
	os.Exit(1)
}

