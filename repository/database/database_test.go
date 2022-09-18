package database_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Albert221/flutter-versions-table/repository/database"
	"github.com/stretchr/testify/assert"
)

func TestFetchAll(t *testing.T) {
	db, err := database.Open("testdata/fetch-all.csv")
	assert.NoError(t, err)

	rows, err := db.FetchAll()
	assert.NoError(t, err)
	assert.Len(t, rows, 1)
	assert.Equal(t, "0.0.1", rows[0].ReleaseTag)
	assert.Equal(t, time.Date(2022, 9, 9, 19, 32, 1, 0, time.UTC), rows[0].ReleaseCommittedAt)
	assert.Equal(t, true, rows[0].IsPrerelease)
	assert.Equal(t, "321xyz", rows[0].EngineCommit)
}

func TestInsertAll(t *testing.T) {
	file := "testdata/insert.csv"

	csvContent := "edge_cursor,release_tag,release_committed_at,is_prerelease,engine_commit"
	os.WriteFile(file, []byte(csvContent), os.ModePerm)
	defer os.Remove(file)

	db, err := database.Open("testdata/insert.csv")
	assert.NoError(t, err)

	row := &database.Row{
		EdgeCursor:         "TMAj",
		ReleaseTag:         "1.0.0",
		ReleaseCommittedAt: time.Date(2022, 9, 9, 19, 47, 54, 0, time.UTC),
		IsPrerelease:       false,
		EngineCommit:       "456ghj",
	}

	err = db.InsertAll([]*database.Row{row})
	assert.NoError(t, err)

	out, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	newLine := strings.Split(string(out), "\n")[1]
	assert.Equal(t, "TMAj,1.0.0,2022-09-09T19:47:54Z,false,456ghj", newLine)
}
