package globals

import (
	"flaver/globals/tools"

	"github.com/spf13/viper"
)

func GetConfig(path ...string) tools.Config {
	return tools.GetConfig(path...)
}

func GetViper(path ...string) *viper.Viper {
	return tools.GetViper(path...)
}
