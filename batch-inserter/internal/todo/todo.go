package todo

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"context"
	"log"
	"sync"
	"sync/atomic"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type BatchLoader struct {
	ch       clickhouse.Conn
	consumer sarama.Consumer
}

func New(ch clickhouse.Conn, consumer sarama.Consumer) *BatchLoader {
	return &BatchLoader{ch: ch, consumer: consumer}

}

type BatchProvider struct {
	ready   []atomic.Bool
	batches [][]domain.Event

	preallocSize uint
	batchSize    uint
	mx           *sync.RWMutex
}

func NewBatchProvider(ctx context.Context, config config.Config) *BatchProvider {
	bprov := &BatchProvider{
		ready:   make([]atomic.Bool, config.BatchProvider.PreallocSize),
		batches: make([][]domain.Event, config.BatchProvider.PreallocSize),

		preallocSize: config.BatchProvider.PreallocSize,
		batchSize:    config.BatchProvider.BatchSize,
		mx:           &sync.RWMutex{},
	}
	allocBatches(bprov.batches, config)

	return bprov

}
func allocBatches(batches [][]domain.Event, config config.Config) {
	for i := range len(batches) {
		batches[i] = make([]domain.Event, config.BatchProvider.BatchSize)
	}
}

type ReturnFunc func()

func (p *BatchProvider) GetBatch() ([]domain.Event, ReturnFunc) {

	p.mx.RLock()
	for i := range len(p.ready) {
		if p.ready[i].CompareAndSwap(false, true) {
			defer p.mx.RUnlock()
			return p.batches[i], func() { p.ready[i].Store(false) }
		}
	}
	p.mx.RUnlock()

	p.mx.Lock()

	p.ready = make([]atomic.Bool, len(p.ready)*2)
	p.batches = make([][]domain.Event, len(p.batches)*2)

	defer p.mx.Unlock()

	p.ready[len(p.ready)-1].Store(true)
	return p.batches[len(p.batches)-1], nil
}
func (bl *BatchLoader) DoWork(ctx context.Context) {
	partitionConsumer, err := bl.consumer.ConsumePartition("test", 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partitionConsumer.Close()
	ch := partitionConsumer.Messages()
	panic("доделать и протестировать, прямо сейчас и на месте ")
	batchProvider := NewBatchProvider(ctx, config.Config{})
	for {
		events := make([]domain.Event, 1000) // a lot of allocations
		// Здесь нужно тчо-то рводе пула слайсов чтобы избежать бессмысленных аллокаций
		// Он будет помечать уже отправленные слайсы как прочитанные и готовые к перезаписи
		// А если таких нет -- аллоцирует дополнительный слайс
		// Должно быть норм
		// Возможно здесь пригодится очередь

		for i := range 1000 {
			msg := <-ch
			//json.Unmarshal(msg.Value, &events[i])
			events[i].Data = msg.Value
			events[i].EID = uuid.New().String()
			events[i].ProviderID = "0"
			events[i].SchemaID = "0"
			events[i].Type = "test_type"
			log.Println("Event received", string(events[i].Data))
			// Вот эта горутина под конец своего исполнения пометит слайс
			// Как готовые к перезаписи
			// []atomic.Bool -- индексы готовых к перезаписи слайсов
			// [][]domain.Event -- готовые к перезаписи
			// Выглядит как thread safe

		}
		go func(events []domain.Event) {
			batch, err := bl.ch.PrepareBatch(ctx, "INSERT INTO event.main_table (eid, provider_id, schema_id, type, data)")
			if err != nil {
				log.Fatalf("Failed to prepare batch: %v", err)
			}
			log.Println("Events received")
			for _, event := range events {
				err := batch.AppendStruct(&event)
				if err != nil {
					log.Fatalf("Failed to append struct: %v", err)
				}
			}
			log.Println("Batch sent")
			err = batch.Send()
			if err != nil {
				log.Fatalf("Failed to send batch: %v", err)
			} // slice will be deallocated here, fix
		}(events)
	}
}
