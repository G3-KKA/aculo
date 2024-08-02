package config

import (
	"time"

	_ "github.com/spf13/viper"
)

// Global config instance
var c Config

// Returns global config instance
func Get() Config {
	return c
}
func InitConfig() error {
	return initConfig()
}

// Hints
//
// 1. use `mapstructure` in Config as if it is a yaml/json tag
// 2. viper can map time not only to string but also to time.Duration

// Configuration Constraints
//
// # Enivronment variables
// - Must be defined, otherwise program shouldn't start
// - Lifetime constants, shouldnt be overridden in runtime
// - Cannot be defaulted
//
// # Configuration File
// - Must exist, have same structure as config.Config, otherwise program shouldn't start
// - May be overridden in runtime or exist in multiple variants across sessions
// - Cannot Be Defaulted
//
// # Command Line Arguments
//   - May not be defined
//   - Lifetime constants, shouldnt be overridden in runtime
//   - Should be defaulted by:
//	    - Type Zero Values
//	    - [-1 , "NO" , "off"] or other kind of negative value

// Use this type to use env decode hook in configuration file
// See config/utilitary.go # envInConfigValuesHook for details
//
// Brief example of usage:
// WORKSPACE = ~/user/goapp
// ${WORKSPACE}/file/path => ~/user/goapp/file/path
type path string

var environment = [...]string{
	// Every path in service works around WORKSPACE
	// Removing this will break every env-based path in service
	"WORKSPACE",
	"CONFIG_FILE",
}

// Represents config file, must be changed manually
// Only public fields
type Config struct {
	Logger     `mapstructure:"Logger"`
	HTTPServer `mapstructure:"HTTPServer"`
}
type Logger struct {
	SyncTimeout time.Duration `mapstructure:"SyncTimeout"`
	Cores       []struct {
		Name           string `mapstructure:"Name"`           // used in LevelWithName
		EncoderLevel   string `mapstructure:"EncoderLevel"`   // production or development
		Path           path   `mapstructure:"Path"`           // everything that getLogFile can handle
		Level          int    `mapstructure:"Level"`          // might be negative, used in LevelWithName
		MustCreateCore bool   `mapstructure:"MustCreateCore"` // false = ignore if core init fails
	} `mapstructure:"Cores"`
}
type HTTPServer struct {
	ListeningAddress string `mapstructure:"ListeningAddress"`
	Port             string `mapstructure:"Port"`
}

// Command line arguments, use pfalg, see example
var flags = [...]flagSetter{}

// Other viper options
var elses = [...]elseSetter{}

// Uses	viper.Set
var override = [...]overrideContainer{}
