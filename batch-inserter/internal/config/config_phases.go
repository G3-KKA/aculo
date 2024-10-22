package config

import (
	"sync"

	"github.com/spf13/viper"
)

// Global config instance

// Initialaise config process
// Every path in service works around single env WORKSPACE
var once = sync.Once{}

func initConfig() (err error) {

	once.Do(func() {
		pipeline := []initPhase{
			setEnv,
			setFlags,
			handleConfigFile,
			bindFlags,
			fillGlobalConfig,
			setElse,
			doOverride,
		}
		// panics only here
		err = execute(pipeline)
	})
	return
}

// Set and immediately validate env variable
func setEnv() error {
	for _, env := range environment {
		err := registerENV(env)
		if err != nil {
			return err
		}
	}
	return nil
}

// Set flags
func setFlags() (err error) {
	for _, flag := range flags {
		flag()
	}
	return nil
}

// Callback on config change , aliases etc.
func setElse() (err error) {
	for _, els := range elses {
		err = els()
		if err != nil {
			return err
		}
	}
	return nil
}

// Do not use, this violates constraints
// If there any way to not override - do not override (C) Me
func doOverride() error {
	for _, over := range override {
		viper.Set(over.Name, over.Value)
	}
	return nil
}
