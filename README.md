## Service Container 

This is an example of a service container implementation in Go. The goal is to build it without using any Voodoo magic, 
code generation at compile time, reflection at runtime and other "Abracadabra": simple, reliable, extensible, 
easy to understand and use.

### Components

See [internal/service](internal/service): this package declares a few interfaces and types described below and 
implements a few example services.

#### Container interface

(comments are omitted, see the source code for more details).

```go
type ShutdownFunc func()

type Container interface {
    RegisterOnShutdown(ShutdownFunc)
    Shutdown()
}
```

`Container` interface must be implemented by the application container. 
See [cmd/my-application/app/service/container.go](cmd/my-application/app/service/container.go)

To make life easier, application-space container can simply embed `service.ShutdownHanlder` type. It implements 
Container interface and its methods and so type, which embeds it will implicitly implement `Container` as well:

```go
type Container struct {
    service.ShutdownHandler
}
```

#### Provider type

```go
type Provider func() error
```

Provider represents a service provider, and it is responsible for instantiating service and setting it to the container.

Each `Provider` must be enclosed into a factory function, which must receive a container and services configuration 
options as its arguments. 

> An instance of the container and other arguments will be persisted in the `Provider` function scope.

Also, for each provider there must be an interface defining service getter and setter methods. Container, 
that is aware of this service, must implement this interface.

```go
type MongoDBAwareContainer interface {
    SetMongoDB(*mongo.Client, ShutdownFunc)
    MongoDB() *mongo.Client
}

func MongoDB(cnt MongoDBAwareContainer, opts *options.ClientOptions) Provider {
    return func() error {
        client, err := mongo.NewClient(opts)
        if err != nil {
            return fmt.Errorf("mongodb service: %w", err)
        }
        cnt.SetMongoDB(client, func() {
            _ = client.Disconnect(context.Background())
        })
        return nil
    }
}
```

`MongoDBAwareContainer` interface sets up a contract between the container and the service provider:

* `Provider` function returned by the `MongoDB()` will instantiate `mongo.Client` and set it to the container 
together with the shutdown function via `SetMongoDB(*mongo.Client, ShutdownFunc)` method, which container must 
implement.
* `MongoDB() *mongo.Client` can be used to retrieve an instance of `mongo.Client` from the container.

Getter method is generally required to access services on the container, as well as by service providers, which depend 
on other services.

Note that service provider factory function accepts an instance of `MongoDBAwareContainer` to enforce the contract.

See advanced example of the Logger service here: [internal/service/logger.go](internal/service/logger.go)

#### Providers using other services (dependencies)

See [internal/service/custom_service.go](internal/service/custom_service.go)

In this example, `CustomService` depends on logger, MongoDB client and Google Cloud Storage:

```go
type CustomService struct {
    mgo *mongo.Client
    log *zap.Logger
    bkt *storage.BucketHandle
}
```

The trick is to specify an interface for it, which embeds dependency interfaces:

```go
type CustomServiceAwareContainer interface {
    MongoDBAwareContainer
    LoggerAwareContainer
    CloudStorageAwareContainer

    SetCustomService(*CustomService)
    CustomService() *CustomService
}
```

In this case, `CustomService` provider will know how to retrieve dependencies from the Container.

### Wiring services to the container

`service.Wire(...)` function is responsible to executing service providers. See usage example here: 
[cmd/my-application/app/service/container.go](cmd/my-application/app/service/container.go)

```go
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
```

Note that the order of the service providers passed into the `service.Wire(...)` function matters. In this case, logger, 
MongoDB and Cloud Storage clients must be instantiated before `CustomService` and `MyService`.