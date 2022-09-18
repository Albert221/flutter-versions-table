package githubapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Albert221/flutter-versions-table/utils"
)

const githubGQLAPIURL = "https://api.github.com/graphql"

func (a *GithubAPI) gqlQuery(query string, vars map[string]any, response any) error {
	requestBody := struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables"`
	}{
		Query:     query,
		Variables: vars,
	}

	reqBuf := new(bytes.Buffer)
	err := json.NewEncoder(reqBuf).Encode(requestBody)

	req, err := http.NewRequest(http.MethodPost, githubGQLAPIURL, reqBuf)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "bearer "+a.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := a.c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	return nil
}

type gqlQueryResponse[T any] struct {
	Data T `json:"data"`
}

func (a *GithubAPI) GetNextFlutterTags(afterCursor string) (tags []*Tag, lastCursor string, err error) {
	// Limit to 10 tags to avoid hitting the rate limit with fetching files
	query := `query($afterCursor: String) {
  repository(name: "flutter", owner: "flutter") {
    refs(
      refPrefix: "refs/tags/"
      orderBy: {field: TAG_COMMIT_DATE, direction: DESC}
      first: 10
      after: $afterCursor
    ) {
      pageInfo {
        endCursor
        hasNextPage
      }
      edges {
		cursor
        node {
          name
		  target {
			... on Commit {
				committedDate
			}
		  }
        }
      }
    }
  }
}`

	type edge struct {
		Cursor string `json:"cursor"`
		Node   struct {
			Name   string `json:"name"`
			Target struct {
				CommittedDate time.Time `json:"committedDate"`
			} `json:"target"`
		} `json:"node"`
	}
	type response struct {
		Repository struct {
			Refs struct {
				PageInfo struct {
					EndCursor   string `json:"endCursor"`
					HasNextPage bool   `json:"hasNextPage"`
				} `json:"pageInfo"`
				Edges []edge `json:"edges"`
			} `json:"refs"`
		} `json:"repository"`
	}

	var r gqlQueryResponse[response]
	err = a.gqlQuery(query, map[string]any{"afterCursor": afterCursor}, &r)
	if err != nil {
		return nil, "", err
	}

	tags = utils.MapSlice(r.Data.Repository.Refs.Edges, func(edge edge) *Tag {
		node := edge.Node
		name := node.Name
		return &Tag{
			EdgeCursor: edge.Cursor,
			Name:       name,
			// https://semver.org/#spec-item-9
			IsPrerelease:  strings.Contains(name, "-"),
			CommittedDate: node.Target.CommittedDate,
		}
	})

	var nextCursor string
	pageInfo := r.Data.Repository.Refs.PageInfo
	if pageInfo.HasNextPage {
		nextCursor = pageInfo.EndCursor
	}

	return tags, nextCursor, nil
}
