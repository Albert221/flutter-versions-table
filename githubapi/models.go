package githubapi

import "time"

type Tag struct {
	Name          string
	IsPrerelease  bool
	CommittedDate time.Time
}
