package main

import (
	"log"
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

		LogLevel int `eng:"LOG_LEVEL,default=5"` //logrus log level
	}

	Discord struct {
		WebhookURL      string `env:"DISCORD_WEBHOOK_URL"`
		MessageTemplate string `env:"DISCORD_MESSAGE_TEMPLATE"`
	}

	Logger *logrus.Logger
}

/*
type Config struct {
	*Environment
	Logger *logrus.Logger
}
*/

func InitLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.Level(Cfg.Server.LogLevel))
	logger.Debugf("log level: %d", logrus.GetLevel())
	logger.Info("logger initialized")
	//log.Println(logger)
	return logger
}

var Cfg Environment = Environment{}

func init() {
	logger := logrus.New()
	err := godotenv.Load()
	if err == nil {
		logger.Info("successfully loaded .env file")
	}
	if _, err := env.UnmarshalFromEnviron(&Cfg); err != nil {
		logger.Fatal(err)
	}
	log.Print(Cfg)
	//Cfg.Logger.SetFormatter(&logrus.JSONFormatter{})
	Cfg.Logger = logger
	Cfg.Logger.SetLevel(logrus.Level(Cfg.Server.LogLevel))
	log.Println(Cfg.Logger.Level)
	log.Println(Cfg.Server.LogLevel)
	Cfg.Logger.Debugf("log level: %d", Cfg.Logger.GetLevel())
	Cfg.Logger.Info("logger initialized")
	Cfg.Logger.Info("successfully initialized configuration")
	Cfg.Logger.Debugf("config: %v", Cfg)
	//Logger.Printf("successfully initialized configuration: %+v", Config)
}
