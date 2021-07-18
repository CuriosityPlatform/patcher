package infrastructure

import (
	"time"

	patcherapi "patcher/api/patcher"
	"patcher/pkg/common/infrastructure/auth"
	"patcher/pkg/patchercli/infrastructure/transport/patcheservice"
)

type DependencyContainer interface {
	PatcherClient() (patcherapi.PatcherServiceClient, error)
}

type DependenciesConfig struct {
	PatcherServiceAddress string
	SigningKey            string
}

func NewDependencyContainer(config DependenciesConfig) DependencyContainer {
	return &dependencyContainer{config: config}
}

type dependencyContainer struct {
	config DependenciesConfig
}

func (container *dependencyContainer) PatcherClient() (patcherapi.PatcherServiceClient, error) {
	return patcheservice.NewClient(
		patcheservice.ClientConfig{PatcherServiceAddress: container.config.PatcherServiceAddress},
		patcheservice.NewAuthorizationInterceptor(container.tokenManager()),
	)
}

func (container *dependencyContainer) tokenManager() auth.TokenManager {
	return auth.NewJwtTokenManager([]byte(container.config.SigningKey), time.Minute)
}
