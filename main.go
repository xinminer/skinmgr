package main

import (
	"os"
	"skinmgr/cmd"
	"skinmgr/internal/log"
)

func main() {
	if err := cmd.NewApp().Run(os.Args); err != nil {
		log.Log.Errorf("start app error: %v", err)
		os.Exit(1)
	}
}
