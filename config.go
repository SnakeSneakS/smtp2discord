package main

import (
	"time"

	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

type Environment struct {
	Auth struct {
		Username string `env:"USERNAME"`
		Password string `env:"PASSWORD"`
	}

	Server struct {
		Port   int    `env:"PORT,default=587"`
		Domain string `env:"DOMAIN,default=localhost"`

		WriteTimeout      time.Duration `env:"WRITE_TIMEOUT,default=10s"`
		ReadTimeout       time.Duration `env:"READ_TIMEOUT,default=10s"`
		MaxMessageBytes   int           `env:"MAX_MESSAGE_BYTES,default=1048576"` //1024 * 1024
		MaxRecipients     int           `env:"MAX_RECIPIENTS,default=10"`
		AllowInsecureAuth bool          `env:"ALLOW_INSECURE_AUTH,default=true"`

		LogLevel int `eng:"LOG_LEVEL,default=5"` //logrus log level
	}

	Discord struct {
		WebhookURL      string `env:"DISCORD_WEBHOOK_URL"`
		MessageTemplate string `env:"DISCORD_MESSAGE_TEMPLATE"`
	}
}

var Config Environment

func init() {
	err := godotenv.Load()
	if err == nil {
		Logger.Printf("successfully loaded .env file: %v", err)
	}
	if _, err := env.UnmarshalFromEnviron(&Config); err != nil {
		Logger.Fatal(err)
	}
	Logger.Printf("successfully initialized configuration")
	//Logger.Printf("successfully initialized configuration: %+v", Config)
}
