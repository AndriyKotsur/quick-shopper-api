package config

import (
	"bytes"
	"embed"
	"log"
	"strings"

	"github.com/spf13/viper"
)

//go:embed config.yml
var f embed.FS // in order to keep static file as part of the binary it has to be embedded

func LoadConfig() {
	viper.SetConfigType("yaml")

	data, _ := f.ReadFile("config.yml")
	err := viper.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		log.Panic(err)
	}

	// all environmental variable names are treated as uppercase despite the case in config. Always!
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
