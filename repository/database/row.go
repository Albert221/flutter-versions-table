package database

import "time"

type Row struct {
	EdgeCursor string

	ReleaseTag         string
	ReleaseCommittedAt time.Time
	IsPrerelease       bool
	EngineCommit       string
}
