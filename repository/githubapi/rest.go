package githubapi

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

const githubRestAPIURL = "https://api.github.com"

func (a *GithubAPI) restGet(path string, response any) error {
	req, err := http.NewRequest(http.MethodGet, githubRestAPIURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "bearer "+a.token)

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

func (a *GithubAPI) FetchFile(ref, path string) (string, error) {
	var response struct {
		Content string `json:"content"`
	}

	err := a.restGet("/repos/flutter/flutter/contents/"+path+"?ref="+ref, &response)
	if err != nil {
		return "", err
	}

	decoded, err := base64.RawStdEncoding.WithPadding('=').DecodeString(response.Content)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}
