package main

import (
	"aculo/frontend-restapi/internal/app"
	log "aculo/frontend-restapi/internal/logger"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal("unhandled error: ", err)
	}
}
