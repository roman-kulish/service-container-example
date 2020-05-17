package service

import (
	"context"
	"os"

	"cloud.google.com/go/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/api/option"

	"github.com/roman-kulish/service-container-example/cmd/my-application/app/app"
	"github.com/roman-kulish/service-container-example/cmd/my-application/app/config"
	"github.com/roman-kulish/service-container-example/internal/service"
)

var _ service.Container = (*Container)(nil)
var _ service.MongoDBAwareContainer = (*Container)(nil)
var _ service.CloudStorageAwareContainer = (*Container)(nil)
var _ service.LoggerAwareContainer = (*Container)(nil)
var _ service.CustomServiceAwareContainer = (*Container)(nil)

type Container struct {
	service.ShutdownHandler

	cfg *config.Config
	mgo *mongo.Client
	gcs *storage.Client
	log *zap.Logger
	csm *service.CustomService
	svc *MyService
}

func NewContainer(cfg *config.Config) (*Container, error) {
	cnt := &Container{cfg: cfg}
	sp := []service.Provider{
		service.Logger(cnt, buildLoggerOptions(cfg)...),
		service.CloudStorage(cnt, context.Background(), option.WithUserAgent(app.ID)),
		service.MongoDB(cnt, cfg.MongoOptions),
		service.CustomServiceFactory(cnt, cfg.Bucket),
		MyServiceFactory(cnt, cfg),
	}
	if err := service.Wire(cnt, sp...); err != nil {
		return nil, err
	}
	return cnt, nil
}

func (c *Container) Config() *config.Config {
	return c.cfg
}

func (c *Container) SetMongoDB(client *mongo.Client) {
	c.mgo = client
}

func (c *Container) MongoDB() *mongo.Client {
	return c.mgo
}

func (c *Container) SetCloudStorage(client *storage.Client) {
	c.gcs = client
}

func (c *Container) CloudStorage() *storage.Client {
	return c.gcs
}

func (c *Container) SetLogger(log *zap.Logger) {
	c.log = log
}

func (c *Container) Logger() *zap.Logger {
	return c.log
}

func (c *Container) Bucket(name string) *storage.BucketHandle {
	return c.gcs.Bucket(name)
}

func (c *Container) SetCustomService(svc *service.CustomService) {
	c.csm = svc
}

func (c *Container) CustomService() *service.CustomService {
	return c.csm
}

func (c *Container) SetMyService(svc *MyService) {
	c.svc = svc
}

func (c *Container) MyService() *MyService {
	return c.svc
}

func buildLoggerOptions(cfg *config.Config) []service.LoggerOption {
	fields := zap.Fields(zap.String("env", string(cfg.Env)))
	errOut := zap.ErrorOutput(zapcore.AddSync(os.Stderr))
	return []service.LoggerOption{
		service.WithFormat(cfg.LogFormat),
		service.WithLevel(cfg.LogLevel),
		service.WithOptions(fields, errOut),
	}
}
