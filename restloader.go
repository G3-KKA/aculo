package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
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
	fmt.Println("size of batch >>>> hardcoded to 1000")
	time.Sleep(150 * time.Millisecond)
	fmt.Println("address >>>> hardcoded to http://localhost:7732/event/?topic=test")
	time.Sleep(150 * time.Millisecond)
	wg := sync.WaitGroup{}
	wg.Add(1000)
	for range 1000 {
		go func() {
			buf := bytes.Buffer{}
			for range 1000 {

				defer wg.Done()
				buf.Write(mariobytes)
				rsp, err := http.Post("http://localhost:7732/event/?topic=test", "application/json", &buf)
				buf.Reset()
				if err != nil {
					panic(err)
				}
				rsp.Body.Close()

				//fmt.Println(string(sliice))
			}
		}()
	}
	wg.Wait()
}
