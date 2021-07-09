package app

type ProjectService interface {
	InitializeProject(configsDir string) (string, error)
	PushCurrentChanges() error
}
