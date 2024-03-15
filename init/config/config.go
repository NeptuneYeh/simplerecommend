package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

var MyConfig *Module

type Module struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RedisAddress        string        `mapstructure:"REDIS_ADDRESS"`
	EmailSenderName     string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress  string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

func NewModule(path string) *Module {
	// 這裡的 path "./" 代表目前工作目錄，也就是你在命令列中執行 Go 指令時所處的目錄。
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	var configModule Module
	err = viper.Unmarshal(&configModule)
	MyConfig = &configModule
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return &configModule
}
