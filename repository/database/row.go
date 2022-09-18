package database

import "time"

type Row struct {
	ReleaseTag         string
	ReleaseCommittedAt time.Time
	IsPrerelease       bool
	EngineCommit       string
}
