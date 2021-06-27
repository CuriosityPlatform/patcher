package app

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	patchServiceLock = "patch-service-lock-%s"
)

var (
	ErrPatchAlreadyApplied              = errors.New("cannot apply applied patch")
	ErrPatchNotFound                    = errors.New("patch not found")
	ErrPatchCantAddPatchWitEmptyContent = errors.New("cannot add patch with empty content")
)

type PatchService interface {
	AddPatch(author, device string, content []byte) error
	ApplyPatch(id PatchID) error
}

func NewPatchService(unitOfWorkFactory UnitOfWorkFactory) PatchService {
	return &patchService{unitOfWorkFactory: unitOfWorkFactory}
}

type patchService struct {
	unitOfWorkFactory UnitOfWorkFactory
}

func (service *patchService) AddPatch(author, device string, content []byte) error {
	if len(content) == 0 {
		return ErrPatchCantAddPatchWitEmptyContent
	}

	return service.executeInUnitOfWork(fmt.Sprintf(patchServiceLock, author), func(provider RepositoryProvider) error {
		repo := provider.PatchRepository()
		return repo.Store(Patch{
			ID:      PatchID(uuid.New()),
			Applied: false,
			Content: content,
			Author:  PatchAuthor(author),
			Device:  Device(device),
		})
	})
}

func (service *patchService) ApplyPatch(id PatchID) error {
	return service.executeInUnitOfWork(fmt.Sprintf(patchServiceLock, uuid.UUID(id).String()), func(provider RepositoryProvider) error {
		repo := provider.PatchRepository()

		patch, err := repo.Find(id)
		if err != nil {
			return err
		}

		if patch.Applied {
			return ErrPatchAlreadyApplied
		}

		patch.Applied = true

		return repo.Store(patch)
	})
}

func (service *patchService) executeInUnitOfWork(lockName string, f func(provider RepositoryProvider) error) error {
	unitOfWork, err := service.unitOfWorkFactory.NewUnitOfWork(lockName)
	if err != nil {
		return err
	}
	defer func() {
		err = unitOfWork.Complete(err)
	}()
	err = f(unitOfWork)
	return err
}
