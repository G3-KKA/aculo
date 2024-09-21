package main

import (
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

const l_bench_topic = "bench_topic"

var l_address []string = []string{"localhost:9092", "localhost:9093"}

var l_reqdata = []byte(`{"id":1, "name":"test"}`)

func main() {
	wg := sync.WaitGroup{}
	for i := range 1 {
		wg.Add(1)
		f := func(ib int) {
			defer wg.Done()
			counter := 0
			go func() {
				for {
					time.Sleep(time.Second * 5)
					log.Printf("me:%d, already send: %d", ib, counter)
				}
			}()

			producer, err := sarama.NewSyncProducer(l_address, nil)
			if err != nil {
				log.Fatal(err.Error())
			}
			defer producer.Close()

			for {
				producer.BeginTxn()
				for range 100 {

					counter++

					if counter%200 == 0 {
						log.Printf("me:%d, already send: %d", ib, counter)
					}
					var msg = sarama.ProducerMessage{
						Topic: l_bench_topic,
						Value: sarama.ByteEncoder(l_reqdata),
					}
					_, _, err := producer.SendMessage(&msg)
					if err != nil {
						log.Fatalf(err.Error())
					}
				}
				producer.CommitTxn()
			}
		}
		go f(i)
	}
	wg.Wait()

}
