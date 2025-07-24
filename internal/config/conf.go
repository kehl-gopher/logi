package config

import (
	"fmt"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/spf13/viper"
)

func LoadConfig(log *utils.Log) *Config {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := bindEnv(viper.GetViper(), BaseConfig{}); err != nil {
		utils.PrintLog(log, fmt.Sprintf("could not bind config keys: %s", err.Error()), utils.ErrorLevel)
		return nil
	}

	baseConfig := &BaseConfig{}
	if err := viper.Unmarshal(&baseConfig); err != nil {
		utils.PrintLog(log, fmt.Sprintf("failed to unmarshal env config: %v", err), utils.FatalLevel)
		return nil
	}

	conf := baseConfig.SetupConfig()
	return conf
}

func bindEnv(v *viper.Viper, conf interface{}) error {
	mapKeys := map[string]interface{}{}
	if err := mapstructure.Decode(conf, &mapKeys); err != nil {
		return err
	}
	for k := range mapKeys {
		if err := v.BindEnv(k); err != nil {
			return err
		}
	}
	return nil
}
