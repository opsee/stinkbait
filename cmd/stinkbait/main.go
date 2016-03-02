package main

import (
	"github.com/opsee/stinkbait/service"
	"github.com/spf13/viper"
)

func main() {
	viper.SetEnvPrefix("stinkbait")
	viper.AutomaticEnv()

	server := service.New(viper.GetStringSlice("memcached_nodes"))
	server.Start(
		viper.GetString("address"),
		viper.GetString("cert"),
		viper.GetString("cert_key"),
	)
}
