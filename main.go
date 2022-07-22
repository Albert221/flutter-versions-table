package main

import (
	"fmt"
	"os"

	"github.com/Albert221/flutter-versions-table/githubapi"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")

	ghAPI := githubapi.NewGithubAPI(token)

	tags, err := ghAPI.GetFlutterTags()
	if err != nil {
		panic(err)
	}

	for _, tag := range tags {
		fmt.Printf("%v\n", tag)
	}
}
