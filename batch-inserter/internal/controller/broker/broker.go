package broker

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/generics/streampool"
	"aculo/batch-inserter/internal/interfaces/shuttable"
	"aculo/batch-inserter/internal/interfaces/txface"
	"aculo/batch-inserter/internal/logger"
	"context"
	"sync"
	"sync/atomic"

	"github.com/IBM/sarama"
)

// Static check
func _() {
	var ( // [Broker]
		_ shuttable.Shuttable          = (*broker)(nil)
		_ txface.Tx[*brokerapi]        = (*broker)(nil)
		_ txface.ApiWrapper[BrokerAPI] = (*broker)(nil)
	)
	var _ BrokerAPI = (*brokerapi)(nil)
}

//
//go:generate mockery --filename=mock_brokerapi.go --name=BrokerAPI --dir=. --structname MockBrokerAPI  --inpackage=true
type BrokerAPI interface {
	//
	// # Initialises topic in underlying MQ
	NewTopic(ctx context.Context) (topic string, err error)
	//
	// # Starts to accumulate messages from given topic
	HandleTopic(ctx context.Context, topic string) (logs <-chan domain.Log, err error)
	//
	// # Stop sending logs to the channel and close it
	StopHandling(ctx context.Context, topic string) (err error)
	//
	// # Delete topic in underlying MQ
	DeleteTopic(ctx context.Context, topic string) (err error)
}

//
//go:generate mockery --filename=mock_broker.go --name=Broker --dir=. --structname MockBroker  --inpackage=true
type Broker interface {
	//
	// Restrict access to broker api via Tx()
	txface.Tx[*brokerapi]
	//
	// Other parts of program *should* accept behaviour, so we provide it as well
	txface.ApiWrapper[BrokerAPI]
	//
	// Everything should
	shuttable.Shuttable
}

// Returns brokerapi wrapped into transaction mechanism
func New(
	ctx context.Context,
	config config.Config,
	logger logger.Logger,
	/* opts ...BrokerOptionFunc, */
) (*broker, error) {

	/* 	brokerOptions := DefaultBrokerOptions()

	   	for _, opt := range opts {
	   		err := opt(&brokerOptions)
	   		if err != nil {
	   			return nil, err
	   		}
	   	} */

	saramaConfig := sarama.NewConfig()

	//
	// Connecting to Kafka
	client, err := sarama.NewClient(config.Broker.Addresses, saramaConfig)
	if err != nil {
		return nil, err
	}
	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		client.Close()
		return nil, err
	}
	// Note:  admin.Close() will also close consumer !
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		client.Close()
		return nil, err
	}

	b := &brokerapi{
		consumer: consumer,
		admin:    admin,
		/* namegen:  brokerOptions.Namegen, */
		pool: streampool.NewStreamPool(),
		mx:   sync.Mutex{},
	}
	tx := &broker{
		b:           b,
		unavailable: atomic.Bool{},
		mx:          sync.RWMutex{},
	}
	return tx, nil

}
