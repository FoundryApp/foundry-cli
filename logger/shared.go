package logger

import "fmt"

func Log(s string, args ...interface{}) {
	fmt.Printf(s, args)
}