package service

import (
	"errors"

	"cloud.google.com/go/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// CustomService is an example of a service, which depends on other services.
type CustomService struct {
	mgo *mongo.Client
	log *zap.Logger
	bkt *storage.BucketHandle
}

// CustomServiceAwareContainer represents a container, which provides CustomService.
type CustomServiceAwareContainer interface {
	MongoDBAwareContainer
	LoggerAwareContainer
	CloudStorageAwareContainer

	SetCustomService(*CustomService)
	CustomService() *CustomService
}

// CustomServiceFactory returns pre-configured CustomService provider.
func CustomServiceFactory(cnt CustomServiceAwareContainer, bucket string) Provider {
	return func() error {
		if bucket == "" {
			return errors.New("cloud storage bucket name is required")
		}
		cnt.SetCustomService(&CustomService{
			mgo: cnt.MongoDB(),
			log: cnt.Logger(),
			bkt: cnt.CloudStorage().Bucket(bucket),
		})
		return nil
	}
}
