package testutils

import (
	"aculo/connector-restapi/internal/config"
	log "aculo/connector-restapi/internal/logger"
	"testing"

	"github.com/spf13/viper"
)

func DefaultSetup(t TInterface, workpacePath string) {
	viper.Set("WORKSPACE", workpacePath)
	viper.Set("CONFIG_FILE", viper.GetString("WORKSPACE")+"/config.yaml")
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	err = log.InitGlobalLogger(config.Get())
	if err != nil {
		panic(err)
	}

}

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
