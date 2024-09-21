package broker

/*
const ( // todo, get this value dynamicaly? cfg ? cap(produser.Messager) ?
	LOG_CHANNEL_BUF_SIZE = 1000
)

type topicHandler struct {
	partConsumer sarama.PartitionConsumer
	down         atomic.Bool
} //  ===================== НУЖЕН ВОРКЕР ПУЛ ========================================

func newTopicHandler(_ context.Context) (*topicHandler, error) {
	handler := &topicHandler{
		partConsumer: nil,
	}
	return handler, nil
}

// [shuttable.Shuttable]
//
// # Closes channel on underlying partition consumer
func (handler *topicHandler) Shutdown() error {
	if handler.partConsumer == nil {
		return unierrors.ErrNotInitialisedYet

	}
	handler.partConsumer.AsyncClose()
	return nil

}

// [shuttable.Shuttable]
func (handler *topicHandler) IsShuttedDown() bool {
	return handler.down.Load()
}

// Handle reads messages from kafka and converts them to [domain.Log].
//
// Then they are send via channel
func (handler *topicHandler) Handle(topic string, consumer sarama.Consumer) (<-chan domain.Log, error) {
	logs := make(chan domain.Log, LOG_CHANNEL_BUF_SIZE)
	partConsumer, err := consumer.ConsumePartition(topic, 1, sarama.OffsetNewest)
	if err != nil {
		return nil, err
	}
	handler.partConsumer = partConsumer

	topicRountine := func() {
		var log domain.Log
		for msg := range partConsumer.Messages() {
			log.Data = msg.Value

			//
			// TODO msg.Headers
			//

			// Зачем у ивента отдельный индентификатор ? Может его лучше на всю сессию один сделать ?
			log.LogID = uuid.New().String()
			log.ProviderID = "0"
			log.SchemaID = "0"
			log.Type = "test_type"
			logs <- log
		}
		handler.down.Store(true)

	}
	go topicRountine()
	return logs, nil
}
*/
