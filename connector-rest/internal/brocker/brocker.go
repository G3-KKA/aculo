package brocker

import (
	"aculo/connector-restapi/internal/config"
	"context"

	"github.com/IBM/sarama"
)

type SendEventRequest struct {
	Topic string
	Event []byte
}
type SendEventResponse struct {
}

//go:generate mockery --filename=mock_brocker.go --name=Brocker --dir=. --structname MockBrocker  --inpackage=true
type Brocker interface {
	SendEvent(context.Context, SendEventRequest) (SendEventResponse, error)
}

type BuildBrockerRequest struct {
}

// TODO : find a way to close connection,
// maybe connection pool, like in PostgresPoll
func New(ctx context.Context, config config.Config, req BuildBrockerRequest) (Brocker, error) {
	producer, err := sarama.NewAsyncProducer(config.Kafka.Addresses, nil)
	if err != nil {
		return nil, err
	}
	repo := &eBroker{

		producer: producer,
	}
	return repo, nil

}

type eBroker struct {
	producer sarama.AsyncProducer
}

// GetEvent implements EventBrocker.
func (brocker *eBroker) SendEvent(ctx context.Context, req SendEventRequest) (SendEventResponse, error) {
	brocker.producer.Input() <- &sarama.ProducerMessage{
		Topic: req.Topic,
		Value: sarama.ByteEncoder(req.Event),
	}
	return SendEventResponse{}, nil
}
