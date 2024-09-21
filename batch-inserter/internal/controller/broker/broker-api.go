package broker

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/generics/streampool"
	"context"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type brokerapi struct {
	consumer sarama.Consumer
	admin    sarama.ClusterAdmin
	namegen  TopicNameGenerator

	pool *streampool.Pool
	mx   sync.Mutex
}

// DeleteTopic implements BrokerAPI.
func (b *brokerapi) DeleteTopic(ctx context.Context, topic string) error {
	b.mx.Lock()
	defer b.mx.Unlock()
	err := b.admin.DeleteTopic(topic)
	if err != nil {
		return err
	}
	return nil

}

// deletes
func (b *brokerapi) StopHandling(ctx context.Context, topic string) error {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.pool.StopWorkerWait(topic)
	return nil
}

const noexpot_READER_TIMER_LIMIT = time.Second * 30

// Expects that topic already exist
func (b *brokerapi) HandleTopic(ctx context.Context, topic string) (<-chan domain.Log, error) {
	b.mx.Lock()
	defer b.mx.Unlock()
	logs := make(chan domain.Log, 1000)
	partConsumer, err := b.consumer.ConsumePartition(topic, 1, sarama.OffsetNewest)
	if err != nil {
		return nil, err
	}
	worker := streampool.PoolFunc(func(stop <-chan struct{}) {
		defer close(logs)
		defer partConsumer.Close()
		var (
			log domain.Log
			msg *sarama.ConsumerMessage
			ok  bool
		)
		timer := time.NewTimer(time.Hour * 24) // just a big value
		for {
			select {
			case <-stop:
				return
			case msg, ok = <-partConsumer.Messages():
			}
			if !ok {
				return
			}
			log.Data = msg.Value

			//
			// TODO msg.Headers
			//

			log.LogID = uuid.New().String() // Зачем у ивента отдельный индентификатор ? Может его лучше на всю сессию один сделать ?
			log.ProviderID = "0"
			log.SchemaID = "0"
			log.Type = "test_type"
			timer.Reset(noexpot_READER_TIMER_LIMIT)
			select {
			case logs <- log:
			case <-timer.C:
				return
			}

		}
	})
	b.pool.Go(topic, worker)
	return logs, nil
}

// Create new topic, not initialising handler
func (b *brokerapi) NewTopic(ctx context.Context) (string, error) {
	b.mx.Lock()
	defer b.mx.Unlock()

	topic := b.namegen.Generate()
	validate := false
	err := b.admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, validate)
	if err != nil {
		return "", err
	}
	return topic, nil

}

func (b *brokerapi) shutdown() error {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.pool.ShutdownWait()

	err := b.admin.Close()

	return err
}
