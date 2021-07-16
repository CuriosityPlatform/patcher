package app

import "time"

type PatchSpec struct {
	Projects []string
	Authors  []string
	Devices  []string
}

type Patch struct {
	ID        PatchID
	Project   string
	Message   string
	Applied   bool
	Author    string
	Device    string
	CreatedAt *time.Time
}

type PatchQueryService interface {
	Query(spec PatchSpec) ([]Patch, error)
}
