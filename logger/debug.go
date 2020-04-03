// +build debug

package logger

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

type PrefixType int

const (
	DebugPrefix PrefixType = iota
	ErrorPrefix
	FatalPrefix
)

var (
	debugFile *os.File
)

func InitDebug(path string) error {
	if path == "" {
		return nil
	}

	dfile, err := os.Create(path)
	if err != nil {
		return err
	}
	debugFile = dfile

	Fdebugln("################## STARTING SESSION")
	return nil
}

func Close() {
	if debugFile == nil {
		return
	}
	debugFile.Close()
}

func Fdebugln(v ...interface{}) {
	if debugFile == nil {
		return
	}

	str := fmt.Sprintf("%s %s", prefix(DebugPrefix), fmt.Sprintln(v...))
	fmt.Fprint(debugFile, str)
}

func FdebuglnError(v ...interface{}) {
	if debugFile == nil {
		return
	}

	str := fmt.Sprintf("%s %s", prefix(ErrorPrefix), fmt.Sprintln(v...))
	fmt.Fprint(debugFile, str)
}

func FdebuglnFatal(v ...interface{}) {
	if debugFile == nil {
		return
	}

	str := fmt.Sprintf("%s %s", prefix(FatalPrefix), fmt.Sprintln(v...))
	fmt.Fprint(debugFile, str)
	// fmt.FPrint(debugFile, runtimeDebug.Stack())
	panic(str)
}

// Doesn't write to the debug file
func Debugln(v ...interface{}) {
	str := fmt.Sprintf("%s %s", prefix(DebugPrefix), fmt.Sprintln(v...))
	fmt.Print(str)
}

// Doesn't write to the debug file
func DebuglnError(v ...interface{}) {
	str := fmt.Sprintf("%s %s", prefix(ErrorPrefix), fmt.Sprintln(v...))
	fmt.Print(str)
}

// Doesn't write to the debug file
func DebuglnFatal(v ...interface{}) {
	str := fmt.Sprintf("%s %s", prefix(FatalPrefix), fmt.Sprintln(v...))
	panic(str)
}

func prefix(t PrefixType) (prefix string) {
	h, m, s := time.Now().Clock()
	timePrefix := fmt.Sprintf("%d:%02d:%02d", h, m, s)

	bold := "\x1b[1m"
	red := "\x1b[31m"
	endSeq := "\x1b[0m"

	switch t {
	case DebugPrefix:
		prefix = fmt.Sprintf("%sDEBUG%s", bold, endSeq)
	case FatalPrefix:
		prefix = fmt.Sprintf("%s%sFATAL%s", red, bold, endSeq)
	case ErrorPrefix:
		prefix = fmt.Sprintf("%s%sERROR%s", red, bold, endSeq)
	default:
		prefix = fmt.Sprintf("%sDEBUG%s", bold, endSeq)
	}

	// We're using 2, to ascend 2 stack frames
	pc, _, line, _ := runtime.Caller(2)
	debugInfo := fmt.Sprintf("[%s:%d]", runtime.FuncForPC(pc).Name(), line)

	return fmt.Sprintf("%s %s %s", prefix, timePrefix, debugInfo)
}
