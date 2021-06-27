package app

type UnitOfWorkFactory interface {
	NewUnitOfWork(lockName string) (UnitOfWork, error)
}

type RepositoryProvider interface {
	PatchRepository() PatchRepository
}

type UnitOfWork interface {
	RepositoryProvider
	Complete(err error) error
}
