package main

import (
	// "net/http"
	// _ "net/http/pprof"

	"foundry/cli/cmd"
	"foundry/cli/config"
	"foundry/cli/logger"
)

func main() {
	// time.Sleep(time.Second * 20)

	if err := config.Init(); err != nil {
		logger.ErrorLoglnFatal("Couldn't init config", err)
	}

	// go func() {
	// 	if err := http.ListenAndServe("localhost:7777", nil); err != nil {
	// 		logger.ErrorLoglnFatal(err)
	// 	}
	// }()

	cmd.Execute()
}
