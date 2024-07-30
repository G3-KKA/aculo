package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	for i := range 10 {
		wg.Add(1)
		go func() {
			for {
				time.Sleep(time.Millisecond * 10)

				resp, err := http.Get("http://localhost:8080/ping")
				if err != nil {
					fmt.Println("err", err)
					fmt.Println("routine id", i)
					wg.Done()
					return
				}
				defer resp.Body.Close()
				bytees, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("err", err)
					fmt.Println("routine id", i)
					wg.Done()
					return
				}
				fmt.Println(string(bytees))
			}
		}()
	}
	wg.Wait()

}
