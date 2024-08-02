package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Event struct {
	Me    string `json:"me" ch:"me"`
	Mario string `json:"mario" ch:"mario"`
}

func main() {
	mem := Event{Me: "mario", Mario: "me"}
	mariobytes, err := json.Marshal(mem)
	if err != nil {
		panic(err)
	}
	buf := bytes.Buffer{}
	fmt.Println("size of batch >>>> hardcoded to 1000")
	time.Sleep(150 * time.Millisecond)
	fmt.Println("address >>>> hardcoded to http://localhost:7732/event/?topic=test")
	time.Sleep(150 * time.Millisecond)

	for range 1000 {
		buf.Write(mariobytes)
		rsp, err := http.Post("http://localhost:7732/event/?topic=test", "application/json", &buf)
		if err != nil {
			panic(err)
		}
		sliice, err := io.ReadAll(rsp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(sliice))
		rsp.Body.Close()
	}
}
