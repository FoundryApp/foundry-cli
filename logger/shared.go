package logger

import (
	"fmt"
	"os"
)

const (
	dfilePath = "/Users/vasekmlejnsky/Developer/foundry/cli/debug.txt"
)

func init() {
	if exists, err := dfileExists(); exists && err != nil {
		// File exists, delete it
		if err := os.Remove(dfilePath); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
}

func dfileExists() (bool, error) {
	if _, err := os.Stat(dfilePath); os.IsNotExist(err) {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

func dfile() (*os.File, error) {

	if exists, err := dfileExists(); exists && err == nil {
		// File exists, open it
		f, err := os.Open(dfilePath)
		if err != nil {
			return nil, err
		}
		return f, nil
	} else if !exists && err == nil {
		// File doesn't exist, create it
		f, err := os.Create(dfilePath)
		if err != nil {
			return nil, err
		}
		return f, nil
	} else {
		// No idea if file exists, got an error
		return nil, err
	}
}

func Log(s string, args ...interface{}) {
	fmt.Printf(s, args)
}

func Logln(args ...interface{}) {
	s := fmt.Sprintf("%s", args...)
	fmt.Println(s)
}

func LogFatal( args ...interface{}) {
	Logln(args)
	os.Exit(1)
}

