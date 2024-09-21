package aculo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
)

var _ io.WriteCloser = (*Logger)(nil)

type (
	// Core Struct for any aculo client
	Logger struct {
		client   sarama.Client
		producer sarama.AsyncProducer
		md       Metadata

		mx  sync.Mutex
		err error

		once sync.Once
		errs chan error
	}

	// Holds
	Metadata struct {
		Topic   string `json:"topic"`
		Address string `json:"address"`
	}
)

// Get metadata from master, connect to kafka
func New(ctx context.Context, masterAddress string) (*Logger, error) {

	var (
		md Metadata
	)

	getmetadata := func() (err error) {

		var (
			url  string
			resp *http.Response
			body []byte
		)

		url = fmt.Sprintf("http://%s/register", masterAddress)

		resp, err = http.Get(url)
		if err != nil {
			// TODO: Retry
			return err
		}

		body = make([]byte, resp.ContentLength)

		read, err := resp.Body.Read(body)
		defer resp.Body.Close()

		if err != nil || read != int(resp.ContentLength) {
			return
		}

		err = sonic.Unmarshal(body, &md)

		return

	}
	if err := getmetadata(); err != nil {
		return nil, err
	}

	logger := Logger{
		client:   nil,
		producer: nil,
		md:       md,
		mx:       sync.Mutex{},
		err:      nil,
		once:     sync.Once{},
		errs:     nil,
	}

	var (
		kafkaCfg *sarama.Config
		client   sarama.Client
		producer sarama.AsyncProducer
	)

	connectToKafka := func() (err error) {

		kafkaCfg = sarama.NewConfig()
		client, err := sarama.NewClient([]string{md.Address}, kafkaCfg)
		if err != nil {
			return
		}
		producer, err = sarama.NewAsyncProducerFromClient(client)
		if err != nil {
			err = errors.Join(client.Close(), err)
		}
		return
	}
	if err := connectToKafka(); err != nil {
		return nil, err

	}

	logger.client = client
	logger.producer = producer

	return &logger, nil

}

// Close underlying client, merge all errors and return
func (l *Logger) Close() (err error) {
	l.mx.Lock()
	defer l.mx.Unlock()

	err = errors.Join(l.producer.Close(), err)
	err = errors.Join(l.client.Close(), err)

	for e := range l.producer.Errors() {
		err = errors.Join(e, err)
	}

	err = errors.Join(l.err, err)

	chanClose := func() {

		close(l.errs)
		l.errs = nil
	}
	l.once.Do(chanClose)

	return

}

// Writes len(buf) as a single message to the topic
func (l *Logger) Write(buf []byte) (n int, err error) {
	msg := &sarama.ProducerMessage{
		Topic: l.md.Topic,
		//Key:       nil,
		Value: sarama.ByteEncoder(buf),
		//Headers:   []sarama.RecordHeader{},
		//Metadata:  nil,
		//Offset:    0,
		//Partition: 0,
		Timestamp: time.Now(),
	}
	// TODO: Is it safe to just write out to the channel?
	// Can deadlock happen here ?
	l.producer.Input() <- msg
	return len(buf), nil
}

func (l *Logger) Errors() <-chan error {
	panic("unimplemented (l *Logger) Errors() <-chan error ")
	/*
	 if l.errs == nil logic

	*/
}
func (l *Logger) Status() (err error) {
	l.mx.Lock()
	defer l.mx.Unlock()

	err = l.err

	if l.errs != nil {
		for len(l.errs) != 0 {
			e := <-l.errs
			err = errors.Join(e, err)
		}
		// TODO: repopulate l.errs with all non-nil errors
	}

	if l.client.Closed() {
		err = errors.Join(err, ErrClientClosed)
	}

	return
}
