package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	for _ = range 1000 {
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