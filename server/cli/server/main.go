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
	defer func() {
		if err := deferFunc(); err != nil {
			log.Fatalln(err)
		}
	}()
	err := server.Run(os.Args)
	if err != nil {
		log.Println(err)
	}
}
