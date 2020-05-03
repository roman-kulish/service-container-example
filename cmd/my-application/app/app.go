package app

import (
	"github.com/roman-kulish/service-container-example/cmd/my-application/app/config"
	"github.com/roman-kulish/service-container-example/cmd/my-application/app/service"
)

// Run is a "real" main, which takes over bootstrapping and
// running this application, logging, graceful shutdown and so on.
func Run(_ *config.Config, cnt *service.Container) error {
	// Make sure container services are gracefully shutdown on exit.
	defer cnt.Shutdown()

	// TODO rest of initialisation

	cnt.Logger().Info("¡Hola comrades! Let's rock and roll!")
	return nil
}