package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type DbConfig struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"postgres"`
	Port     string `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	Database string `yaml:"db" env:"POSTGRES_DB" env-default:"plat"`
	Username string `yaml:"user" env:"POSTGRES_USER" env-default:"platform"`
	Password string `yaml:"pass" env:"POSTGRES_PW"`
}

var instance *DbConfig
var once sync.Once

func GetDbConfig() *DbConfig {
	once.Do(
		func() {
			instance = &DbConfig{}
			if err := cleanenv.ReadConfig("./db_config.yml", instance); err != nil {
				help, _ := cleanenv.GetDescription(instance, nil)
				fmt.Println(help)
			}
		},
	)
	return instance
}
