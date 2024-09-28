package config

import "reflect"

// # ReadInConfig reads in internal global config from CONFIG_FILE.
//
// Returns it or error if got any.
func ReadInConfig() (Config, error) {
	err := initConfig()
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

// # Get returns already initialized config
//
// Return [ErrZeroValueConfig] if [ReadInConfig] not been called at least once or failed.
func Get() (Config, error) {
	var zero Config
	if reflect.DeepEqual(c, zero) {
		return Config{}, ErrZeroValueConfig
	}

	return c, nil
}
