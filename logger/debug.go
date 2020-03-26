// +build debug

package logger

import (
	"fmt"
)

func Debugf(s string, args ...interface{}) {
	// fmt.Printf("DEBUG: " + s, args...)
}

func Fdebugln(s string) {
	// TODO: Open file once
	// use mutex lock

	df, err := dfile()
	if err != nil {
		panic(err)
	}
	defer df.Close()

	fmt.Fprintln(df, s)
}

func Debugln(args ...interface{}) {
	// s := fmt.Sprintf("DEBUG: %s", args...)
	// fmt.Println(s)
}
