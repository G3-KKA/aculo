package main

import (
	"aculo/connector-restapi/internal/app"
	log "aculo/connector-restapi/internal/logger"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal("unhandled error: ", err)
	}
}
