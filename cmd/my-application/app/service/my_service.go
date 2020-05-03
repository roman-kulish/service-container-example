package service

import (
	"errors"

	"cloud.google.com/go/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/roman-kulish/service-container-example/cmd/my-application/app/config"
	"github.com/roman-kulish/service-container-example/internal/service"
)

// MyService implements application specific internal services.
type MyService struct {
	mgo *mongo.Client
	log *zap.Logger
	bkt *storage.BucketHandle
}

// MyServiceFactory returns pre-configured MyService provider.
func MyServiceFactory(cnt *Container, cfg *config.Config) service.Provider {
	return func() error {
		if cfg.Bucket == "" {
			return errors.New("cloud storage bucket name is required")
		}
		cnt.SetMyService(&MyService{
			mgo: cnt.MongoDB(),
			log: cnt.Logger(),
			bkt: cnt.Bucket(cfg.Bucket),
		})
		return nil
	}
}