package config

import (
	"github.com/spf13/viper"
)

// Initialaise config process.
// Every path in service works around single env WORKSPACE.
func initConfig() error {

	pipeline := []initPhase{
		setEnv,
		setFlags,
		handleConfigFile,
		bindFlags,
		fillGlobalConfig,
		setElse,
		doOverride,
	}
	err := execute(pipeline)

	if err != nil {
		return err
	}

	return err
}

// Set and immediately validate env variable.
func setEnv() error {
	for _, env := range environment {
		err := registerENV(env)
		if err != nil {
			return err
		}
	}

	return nil
}

// Set flags.
func setFlags() error {
	for _, flag := range flags {
		flag()
	}

	return nil
}

// Callback on config change , aliases etc.
func setElse() error {
	for _, els := range elses {

		err := els()
		if err != nil {
			return err
		}
	}

	return nil
}

// # Do not use, this violates constraints!
// If there any way to not override - do not override (C) Me.
func doOverride() error {
	for _, over := range override {
		viper.Set(over.Key, over.Value)
	}

	return nil
}
