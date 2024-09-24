package config

import (
	"reflect"
	"time"

	_ "github.com/spf13/viper"
)

func ReadInConfig() (Config, error) {
	err := initConfig()
	if err != nil {
		return Config{}, err
	}
	return c, nil
}
func Get() (Config, error) {
	if reflect.DeepEqual(c, Config{}) {
		return Config{}, ErrZeroValueConfig
	}
	return c, nil
}

// Configuration Constraints
//
// # Enivronment variables
//   - Must be defined, otherwise application shouldn't start
//   - Constant, shouldnt be overridden in runtime
//   - Should not have a default value
//
// # Configuration File
//   - Must exist, have same structure as config.Config, otherwise application shouldn't start
//   - May be overridden in runtime or exist in multiple variants across application parts
//   - Should not have a default value
//
// # Command Line Arguments
//   - May not be defined
//   - Lifetime constants, shouldnt be overridden in runtime
//   - Should be defaulted by one of the following:
//	    - Type Zero Values
//	    - [-1 , "NO" , "off"] or other kind of negative value

// Signalises that config field may contain env signature,
// and it must be replaced with value of the env.
//
// WORKSPACE = '~/user/goapp'
//
// ${WORKSPACE}/file/path   -->    ~/user/goapp/file/path
type EnvString string

// Represents config file, must be changed manually
// Only public fields
type Config struct {
	L Logger `mapstructure:"Logger"`
}
type Logger struct {
	SyncTimeout time.Duration `mapstructure:"SyncTimeout"`
	Cores       []LoggerCore  `mapstructure:"Cores"`
}
type LoggerCore struct {
	Name           string    `mapstructure:"Name"`
	EncoderLevel   string    `mapstructure:"EncoderLevel"` // production or development
	Path           EnvString `mapstructure:"Path"`
	Level          int       `mapstructure:"Level"` // might be negative
	MustCreateCore bool      `mapstructure:"MustCreateCore"`
}

// Environment variables validates automatically
var environment = [...]string{
	// Every path in service works around WORKSPACE
	// Removing this will break every env-based path in service
	"WORKSPACE",
	"CONFIG_FILE",
}

// Command line arguments, use pfalg, see example
var flags = [...]flagSetter{}

// Other viper options
var elses = [...]elseSetter{}

// Uses	viper.Set
var override = [...]overrideContainer{}

// Container for config override
type overrideContainer struct {
	Key   string
	Value any
}

// Use pflag to bind
type flagSetter func()

// Other options
type elseSetter func() error
