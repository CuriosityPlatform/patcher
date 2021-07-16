package app

import (
	"time"

	"github.com/google/uuid"
)

type PatchID uuid.UUID
type Project string
type Message []byte
type PatchContent []byte
type PatchAuthor string
type Device string

type Patch struct {
	ID      PatchID
	Project Project
	Message Message
	Applied bool
	// empty string if you use PatchRepository.Find
	Content   PatchContent
	Author    PatchAuthor
	Device    Device
	CreatedAt *time.Time
}

type PatchRepository interface {
	Find(id PatchID) (Patch, error)
	Store(patch Patch) error
}
