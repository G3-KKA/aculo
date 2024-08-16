package testutils

import (
	"testing"

	"github.com/spf13/viper"
)

type TInterface interface {
	T() *testing.T
}

func ThisIsUnitTest(t TInterface) {
	_ = viper.BindEnv("INTEGRATION_TEST")

	if viper.Get("INTEGRATION_TEST") != nil {
		t.T().Skip("Skipping regular test")
	}
}
func ThisIsIntegrationTest(t TInterface) {
	_ = viper.BindEnv("INTEGRATION_TEST")

	if viper.Get("INTEGRATION_TEST") == nil {
		t.T().Skip("Skipping integration test")
	}
}
