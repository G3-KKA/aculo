package config

type Config struct {
	Kafka         `mapstructure:"Kafka"`
	Clickhouse    `mapstructure:"Clickhouse"`
	BatchProvider `mapstructure:"BatchProvider"`
}
type Kafka struct {
}
type Clickhouse struct {
}
type BatchProvider struct {
	PreallocSize uint `mapstructure:"PreallocSize"`
	BatchSize    uint `mapstructure:"BatchSize"`
}
