package config

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap/zapcore"

	"github.com/roman-kulish/service-container-example/cmd/my-application/app/app"
	"github.com/roman-kulish/service-container-example/internal/service"
)

const (
	EnvLocal       Environment = "local"
	EnvDevelopment Environment = "development"
	EnvStaging     Environment = "staging"
	EnvProduction  Environment = "production"

	defaultEnv       = EnvLocal
	defaultLogFormat = service.LogFormatText
	defaultLogLevel  = zapcore.InfoLevel
)

type Environment string

type Config struct {
	Env          Environment
	LogLevel     zapcore.Level
	LogFormat    service.LogFormat
	MongoOptions *options.ClientOptions
	Bucket       string
}

func newConfig() *Config {
	return &Config{
		Env:       defaultEnv,
		LogLevel:  defaultLogLevel,
		LogFormat: defaultLogFormat,
		MongoOptions: options.Client().SetAppName(app.ID),
	}
}
