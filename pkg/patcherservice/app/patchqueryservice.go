package app

import "time"

type PatchSpecification struct {
	PatchIDS    []PatchID
	Authors     []PatchAuthor
	Devices     []Device
	After       *time.Time
	Before      *time.Time
	ShowApplied *bool
}

type PatchQueryService interface {
	GetPatch(id PatchID) (Patch, error)
	GetPatchContent(id PatchID) (PatchContent, error)
	GetPatches(spec PatchSpecification) ([]Patch, error)
}
