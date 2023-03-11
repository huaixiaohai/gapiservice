package config

import (
	"flag"

	"github.com/spf13/viper"
)

const (
	Develop    = "develop"
	Release    = "release"
	Production = "production"
)

var (
	file string
	C    *Config
)

type Config struct {
	Env string
	ServerName string
	Log           *Log
}


type Log struct {
	Level    uint32
	Output   string
	FilePath string
}


func init() {
	flag.StringVar(&file, "c", "./config/config.yaml", "config file path")
	flag.Parse()
	viper.SetConfigFile(file)
	err := viper.ReadInConfig()
	if err != nil {
		println(err)
		panic(err)
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		println(err)
		panic(err)
	}
}
