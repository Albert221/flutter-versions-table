package repository

import (
	"github.com/Albert221/flutter-versions-table/repository/database"
	"github.com/Albert221/flutter-versions-table/repository/githubapi"
	"github.com/Albert221/flutter-versions-table/utils"
	"github.com/pkg/errors"
)

type CachingRepository struct {
	dbRepo *database.DB
	ghAPI  *githubapi.GithubAPI
}

func NewCaching(dbRepo *database.DB, githubAPI *githubapi.GithubAPI) *CachingRepository {
	return &CachingRepository{
		dbRepo: dbRepo,
		ghAPI:  githubAPI,
	}
}

const engineVersionFile = "bin/internal/engine.version"

func (c *CachingRepository) FetchAll() ([]*FlutterVersion, error) {
	versions := []*FlutterVersion{}

	// Read database versions
	dbRows, err := c.dbRepo.FetchAll()
	if err != nil {
		return nil, errors.Wrap(err, "error during fetching versions from database")
	}

	versions = append(versions, utils.MapSlice(dbRows, dbModelToRepositoryModel)...)
	latestTag := versions[0].TagName

	// Read versions from API
	afterCursor := ""
outerFor:
	for {
		var tags []*githubapi.Tag
		tags, afterCursor, err = c.ghAPI.GetNextFlutterTags(afterCursor)
		if err != nil {
			return nil, errors.Wrap(err, "error during fetching next flutter tags")
		}
		if len(tags) == 0 {
			break
		}

		for _, tag := range tags {
			// If we started overlapping database tags with API ones, break.
			if latestTag == tag.Name {
				break outerFor
			}

			engineRef, err := c.ghAPI.FetchFile(tag.Name, engineVersionFile)
			if err != nil {
				return nil, errors.Wrap(err, "error fetching file with engine version")
			}

			model := ghAPIModelsToRepositoryModel(tag, engineRef)

			versions = append(versions, model)
		}
	}

	return versions, nil
}

const (
	tagURLPrefix          = "https://github.com/flutter/flutter/tree/"
	engineCommitURLPrefix = "https://github.com/flutter/engine/tree/"
)

func dbModelToRepositoryModel(row *database.Row) *FlutterVersion {
	return &FlutterVersion{
		TagName:      row.ReleaseTag,
		TagURL:       tagURLPrefix + row.ReleaseTag,
		IsPrerelease: row.IsPrerelease,
		CommitedAt:   row.ReleaseCommittedAt,

		EngineCommitHash: row.EngineCommit,
		EngineCommitURL:  engineCommitURLPrefix + row.EngineCommit,
	}
}

func ghAPIModelsToRepositoryModel(tag *githubapi.Tag, engineRef string) *FlutterVersion {
	return &FlutterVersion{
		TagName:      tag.Name,
		TagURL:       tagURLPrefix + tag.Name,
		IsPrerelease: tag.IsPrerelease,
		CommitedAt:   tag.CommittedDate,

		EngineCommitHash: engineRef,
		EngineCommitURL:  engineCommitURLPrefix + engineRef,
	}
}
