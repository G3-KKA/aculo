package todo

import (
	"aculo/batch-inserter/domain"
	"context"
	"log"

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

func (bl *BatchLoader) DoWork(ctx context.Context) {
	partitionConsumer, err := bl.consumer.ConsumePartition("test", 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partitionConsumer.Close()
	ch := partitionConsumer.Messages()
	for {
		// USELESS, USE BATCH PROVIDER
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
		// Repository logic ( aka service.InsertBatch(ctx, batch) => repo.InsertBatch(ctx, batch) )
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
