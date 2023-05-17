package tools

import (
	"flag"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type ViperSettings struct {
	viper  *viper.Viper
	config Config
}

func (this ViperSettings) getFileName(filePath ...string) string {
	if len(filePath) != 0 {
		return filePath[0]
	}
	result := ""
	flag.StringVar(&result, "c", "", "choose config file.")
	flag.Parse()

	if len(result) != 0 {
		log.Println("Config from command arguemnt")
		return result
	} else if configEnv := this.viper.GetString("CONFIG_PATH"); configEnv != "" {
		log.Println("Config from enviroment [CONFIG_PATH]")
		return configEnv
	}
	log.Println("Config using default path")
	if len(os.Args) == 0 || !strings.HasSuffix(os.Args[0], "test") {
		// not run go test
		return "config.yaml"
	}
	var directory string
	var err error
	if directory, err = GetWorkDirectory(); err != nil {
		return ""
	}
	return path.Join(directory, "config_localtest.yaml")
}

func (this *ViperSettings) SetConfigFile() {
	if configFileName == "" {
		configFileName = this.getFileName()
	}
	log.Printf("Config file path: %v\n", configFileName)
	this.viper.SetConfigFile(configFileName)
}

func (this *ViperSettings) SetConfigType() {
	this.viper.SetConfigType("yaml")
}

func (this *ViperSettings) AutomaticEnv() {
	this.viper.AutomaticEnv()
}

func (this *ViperSettings) SetupConfig() error {
	if err := this.viper.ReadInConfig(); err != nil {
		return err
	}
	return this.viper.Unmarshal(&this.config)
}

func (this *ViperSettings) Setup() error {
	this.SetConfigFile()
	this.SetConfigType()
	this.AutomaticEnv()
	return this.SetupConfig()
}

func (this ViperSettings) GetConfig() *Config {
	return &this.config
}

func (this ViperSettings) GetViper() *viper.Viper {
	return this.viper
}

var (
	viperMux       sync.Mutex
	_viper          *viper.Viper
	config         *Config
	configFileName string
)

func SetupViperConfig(path ...string) {
	if _viper != nil && config != nil {
		return
	}
	viperMux.Lock()
	defer viperMux.Unlock()
	if _viper != nil && config != nil {
		return
	}
	result := ViperSettings{viper: viper.New()}
	if err := result.Setup(); err != nil {
		log.Fatalln(err.Error())
	}
	_viper = result.GetViper()
	config = result.GetConfig()
	return
}

func GetConfig(path ...string) Config {
	if config == nil {
		SetupViperConfig(path...)
	}
	return *config
}

func GetViper(path ...string) *viper.Viper {
	if _viper == nil {
		SetupViperConfig(path...)
	}
	return _viper
}
