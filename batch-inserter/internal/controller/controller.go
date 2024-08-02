package controller

import (
	"aculo/batch-inserter/internal/config"
	log "aculo/batch-inserter/internal/logger"
	"aculo/batch-inserter/internal/service"
	"context"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

//go:generate mockery --filename=mock_controller.go --name=Controller --dir=. --structname MockController  --inpackage=true
type Controller interface {
	Serve() error
}
type cntrllr struct {
	service       service.Service
	config        config.Config
	batchprovider *BatchProvider
}

func New(ctx context.Context, config config.Config, service service.Service) (Controller, error) {
	return &cntrllr{
		service:       service,
		config:        config,
		batchprovider: NewBatchProvider(ctx, config),
	}, nil
}

func (cntrllr *cntrllr) Serve() error {
	brockerCfg := cntrllr.config.Brocker
	consumer, close, err := assembleConsumer(brockerCfg)
	if err != nil {
		return err
	}
	defer close()
	ch := consumer.Messages()
	for {
		ctx := context.TODO()
		events, returnbatch := cntrllr.batchprovider.GetBatch()
		// USELESS, USE BATCH PROVIDER
		// Здесь нужно тчо-то рводе пула слайсов чтобы избежать бессмысленных аллокаций
		// Он будет помечать уже отправленные слайсы как прочитанные и готовые к перезаписи
		// А если таких нет -- аллоцирует дополнительный слайс
		// Должно быть норм
		// Возможно здесь пригодится очередь

		for i := range brockerCfg.BatchSize {
			msg := <-ch
			// TODO msg.Headers
			log.Debug("handwriting event, make it make sence, idiot, use message headers ")
			events[i].Data = msg.Value
			events[i].EID = uuid.New().String() // Зачем у ивента отдельный индентификатор ? Может его лучше на всю сессию один сделать ?
			events[i].ProviderID = "0"
			events[i].SchemaID = "0"
			events[i].Type = "test_type"
			log.Debug("Event received", string(events[i].Data))
			// Вот эта горутина под конец своего исполнения пометит слайс
			// Как готовые к перезаписи
			// []atomic.Bool -- индексы готовых к перезаписи слайсов
			// [][]domain.Event -- готовые к перезаписи
			// Выглядит как thread safe

		}
		go func() {
			defer returnbatch()
			err := cntrllr.service.SendBatch(ctx, events)
			if err != nil {
				log.Info("failed to send batch: %v", err)
			}
		}()

	}
	// service.SendBatch(ctx, batch)
}

type mustDeferFunc func()

func assembleConsumer(cfg config.Brocker) (sarama.PartitionConsumer, mustDeferFunc, error) {
	consumer, err := sarama.NewConsumer(cfg.Addresses, nil)
	if err != nil {
		return nil, nil, err
	}
	partitionConsumer, err := consumer.ConsumePartition(cfg.Topic, 0, sarama.OffsetOldest)
	if err != nil {
		consumer.Close()
		return nil, nil, err
	}
	return partitionConsumer, func() {
		err := partitionConsumer.Close()
		if err != nil {
			log.Info("failed to close partition consumer: %v", err)
		}
		err = consumer.Close()
		if err != nil {
			log.Info("failed to close consumer: %v", err)
		}
	}, nil
}
