package githubapi

import "time"

type Tag struct {
	EdgeCursor string

	Name          string
	IsPrerelease  bool
	CommittedDate time.Time
}
