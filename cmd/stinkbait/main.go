package main

import (
	"github.com/opsee/stinkbait/limiter"
	"github.com/opsee/stinkbait/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

func main() {
	viper.SetEnvPrefix("stinkbait")
	viper.AutomaticEnv()

	limiter, err := limiter.New(limiter.Config{
		GeneratorBucketCapacity: int64(120),
		GeneratorBucketInterval: time.Second,
		TokenCacheSize:          1024,
		TokenBucketCapacity:     int64(10),
		TokenBucketInterval:     10 * time.Second,
		HostCacheSize:           1024,
		HostBucketCapacity:      int64(10),
		HostBucketInterval:      10 * time.Second,
	})

	if err != nil {
		log.Fatal(err, " failed to initialize limiter")
	}

	server := service.New(limiter, viper.GetStringSlice("memcached_nodes"))
	server.Start(
		viper.GetString("address"),
		viper.GetString("cert"),
		viper.GetString("cert_key"),
	)
}
