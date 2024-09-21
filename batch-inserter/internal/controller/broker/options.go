package broker

/* type (

	//
	// Used in OptionFunction pattern
	BrokerOptionFunc = option.OptionFunc[BrokerOptions]

	//
	// Every field may be modified using [BrokerOptionFunc].
	//
	// Generic setters for all fields are provided in this package.
	BrokerOptions struct {
		Namegen TopicNameGenerator
	}
)

// Return default [BrokerOptions].
func DefaultBrokerOptions() BrokerOptions {
	options := BrokerOptions{
		Namegen: &logTopicNamegen{},
	}
	return options
}

// Set custom name generator
func WithNameGenerator(namegen TopicNameGenerator) BrokerOptionFunc {
	opt := BrokerOptionFunc(func(bc *BrokerOptions) error {
		bc.Namegen = namegen
		return nil
	})
	return opt
} */
