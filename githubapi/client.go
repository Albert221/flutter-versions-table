package githubapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Albert221/flutter-versions-table/utils"
)

type GithubAPI struct {
	c     http.Client
	token string
}

func NewGithubAPI(token string) *GithubAPI {
	return &GithubAPI{
		c:     http.Client{},
		token: token,
	}
}

const githubAPIURL = "https://api.github.com/graphql"

func (a *GithubAPI) query(query string, response any) error {
	requestBody := struct {
		Query string `json:"query"`
	}{Query: query}

	reqBuf := new(bytes.Buffer)
	err := json.NewEncoder(reqBuf).Encode(requestBody)

	req, err := http.NewRequest(http.MethodPost, githubAPIURL, reqBuf)
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

type Tag struct {
	Name         string
	IsPrerelease bool
}

type queryResponse[T any] struct {
	Data T `json:"data"`
}

func (a *GithubAPI) GetFlutterTags() ([]*Tag, error) {
	query := `{
  repository(name: "flutter", owner: "flutter") {
    createdAt
    refs(
      refPrefix: "refs/tags/"
      orderBy: {field: TAG_COMMIT_DATE, direction: DESC}
      first: 100
    ) {
      edges {
        node {
          name
        }
      }
    }
  }
}`

	type edge struct {
		Node struct {
			Name string `json:"name"`
		} `json:"node"`
	}
	type response struct {
		Repository struct {
			Refs struct {
				TotalCount int    `json:"totalCount"`
				Edges      []edge `json:"edges"`
			} `json:"refs"`
		} `json:"repository"`
	}

	var r queryResponse[response]
	err := a.query(query, &r)
	if err != nil {
		return nil, err
	}

	tags := utils.MapSlice(r.Data.Repository.Refs.Edges, func(edge edge) *Tag {
		name := edge.Node.Name
		return &Tag{
			Name: name,
			// https://semver.org/#spec-item-9
			IsPrerelease: strings.Contains(name, "-"),
		}
	})

	return tags, nil
}
