package brocker

import (
	"aculo/connector-restapi/internal/config"
	"aculo/connector-restapi/internal/request"
	"context"

	"github.com/IBM/sarama"
)

//go:generate mockery --filename=mock_brocker.go --name=Brocker --dir=. --structname MockBrocker  --inpackage=true
type Brocker interface {
	SendEvent(ctx context.Context, req request.SendEventRequest) (request.SendEventResponse, error)
}

// TODO : find a way to close connection,
// maybe connection pool, like in PostgresPoll
func New(ctx context.Context, config config.Config) (*eBroker, error) {
	producer, err := sarama.NewAsyncProducer(config.Kafka.Addresses, nil)
	if err != nil {
		return nil, err
	}

	brocker := &eBroker{

		producer: producer,
	}
	return brocker, nil

}

type eBroker struct {
	producer       sarama.AsyncProducer // ЛИШНЕЕ !!!
	internalChanel chan *sarama.ProducerMessage
}

// GetEvent implements EventBrocker.
func (brocker *eBroker) SendEvent(ctx context.Context, req request.SendEventRequest) (request.SendEventResponse, error) {

	brocker.producer.Input() <- &sarama.ProducerMessage{
		Topic: req.Topic,
		Value: sarama.ByteEncoder(req.Event),
	}
	return request.SendEventResponse{}, nil
}
