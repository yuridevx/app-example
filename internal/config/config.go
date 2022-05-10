package config

import (
	"github.com/spf13/viper"
)

type Relic struct {
	License string
	AppName string
	Enabled bool
}

type Pyroscope struct {
	ServerAddress string
	AppPrefix     string
}

type Config struct {
	Relic         Relic
	Pyroscope     Pyroscope
	ListenGrpc    string
	ListenHttp    string
	ListenMetrics string
	Timezone      string
}

func NewConfig() *Config {
	conf := NewDefaultConfig()
	InitViper(conf)
	return conf
}

func InitViper(conf *Config) {
	viper.SetConfigName("app.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/company/")
	viper.AddConfigPath("$HOME/.company")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(conf)
	if err != nil {
		return
	}
}

func NewDefaultConfig() *Config {
	return &Config{
		ListenGrpc:    ":10000",
		ListenHttp:    ":8080",
		ListenMetrics: ":9090",
		Timezone:      "UTC",
	}
}
