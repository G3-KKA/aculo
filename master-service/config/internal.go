package config

import (
	"master-service/internal/errspec"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Global config, do not try to access it.
var c Config

type (
	// Phase of config initialisation.
	initPhase func() error

	// Container for config override.
	overrideContainer struct {
		Key   string
		Value any
	}

	// Other options.
	elseSetter func() error
)

func execute(pipeline []initPhase) error {
	for _, phase := range pipeline {
		err := phase()
		if err != nil {
			return err
		}
	}

	return nil
}

// Adds validation to env binding.
func registerENV(input ...string) (err error) {
	err = viper.BindEnv(input...)
	if err != nil {
		return err
	}
	for _, env := range input {
		// Type-free validation.
		// Not defined integer or bool should be "" as well.
		envalue := viper.GetString(env)
		if envalue == "" {

			return errspec.MsgValue(ErrEnvNotDefined, "not defined", env)
		}
	}

	return nil
}

func fillGlobalConfig() error {
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// Will be called one after another.

	// Do not try to put them separately, ComposeDecode() is crucial.
	hooks := []mapstructure.DecodeHookFunc{
		envReplaceHook(),
		mapstructure.StringToTimeDurationHookFunc(),
	}
	composeHook := mapstructure.ComposeDecodeHookFunc(hooks...)
	err = viper.Unmarshal(&c, viper.DecodeHook(composeHook))
	if err != nil {
		return err
	}

	return nil
}

// Sets config file name and extension.
func handleConfigFile() error {

	viper.SetConfigFile(viper.GetString("CONFIG_FILE"))

	return nil
}

// Parse ${ENV}/dir/file kind of path,
// Only works if variable type is path, see ./config.go.
func envReplaceHook() mapstructure.DecodeHookFuncType {
	hook := mapstructure.DecodeHookFuncType(
		func(
			f reflect.Type,
			t reflect.Type,
			data any,
		) (any, error) {
			// Skip other types of data.
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

			// Search for '${...}' in string.
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
			if closeBracket < openBracket { // ...}${... check.
				return data, nil
			}

			beforeEnv := dataString[:dollar]
			afterEnv := dataString[closeBracket+1:]

			env := dataString[openBracket+1 : closeBracket]
			ret = beforeEnv + viper.GetString(env) + afterEnv

			return ret, nil
		})

	return hook
}
