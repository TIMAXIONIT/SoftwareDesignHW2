package container

import (
	"pos/internal/application"
	"pos/internal/domain"
	"pos/internal/infrastructure/persistence"
)

type Container struct {
	App *application.Facade
}

func New() *Container {
	repository := persistence.NewInMemoryRepository()

	strategy := domain.WeightedAverageStrategy{}
	calculator := domain.NewResultDomainService(strategy)

	factory := application.NewEvaluationFactory(
		repository,
		repository,
	)

	service := application.NewService(
		repository,
		repository,
		repository,
		calculator,
		factory,
	)

	facade := application.NewFacade(service)

	return &Container{
		App: facade,
	}
}
