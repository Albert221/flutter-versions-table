package database_test

import (
	"testing"

	"github.com/Albert221/flutter-versions-table/database"
	"github.com/stretchr/testify/assert"
)

func TestFetching(t *testing.T) {
	db, err := database.Open("testdata/releases.csv")
	assert.NoError(t, err)

	rows, err := db.FetchAll()
	assert.NoError(t, err)
	assert.Len(t, rows, 1)
	assert.Equal(t, "0.0.1", rows[0].ReleaseTag)
	assert.Equal(t, 9, rows[0].ReleaseCommittedAt.Day())
	assert.Equal(t, true, rows[0].IsPrerelease)
	assert.Equal(t, "321xyz", rows[0].EngineCommit)
}
