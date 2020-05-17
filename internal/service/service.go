package service

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/api/option"
)

// MongoDBAwareContainer represents a container, which provides MongoDB client.
type MongoDBAwareContainer interface {
	Container

	SetMongoDB(*mongo.Client)
	MongoDB() *mongo.Client
}

// MongoDB returns pre-configured MongoDB service provider.
func MongoDB(cnt MongoDBAwareContainer, opts *options.ClientOptions) Provider {
	return func() error {
		client, err := mongo.NewClient(opts)
		if err != nil {
			return fmt.Errorf("mongodb service: %w", err)
		}
		cnt.SetMongoDB(client)
		cnt.RegisterOnShutdown(func() {
			_ = client.Disconnect(context.Background())
		})
		return nil
	}
}

// CloudStorageAwareContainer represents a container, which provides
// Google Cloud Storage client.
type CloudStorageAwareContainer interface {
	SetCloudStorage(*storage.Client)
	CloudStorage() *storage.Client
}

// CloudStorage returns pre-configured Google Cloud Storage service provider.
func CloudStorage(cnt CloudStorageAwareContainer, ctx context.Context, opts ...option.ClientOption) Provider {
	return func() error {
		client, err := storage.NewClient(ctx, opts...)
		if err != nil {
			return fmt.Errorf("cloud storage service: %w", err)
		}
		cnt.SetCloudStorage(client)
		return nil
	}
}

// LoggerAwareContainer represents a container, which provides logger.
type LoggerAwareContainer interface {
	Container

	SetLogger(*zap.Logger)
	Logger() *zap.Logger
}

// Logger returns pre-configured Logger service provider.
func Logger(cnt LoggerAwareContainer, opts ...LoggerOption) Provider {
	cfg := newLoggerConfig()
	for _, fn := range opts {
		fn.apply(cfg)
	}
	if cfg.out == nil {
		cfg.out = os.Stdout
	}

	return func() error {
		zapConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		var encoder zapcore.Encoder
		switch cfg.format {
		case LogFormatJSON:
			zapConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
			encoder = zapcore.NewJSONEncoder(zapConfig)
		case LogFormatText, "":
			zapConfig.EncodeLevel = zapcore.CapitalLevelEncoder
			encoder = zapcore.NewConsoleEncoder(zapConfig)
		default:
			return fmt.Errorf(`unsupported log format "%s"`, cfg.format)
		}

		sync := zapcore.AddSync(cfg.out)
		zapCore := zapcore.NewCore(encoder, sync, cfg.level)
		logger := zap.New(zapCore, cfg.opts...)

		cnt.SetLogger(logger)
		cnt.RegisterOnShutdown(func() {
			_ = logger.Sync()
		})
		return nil
	}
}
