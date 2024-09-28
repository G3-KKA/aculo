package config

import (
	"time"
)

// Environment variables validated automatically.
var (
	environment = [...]string{
		// Every path in service works around WORKSPACE,
		// Removing this will break every env-based path in service.
		"WORKSPACE",
		"CONFIG_FILE",
	}

	// Command line arguments, use pfalg, see example.
	flags = [...]flagSetter{}

	// Other viper options.
	elses = [...]elseSetter{}

	// Uses	viper.Set.
	override = [...]overrideContainer{}
)

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

// Signalizes that config field may contain env signature,
// and it must be replaced with value of the env.
type EnvString string

// Example:
// WORKSPACE = '~/user/goapp'
//
// ${WORKSPACE}/file/path   -->    ~/user/goapp/file/path.

// Represents expected contents of configuration file.
type Config struct {
	L Logger     `mapstructure:"Logger"`
	C Controller `mapstructure:"Controller"`
}
type (
	Logger struct {
		SyncTimeout time.Duration `mapstructure:"SyncTimeout"`
		Cores       []LoggerCore  `mapstructure:"Cores"`
	}
	LoggerCore struct {
		Name           string    `mapstructure:"Name"`
		EncoderLevel   string    `mapstructure:"EncoderLevel"` // production or development.
		Path           EnvString `mapstructure:"Path"`
		Level          int8      `mapstructure:"Level"` // might be negative.
		MustCreateCore bool      `mapstructure:"MustCreateCore"`
	}
	Controller struct {
		GRPCServer `mapstructure:"GRPCServer"`
		HTTPServer `mapstructure:"HTTPServer"`
	}
	GRPCServer struct {
		Address string `mapstructure:"Address"`
	}
	HTTPServer struct {
		Address string `mapstructure:"Address"`
	}
)
