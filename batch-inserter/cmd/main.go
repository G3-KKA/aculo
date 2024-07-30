package main

import (
	"aculo/batch-inserter/internal/app"
	"log"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal("unhandled error: ", err)
	}
}
