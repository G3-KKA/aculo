package EXAMPLE_DO_NOT_COPY

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type path string

type Config struct {
	BasicField  string `mapstructure:"basic_field"`
	Example     path   `mapstructure:"example"` // ENV hook will be applied
	InnerStruct struct {
		Field time.Duration `mapstructure:"field"` // works as intended
	} `mapstructure:"inner_struct"`
}

var environment = [...]string{
	// Every path in service works around WORKSPACE
	// Removing this will break every env-based path in service
	"WORKSPACE",
	"CONFIG_FILE",
	// additional envs
	"GOVERSION",
	"OS",
}

var flags = [...]flagSetter{
	// parceable flags, defaults are negative
	func() { pflag.Bool("enable_debug", false, "Define if debug info is enabled") },
}

// Other viper options
var elses = [...]elseSetter{
	func() error {

		// Callback on config change
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Println("Config file changed:", e.Name)
		})
		viper.WatchConfig()

		return nil
	},
	func() error {
		// Aliases example
		viper.RegisterAlias("enable_debug", "debug")
		return nil
	},
}

// Uses	viper.Set
var toOverride = [...]overrideContainer{
	// type-free override, dangerous
	{"enable_debug", true},
	{"dummy", "dummy"},
}
