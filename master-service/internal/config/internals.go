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

var c Config

func execute(pipeline []initPhase) error {
	for _, phase := range pipeline {
		err := phase()
		if err != nil {
			return fmt.Errorf("cannot init config %e", err)
		}
	}
	return nil
}

// Adds validation to env binding
func registerENV(input ...string) (err error) {
	err = viper.BindEnv(input...)
	if err != nil {
		return err
	}
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
			if t != reflect.TypeOf(EnvString("")) {
				return data, nil
			}
			var (
				dataString string
				ret        string

				dollar       int
				openBracket  int
				closeBracket int
			)

			dataString, _ = data.(string)

			// Search for '$' in string
			dollar = strings.IndexByte(dataString, '$')
			openBracket = strings.IndexByte(dataString, '{')
			closeBracket = strings.IndexByte(dataString, '}')

			if closeBracket == -1 || openBracket == -1 || dollar == -1 {
				return data, nil
			}
			check := strings.Index(dataString, "${")
			if check == -1 || check != dollar {
				return data, nil
			}
			if closeBracket < openBracket { // }${
				return data, nil
			}

			beforeEnv := dataString[:dollar]
			// somevalue${ENV}
			/* 			if len(dataString) <= closeBracket+1 {
				// overflow check
				return data, nil
			} */
			afterEnv := dataString[closeBracket+1:] // Dangerous

			env := dataString[openBracket+1 : closeBracket]
			ret = beforeEnv + viper.GetString(env) + afterEnv

			return ret, nil
		})
	return hook

}
