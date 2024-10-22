package testutils

import (
	"aculo/frontend-restapi/internal/config"
	log "aculo/frontend-restapi/internal/logger"

	"github.com/spf13/viper"
)

func DefaultPreTestSetup(workpacePath string) {
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
