package changelog

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

// Cf. https://raw.githubusercontent.com/olivierlacan/keep-a-changelog/main/CHANGELOG.md
const changelogTemplate = `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

[Unreleased]: {{ .RepositoryUrl }}
`

func createChangelog(repositoryUrl string) (string, error) {
	tmpl, err := template.New("changelog").Parse(changelogTemplate)
	if err != nil {
		return "", fmt.Errorf("error reading changelog template: %s", err)
	}

	var builder strings.Builder
	err = tmpl.Execute(&builder, map[string]string{
		"RepositoryUrl": repositoryUrl,
	})
	if err != nil {
		return "", fmt.Errorf("error rendering changelog template: %s", err)
	}

	return builder.String(), nil
}

func Create(name string) {
	if _, err := os.Stat(name); err == nil {
		log.Fatalf("File %s already exist", name)
	}

	repositoryUrl, err := getRepositoryUrl()
	if err != nil {
		log.Fatalf("error getting repository url: %s", err)
	}

	changelog, err := createChangelog(repositoryUrl)
	if err != nil {
		log.Fatalf("error creating changelog %s: %s", name, err)
	}

	err = os.WriteFile(name, []byte(changelog), 0o644)
	if err != nil {
		log.Fatalf("error writing changelog %s: %s", name, err)
	}

	fmt.Printf("Changelog written to %s\n", name)

	fmt.Println(changelog)
}
