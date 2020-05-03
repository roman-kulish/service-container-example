package config

import (
	"fmt"
	"os"

	"github.com/roman-kulish/service-container-example/internal/service"
)

const (
	envEnv       = "ENV"
	envLogLevel  = "LOG_LEVEL"
	envLogFormat = "LOG_FORMAT"
	envMongoURI  = "MONGODB_URI"
	envBucket    = "STORAGE_BUCKET"
)

var envVars = []string{
	envEnv,
	envLogLevel,
	envLogFormat,
	envMongoURI,
	envBucket,
}

func NewFromEnv() (*Config, error) {
	cfg := newConfig()

	for _, key := range envVars {
		val := os.Getenv(key)
		if val == "" {
			continue
		}

		var err error
		switch key {
		case envEnv:
			cfg.Env, err = parseEnvironment(val)
		case envLogLevel:
			err = cfg.LogLevel.Set(val)
		case envLogFormat:
			cfg.LogFormat, err = parseLogFormat(val)
		case envMongoURI:
			cfg.MongoOptions.ApplyURI(val)
		case envBucket:
			cfg.Bucket = val
		}
		if err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

func parseEnvironment(val string) (Environment, error) {
	switch v := Environment(val); v {
	case EnvLocal, EnvDevelopment, EnvStaging, EnvProduction:
		return v, nil
	}
	return "", fmt.Errorf(`environment "%s" is not valid`, val)
}

func parseLogFormat(val string) (service.LogFormat, error) {
	v := service.LogFormat(val)
	if v != service.LogFormatText && v != service.LogFormatJSON {
		return "", fmt.Errorf(`log format "%s" is not supported`, val)
	}
	return v, nil
}
