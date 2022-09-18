package main

import (
	"os"
	"text/template"

	"github.com/Albert221/flutter-versions-table/repository/githubapi"
	"github.com/Albert221/flutter-versions-table/utils"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")

	ghAPI := githubapi.NewGithubAPI(token)

	tags, _, err := ghAPI.GetNextFlutterTags("")
	if err != nil {
		panic(err)
	}

	tmpl := template.Must(template.ParseFiles("template/table.gohtml"))

	err = os.MkdirAll("docs", os.ModePerm)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("docs/index.html")
	if err != nil {
		panic(err)
	}

	vm := viewmodel{Tags: utils.MapSlice(tags, func(aTag *githubapi.Tag) tag {
		engineRef, err := ghAPI.FetchFile(aTag.Name, "bin/internal/engine.version")
		if err != nil {
			panic(err)
		}

		return tag{
			Name:         aTag.Name,
			IsPrerelease: aTag.IsPrerelease,
			EngineRef:    engineRef,
		}
	})}

	err = tmpl.Execute(file, vm)
	if err != nil {
		panic(err)
	}
}

type tag struct {
	Name         string
	IsPrerelease bool
	EngineRef    string
}

type viewmodel struct {
	Tags []tag
}
