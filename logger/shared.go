package logger

import (
	"fmt"
	"os"
)

func Log(s string, args ...interface{}) {
	fmt.Printf(s, args)
}

func Logln(s string, args ...interface{}) {
	// s := fmt.Sprintf("%v", args...)
	fmt.Println(s, args)
}

func LogFatal(s string, args ...interface{}) {
	Logln(s, args)
	os.Exit(1)
}