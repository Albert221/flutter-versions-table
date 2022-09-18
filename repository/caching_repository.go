package repository

import (
	"strings"

	"github.com/Albert221/flutter-versions-table/log"
	"github.com/Albert221/flutter-versions-table/repository/database"
	"github.com/Albert221/flutter-versions-table/repository/githubapi"
	"github.com/Albert221/flutter-versions-table/utils"
	"github.com/pkg/errors"
)

type CachingRepository struct {
	dbRepo *database.DB
	ghAPI  *githubapi.GithubAPI
	logger *log.Logger
}

func NewCaching(
	dbRepo *database.DB,
	githubAPI *githubapi.GithubAPI,
	logger *log.Logger,
) *CachingRepository {
	return &CachingRepository{
		dbRepo: dbRepo,
		ghAPI:  githubAPI,
		logger: logger,
	}
}

func (c *CachingRepository) FetchAll() ([]*FlutterVersion, error) {
	versions := []*FlutterVersion{}

	// Read database versions
	c.logger.Info("Fetching versions from local database")
	dbRows, err := c.dbRepo.FetchAll()
	if err != nil {
		return nil, errors.Wrap(err, "error during fetching versions from database")
	}

	versions = append(versions, utils.MapSlice(dbRows, dbModelToRepositoryModel)...)
	var latestTag string
	if len(versions) > 0 {
		c.logger.Sub().Info("Found %d versions", len(versions))
		latestTag = versions[0].TagName
	} else {
		c.logger.Sub().Info("Found no versions")
	}

	// Read versions from API
	afterCursor := ""
outerFor:
	for {
		c.logger.Info(`Fetching Flutter tags from GitHubAPI after cursor "%s"`, afterCursor)

		var tags []*githubapi.Tag
		tags, afterCursor, err = c.ghAPI.GetNextFlutterTags(afterCursor)
		if err != nil {
			return nil, errors.Wrap(err, "error during fetching next flutter tags")
		}

		if len(tags) == 0 {
			c.logger.Sub().Info("Found no Flutter tags")
			break
		}
		c.logger.Sub().Info("Found %d Flutter tags", len(tags))

		lastVersionIndex := len(versions)
		for _, tag := range tags {
			// If we started overlapping database tags with API ones, break.
			if latestTag == tag.Name {
				c.logger.Sub().Info("Flutter tag %s was already in local database, breaking", tag.Name)
				break outerFor
			}

			c.logger.Sub().Info("Fetching engine commit for Flutter tag %s", tag.Name)
			engineRef, err := c.ghAPI.FetchEngineCommit(tag.Name)
			if err != nil {
				return nil, errors.Wrap(err, "error fetching file with engine version")
			}
			engineRef = strings.Trim(engineRef, " \n\r")

			model := ghAPIModelsToRepositoryModel(tag, engineRef)

			versions = append(versions, model)
		}

		// Insert API models to database
		rowsToInsert := utils.MapSlice(versions[lastVersionIndex:], repositoryModelToDBModel)
		c.logger.Info("Inserting %d versions to local database", len(rowsToInsert))
		err = c.dbRepo.InsertAll(rowsToInsert)
		if err != nil {
			return nil, errors.Wrap(err, "error during inserting flutter versions into database")
		}

		if afterCursor == "" {
			break
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

func repositoryModelToDBModel(model *FlutterVersion) *database.Row {
	return &database.Row{
		ReleaseTag:         model.TagName,
		ReleaseCommittedAt: model.CommitedAt,
		IsPrerelease:       model.IsPrerelease,
		EngineCommit:       model.EngineCommitHash,
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
