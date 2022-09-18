package repository

import "time"

type FlutterVersion struct {
	edgeCursor string

	TagName      string
	TagURL       string
	IsPrerelease bool
	CommitedAt   time.Time

	EngineCommitHash string
	EngineCommitURL  string

	DartSDKCommitHash string
	DartSDKCommitURL  string
}
