package utils

import (
	"fmt"
	"strings"
)

func ParseRepo(str string) (user, repo string, err error) {
	parts := strings.Split(str, "/")

	if len(parts) != 2 {
		err = fmt.Errorf("Error: Invalid repository (e.g. octocat/hello-world). REPO: %s", str)
		return
	}

	user = parts[0]
	repo = parts[1]
	return
}

func ParseId(str, example string) (user, repo, id string, err error) {
	parts := strings.Split(str, "/")

	if len(parts) < 3 {
		err = fmt.Errorf(
			"Error: Invalid identity (e.g. octocat/hello-world/%s).",
			example,
		)
		return
	}

	user = parts[0]
	repo = parts[1]

	id = strings.Join(parts[2:], "/")

	return
}

func ParseOrgId(str, example string) (organization, id string, err error) {
	parts := strings.Split(str, "/")
	if len(parts) < 2 {
		err = fmt.Errorf(
			"Error: Invalid Organization Identity (e.g. octocat/%s)",
			example,
			)
		return
	}

	organization = parts[0]
	id = strings.Join(parts[1:], "/")

	return
}

func Bool(val bool) *bool {
	return &val
}
