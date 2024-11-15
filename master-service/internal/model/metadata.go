package model

// KafkaMetadata is an information returned to the client,
// for its future connection to the Kafka.
type KafkaMetadata struct {
	Address string
	Topic   string
}
