package testmark

import (
	"testing"

	"github.com/spf13/viper"
)

type TInterface interface {
	T() *testing.T
}
type mark int

const (
	UNIT_TEST mark = 1 << iota
	INTEGRATION_TEST
)

// Will ignore any test, if requirements not match
func MarkAs(m mark, t TInterface) {
	switch m {
	case UNIT_TEST:
		unit(t)
	case INTEGRATION_TEST:
		// Expects INTEGRATION_TEST env to be set
		intergration(t)
	}

}
func unit(t TInterface) {
	_ = viper.BindEnv("INTEGRATION_TEST")

	if viper.Get("INTEGRATION_TEST") != nil {
		t.T().Skip("Skipping unit test")
	}
}
func intergration(t TInterface) {
	_ = viper.BindEnv("INTEGRATION_TEST")

	if viper.Get("INTEGRATION_TEST") == nil {
		t.T().Skip("Skipping integration test")
	}
}
