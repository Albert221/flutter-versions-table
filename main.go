package main

import (
	"os"
	"text/template"

	"github.com/Albert221/flutter-versions-table/githubapi"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")

	ghAPI := githubapi.NewGithubAPI(token)

	tags, err := ghAPI.GetFlutterTags()
	if err != nil {
		panic(err)
	}

	tmpl := template.Must(template.ParseFiles("template/table.gohtml"))

	file, err := os.Create("docs/index.html")
	if err != nil {
		panic(err)
	}

	vm := viewmodel{Tags: tags}
	err = tmpl.Execute(file, vm)
	if err != nil {
		panic(err)
	}
}

type viewmodel struct {
	Tags []*githubapi.Tag
}
