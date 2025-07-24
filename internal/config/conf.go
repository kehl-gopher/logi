package config

import (
	"fmt"

	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/spf13/viper"
)

func LoadConfig(log *utils.Log) *Config {
	// viper.SetConfigFile(".env")
	// viper.SetConfigType("env")
	viper.AutomaticEnv()

	baseConfig := &BaseConfig{}
	if err := viper.Unmarshal(&baseConfig); err != nil {
		utils.PrintLog(log, fmt.Sprintf("failed to unmarshal env config: %v", err), utils.FatalLevel)
		return nil
	}

	conf := baseConfig.SetupConfig()
	return conf
}
