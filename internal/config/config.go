package config

import "github.com/spf13/viper"

type Config struct {
	Mysql struct {
		DSN               string
		PoolMin           int
		PoolMax           int
		IdleMax           int
		MaxIdleTimeMinute int
		MaxLifeTimeMinute int
	}
	Redis struct {
		Host string
	}
	Migration struct {
		Source string
		DSN    string
	}
}

func New(configName string) (*Config, error) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
