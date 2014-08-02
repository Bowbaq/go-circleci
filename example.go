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

	recent_builds, err := circleci.RecentBuilds("", "", "", 30, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("All builds", len(recent_builds))

	github_builds, err := circleci.RecentBuilds("github", "github", "", 30, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GitHub builds", len(github_builds))

	master_builds, err := circleci.RecentBuilds("github", "github", "master", 30, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Master builds", len(master_builds))

	details, err := circleci.BuildDetails("github", "github", master_builds[0].BuildNum)

	fmt.Println("Details", details)
}
