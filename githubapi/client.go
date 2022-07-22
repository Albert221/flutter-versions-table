package githubapi

import (
	"net/http"
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
