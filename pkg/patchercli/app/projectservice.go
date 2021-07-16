package app

import "github.com/pkg/errors"

var (
	ErrNoPatchesForProject      = errors.New("no patches for project")
	ErrTooManyPatchesForProject = errors.New("too many patches for project")
)

type (
	PatchID string
)

type ApplyPatchParam struct {
	PatchID   *PatchID
	WithApply bool
}

type PushCurrentChangesParam struct {
	Message string
	NoReset bool
}

type ProjectService interface {
	InitializeProject(configsDir string) (string, error)
	ApplyPatch(param ApplyPatchParam) error
	PushCurrentChanges(param PushCurrentChangesParam) error
}
