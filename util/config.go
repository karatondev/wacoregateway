package util

import (
	"github.com/spf13/viper"
)

const ConfigName = "config"
const ConfigType = "yaml"

var Configuration Config

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	}
	Logger struct {
		Dir        string `mapstructure:"dir"`
		FileName   string `mapstructure:"file_name"`
		MaxBackups int    `mapstructure:"max_backups"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxAge     int    `mapstructure:"max_age"`
		Compress   bool   `mapstructure:"compress"`
		LocalTime  bool   `mapstructure:"local_time"`
		Level      string `mapstructure:"level"`
	} `mapstructure:"logger"`
	Postgres struct {
		Host     string   `mapstructure:"host"`
		Port     int      `mapstructure:"port"`
		Username string   `mapstructure:"username"`
		Password string   `mapstructure:"password"`
		Database string   `mapstructure:"database"`
		Options  []string `mapstructure:"options"`
	} `mapstructure:"postgres"`
	AMQP struct {
		Scheme        string `mapstructure:"scheme"`
		Host          string `mapstructure:"host"`
		Port          int    `mapstructure:"port"`
		Username      string `mapstructure:"username"`
		Password      string `mapstructure:"password"`
		Concurrent    int    `mapstructure:"concurrent"`
		PrefetchCount int    `mapstructure:"prefetch_count"`
		PrefetchSize  int    `mapstructure:"prefetch_size"`
		Global        bool   `mapstructure:"global"`
	} `mapstructure:"amqp"`
	Queues struct {
		EventHandlerQueue  string `mapstructure:"event_handler_queue"`
		MessagesEventQueue string `mapstructure:"messages_event_queue"`
		ReceiptsQueue      string `mapstructure:"receipt_event_queue"`
		QRHandlerQueue     string `mapstructure:"qr_handler_queue"`
	} `mapstructure:"queues"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (cfg *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(ConfigName)
	viper.SetConfigType(ConfigType)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	Configuration = config
	return &config, nil
}
