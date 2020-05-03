package service

import (
	"sync"
)

// ShutdownFunc represents a function called on a container Shutdown().
type ShutdownFunc func()

// Container represents a service container.
type Container interface {
	// RegisterOnShutdown allows to register a function, which is executed
	// when container Shutdown() is called.
	RegisterOnShutdown(ShutdownFunc)

	// Shutdown allows to gracefully shutdown container services.
	//
	// Shutdown() must be executed before application exits and it invokes
	// shutdown functions, registered with RegisterOnShutdown(). These
	// functions are executed sequentially in the reverse order, as they
	// were registered.
	//
	// Once Shutdown() is called, container services will become unusable.
	Shutdown()
}

// Provider represents a service provider, which is responsible for
// constructing service and setting it to the container.
type Provider func() error

// ShutdownHandler provides methods for registering and calling service
// shutdown functions.
//
// ShutdownHandler implements Container interface and must be used as
// embedded object with specific container implementation.
type ShutdownHandler struct {
	mu         sync.Mutex
	onShutdown []ShutdownFunc
}

func (h *ShutdownHandler) RegisterOnShutdown(fn ShutdownFunc) {
	h.mu.Lock()
	h.onShutdown = append(h.onShutdown, fn)
	h.mu.Unlock()
}

func (h *ShutdownHandler) Shutdown() {
	h.mu.Lock()
	for i := len(h.onShutdown) - 1; i >= 0; i-- {
		h.onShutdown[i]()
	}
	// prevent panicking, if Shutdown() is called twice.
	h.onShutdown = nil
	h.mu.Unlock()
}

// BuildContainer executes given service providers and returns an error,
// if any of the providers fails.
func BuildContainer(cnt Container, sp ...Provider) error {
	for _, fn := range sp {
		if fn == nil {
			cnt.Shutdown()
			// using it wrong, fail fast at services bootstrapping.
			panic("provider function cannot be nil")
		}
		if err := fn(); err != nil {
			cnt.Shutdown()
			return err
		}
	}
	return nil
}
