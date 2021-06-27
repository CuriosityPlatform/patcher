package infrastructure

import (
	commonmysql "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"

	"patcher/pkg/patcherservice/app"
	"patcher/pkg/patcherservice/infrastructure/mysql"
	"patcher/pkg/patcherservice/infrastructure/mysql/query"
)

type DependencyContainer interface {
	PatchService() app.PatchService
	PatchQueryService() app.PatchQueryService
}

func NewDependencyContainer(client commonmysql.TransactionalClient) DependencyContainer {
	return &dependencyContainer{
		patchService:      app.NewPatchService(unitOfWorkFactory(client)),
		patchQueryService: query.NewPatchQueryService(client),
	}
}

type dependencyContainer struct {
	patchService      app.PatchService
	patchQueryService app.PatchQueryService
}

func (container *dependencyContainer) PatchService() app.PatchService {
	return container.patchService
}

func (container *dependencyContainer) PatchQueryService() app.PatchQueryService {
	return container.patchQueryService
}

func unitOfWorkFactory(client commonmysql.TransactionalClient) app.UnitOfWorkFactory {
	return mysql.NewUnitOfFactory(client)
}
