package main

import (
	"log"
	"os"
	"runtime"

	"github.com/bububa/visiondb/server/app"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	server, deferFunc := app.NewApp()
	defer deferFunc()
	err := server.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
