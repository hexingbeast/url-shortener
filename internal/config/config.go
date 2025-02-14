package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)


type Config struct {
    Env string `yaml:"env" env-default:"local"`
    StoragePath string `yaml:"storage_path" env-required:"true"`
    HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
    Address string `yaml:"address" env-default:"localhost:8082"`
    Timeout time.Duration `yaml:"timeout" env-default:"4s"`
    IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
    // пользователь храниться в явном виде
    User string `yaml:"user" env-required:"true"`
    // пароль будет храниться в секретах на github
    // переменная окружения HTTP_SERVER_PASSWORD должна быть записана полностью
    // так как нету вложенности и нельзя в сруктуре Config написать HTTP_SERVER, а тут PASSWORD
    Password string `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}


// приставка Must в названии используется в тех случаях, когда функция
// вместо возврата ошибки начинает паниковать
func MustLoad() *Config {
    configPath := os.Getenv("CONFIG_PATH")
    if configPath == "" {
        log.Fatal("CONFIG_PATH is not set")
    }
    
    // check if file exists
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        log.Fatalf("config file does not exists: %s", configPath)
    }

    var cfg Config
    
    if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
        log.Fatalf("cannot read config: %s", err)
    }

    return &cfg
}
