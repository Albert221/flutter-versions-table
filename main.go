package main

import (
	"os"
	"path"
	"text/template"

	"github.com/Albert221/flutter-versions-table/repository"
	"github.com/Albert221/flutter-versions-table/repository/database"
	"github.com/Albert221/flutter-versions-table/repository/githubapi"
	"github.com/pkg/errors"
)

const (
	csvDataFile       = "docs/data.csv"
	tableTemplateFile = "template/table.gohtml"
	tableOutputFile   = "docs/index.html"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")

	dbRepo, err := database.Open(csvDataFile)
	if err != nil {
		panic(err)
	}
	ghAPI := githubapi.NewGithubAPI(token)
	cachingRepo := repository.NewCaching(dbRepo, ghAPI)

	flutterVersions, err := cachingRepo.FetchAll()
	if err != nil {
		panic(err)
	}

	vm := viewModel{Versions: flutterVersions}
	err = renderView(tableTemplateFile, tableOutputFile, vm)
	if err != nil {
		panic(err)
	}
}

func renderView(templatePath, outputPath string, viewModel any) error {
	tmpl := template.Must(template.ParseFiles(templatePath))

	outputDir := path.Dir(outputPath)

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "error during creating output file parent directories")
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return errors.Wrap(err, "error during creating output file")
	}

	err = tmpl.Execute(file, viewModel)
	return errors.Wrap(err, "error during executing a template")
}

type viewModel struct {
	Versions []*repository.FlutterVersion
}
