package config

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Phase of config initialisation
type initPhase func() error

// Container for config override
type overrideContainer struct {
	name  string
	value any
}

// Use pflag to bind
type flagSetter func()

// Other options
type elseSetter func() error

func execute(pipeline []initPhase) error {
	for _, phase := range pipeline {
		err := phase()
		if err != nil {
			return fmt.Errorf("cannot init condfig %e", err)
		}
	}
	return nil
}

// Adds validation to env binding
func registerENV(input ...string) error {
	viper.BindEnv(input...)
	for _, env := range input {
		// Type-free validation
		// Not defined integer or bool should be "" as well
		envalue := viper.GetString(env)
		if envalue == "" {
			return fmt.Errorf("%s is not defined", env)
		}
	}
	return nil
}

// Wraps viper.BindPFlags()
func bindFlags() error {
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return fmt.Errorf("cannot bind flags: %v", err)
	}
	return nil
}
func fillGlobalConfig() error {

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("cannot read config file: %v", err)
	}

	// Will be called one after another
	// Do not try to put them separately
	// ComposeDecode in crucial
	hooks := []mapstructure.DecodeHookFunc{
		envReplaceHook(),
		mapstructure.StringToTimeDurationHookFunc(),
	}
	composeHook := mapstructure.ComposeDecodeHookFunc(hooks...)
	err = viper.Unmarshal(&c, viper.DecodeHook(composeHook))
	if err != nil {
		return fmt.Errorf("cannot unmarshal config file contents: %v", err)
	}
	return nil
}

// Parse config file path for  ext
// TODO filepath.EXT()
func extFromPath(path string) string {
	dotIndex := strings.LastIndexByte(path, '.')
	if dotIndex == -1 {
		return ""
	}
	return path[dotIndex+1:]
}

// Parse config file path for name
func nameFromPath(path string) string {
	dotIndex := strings.LastIndexByte(path, '.')
	if dotIndex == -1 {
		return ""
	}
	slashIndex := strings.LastIndexByte(path[:dotIndex], '/')
	return path[slashIndex+1 : dotIndex]
}

// Sets config file name and extention
// Change only if something breaks
func handleConfigFile() error {
	configFileEnv := viper.GetString("CONFIG_FILE")

	// Getting parts of path
	name := nameFromPath(configFileEnv)
	ext := extFromPath(configFileEnv)
	dir := filepath.Dir(configFileEnv)

	// Setting Config
	viper.AddConfigPath(dir)
	viper.SetConfigName(name)
	viper.SetConfigType(ext)
	return nil
}

// Parse ${ENV}/dir/file kind of path,
// Only works if variable type is path, see ./config.go
func envReplaceHook() mapstructure.DecodeHookFuncType {
	hook := mapstructure.DecodeHookFuncType(
		func(
			f reflect.Type,
			t reflect.Type,
			data any,
		) (any, error) {
			// Skip other types of data
			if f.Kind() != reflect.String {
				return data, nil
			}
			if t != reflect.TypeOf(path("")) {
				return data, nil
			}

			// example
			// #config.yaml
			// key : ${WORKSPACE}/file/path
			// #.env
			// WORKSPACE = ~/user/goapps/thisapp/internal

			// viper gives us not 'path' type but regular string
			dataString := data.(string)

			// Search for '$' in string
			dollarIndex := strings.Index(dataString, "$")
			if dollarIndex == -1 {
				return data, nil
			}
			dataString = dataString[dollarIndex:]

			// Getting indexes of '{' and '}'
			openBracket := strings.Index(dataString[dollarIndex:], "{")
			closeBracket := strings.Index(dataString[dollarIndex:], "}")
			if closeBracket == -1 || openBracket == -1 {
				return data, nil
			}

			// After this operation envName == WORKSPACE
			envName := dataString[openBracket+1 : closeBracket]

			// ~/user/goapps/thisapp/internal + /file/path
			// ENV we trying to get should not contain '/'
			// and actual data we want to get should start with'/'
			// This is the most common case for any path operation
			path := viper.GetString(envName) + dataString[closeBracket+1:]
			return path, nil
		})
	return hook

}
