package main

import (
	"time"

	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Environment struct {
	Auth struct {
		Username string `env:"USERNAME"`
		Password string `env:"PASSWORD"`
	}

	Server struct {
		Addr   string `env:"ADDR,default=:587"`
		Domain string `env:"DOMAIN,default=localhost"`

		WriteTimeout      time.Duration `env:"WRITE_TIMEOUT,default=10s"`
		ReadTimeout       time.Duration `env:"READ_TIMEOUT,default=10s"`
		MaxMessageBytes   int           `env:"MAX_MESSAGE_BYTES,default=1048576"` //1024 * 1024
		MaxRecipients     int           `env:"MAX_RECIPIENTS,default=1"`
		AllowInsecureAuth bool          `env:"ALLOW_INSECURE_AUTH,default=true"`

		LogLevel int `env:"LOG_LEVEL,default=4"` //logrus log level
	}

	Discord struct {
		WebhookURL      string `env:"DISCORD_WEBHOOK_URL"`
		MessageTemplate string `env:"DISCORD_MESSAGE_TEMPLATE"`
	}
}

type Config struct {
	*Environment
	Logger *logrus.Logger
}

var Cfg Config = Config{
	Environment: &Environment{},
	Logger:      logrus.New(),
}

func init() {

	err := godotenv.Load()
	if err == nil {
		Cfg.Logger.Info("successfully loaded .env file")
	}
	if _, err := env.UnmarshalFromEnviron(Cfg.Environment); err != nil {
		Cfg.Logger.Fatal(err)
	}
	Cfg.Logger.SetFormatter(&logrus.JSONFormatter{})
	Cfg.Logger.SetLevel(logrus.Level(Cfg.Server.LogLevel))
	//log.Println(Cfg.Logger.Level)
	//log.Println(Cfg.Server.LogLevel)
	Cfg.Logger.Info("successfully initialized configuration")
	Cfg.Logger.Debugf("config: %v", *Cfg.Environment)
	//Logger.Printf("successfully initialized configuration: %+v", Config)
}
