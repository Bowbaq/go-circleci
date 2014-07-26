package main

import (
	"fmt"

	"github.com/Bowbaq/go-circleci/circleci"
)

func main() {
	circleci := circleci.NewClient("<token>")

	projects, err := circleci.Projects()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, project := range projects {
		fmt.Println(project.VCSURL)
	}
}
