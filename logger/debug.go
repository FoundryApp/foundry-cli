// +build debug

package logger

import (
	"fmt"
	"os"
	"runtime"
)

type PrefixType int

const (
	DebugPrefix PrefixType = iota
	FatalPrefix
)

var (
	debugFile *os.File

	fatalPrefix = ""
	debugPrefix = ""
)

func InitDebug(path string) {
	if path == "" { return }

	dfile, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	debugFile = dfile
}

func Close() {
	if debugFile == nil { return }
	debugFile.Close()
}

func Fdebugln(v ...interface{}) {
	if debugFile == nil { return }

	str := fmt.Sprintf("%s %s", prefix(DebugPrefix), fmt.Sprintln(v...))
	fmt.Fprint(debugFile, str)
}

func FdebuglnFatal(v ...interface{}) {
	if debugFile == nil { return }

	str := fmt.Sprintf("%s %s", prefix(FatalPrefix), fmt.Sprintln(v...))
	fmt.Fprint(debugFile, str)
	os.Exit(1)
}

func prefix(t PrefixType) (prefix string) {
	bold 		:= "\x1b[1m"
	red 		:= "\x1b[31m"
	endSeq 	:= "\x1b[0m"

	switch t {
	case DebugPrefix:
		prefix = fmt.Sprintf("%sDEBUG%s", bold, endSeq)
	case FatalPrefix:
		prefix = fmt.Sprintf("%s%sFATAL%s", red, bold, endSeq)
	default:
		prefix = fmt.Sprintf("%sDEBUG%s", bold, endSeq)
	}

	// We're using 2, to ascend 2 stack frames
	pc, _, line, _ := runtime.Caller(2)
	debugInfo := fmt.Sprintf("[%s:%d]", runtime.FuncForPC(pc).Name(), line)

	return fmt.Sprintf("%s %s", prefix, debugInfo)
}
