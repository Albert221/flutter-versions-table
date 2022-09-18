package database

import (
	"database/sql"
	"path"
	"time"

	"github.com/pkg/errors"

	_ "github.com/mithrandie/csvq-driver"
)

type DB struct {
	db      *sql.DB
	csvName string
}

func Open(csvPath string) (*DB, error) {
	dirName := path.Dir(csvPath)
	csvName := path.Base(csvPath)

	sqlDB, err := sql.Open("csvq", dirName)
	if err != nil {
		return nil, errors.Wrap(err, "could not open database")
	}

	return &DB{
		db:      sqlDB,
		csvName: "`" + csvName + "`",
	}, nil
}

func (d *DB) FetchAll() ([]*Row, error) {
	query := `
		SELECT
			release_tag,
			release_committed_at,
			is_prerelease,
			engine_commit
		FROM ` + d.csvName + `
		ORDER BY release_committed_at DESC
	`

	dbRows, err := d.db.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "error running a query")
	}

	rows := []*Row{}
	for dbRows.Next() {
		var row Row
		var releaseCommittedAt string

		err := dbRows.Scan(
			&row.ReleaseTag,
			&releaseCommittedAt,
			&row.IsPrerelease,
			&row.EngineCommit,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error scanning query results")
		}

		row.ReleaseCommittedAt, err = time.Parse(time.RFC3339, releaseCommittedAt)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing release_committed_at")
		}

		rows = append(rows, &row)
	}

	return rows, nil
}

func (d *DB) Insert(row *Row) error {
	_, err := d.db.Exec(`
		INSERT INTO `+d.csvName+`(
			release_tag,
			release_committed_at,
			is_prerelease,
			engine_commit
		) VALUES (
			?, ?, ?, ?
		)
	`,
		row.ReleaseTag,
		row.ReleaseCommittedAt.Format(time.RFC3339),
		row.IsPrerelease,
		row.EngineCommit,
	)
	return errors.Wrap(err, "error during inserting row")
}
