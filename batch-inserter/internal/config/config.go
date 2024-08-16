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

// Hints
//
// 1. use `mapstructure` in Config as if it is a yaml/json tag
// 2. viper can map time not only to string but also to time.Duration

// Configuration Constraints
//
// # Enivronment variables
// - Must be defined, otherwise application shouldn't start
// - Constant, shouldnt be overridden in runtime
// - Should not have a default value
//
// # Configuration File
// - Must exist, have same structure as config.Config, otherwise application shouldn't start
// - May be overridden in runtime or exist in multiple variants across application parts
// - Should not have a default value
//
// # Command Line Arguments
//   - May not be defined
//   - Lifetime constants, shouldnt be overridden in runtime
//   - Should be defaulted by one of the following:
//	    - Type Zero Values
//	    - [-1 , "NO" , "off"] or other kind of negative value

// Use this type to use env decode hook in configuration file
// See config/utilitary.go # envInConfigValuesHook for details
//
// # Brief example of potential usage:
//
// WORKSPACE = '~/user/goapp'
// ${WORKSPACE}/file/path => ~/user/goapp/file/path
type path string

// Represents config file, must be changed manually
// Only public fields
type Config struct {
	Logger     `mapstructure:"Logger"`
	Broker     `mapstructure:"Broker"`
	Repository `mapstructure:"Repository"`
}
type Logger struct {
	SyncTimeout time.Duration `mapstructure:"SyncTimeout"`
	Cores       []LoggerCore  `mapstructure:"Cores"`
}
type LoggerCore struct {
	Name           string `mapstructure:"Name"`           // used in LevelWithName
	EncoderLevel   string `mapstructure:"EncoderLevel"`   // production or development
	Path           path   `mapstructure:"Path"`           // everything that getLogFile can handle
	Level          int    `mapstructure:"Level"`          // might be negative, used in LevelWithName
	MustCreateCore bool   `mapstructure:"MustCreateCore"` // false = ignore if core init fails
}
type Broker struct {
	Addresses     []string `mapstructure:"Addresses"`
	BatchSize     int      `mapstructure:"BatchSize"`
	Topic         string   `mapstructure:"Topic"`
	BatchProvider `mapstructure:"BatchProvider"`
}

type BatchProvider struct {
	PreallocSize int `mapstructure:"PreallocSize"`
}
type Repository struct {
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
	Name  string
	Value any
}

// Use pflag to bind
type flagSetter func()

// Other options
type elseSetter func() error
