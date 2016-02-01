package main

import (
	"fmt"

	"os"

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
		for _, project := range projects {
			artifacts, err := project.Artifacts(circleci, "master")
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, artifact := range artifacts {
				file, err := artifact.Download(circleci, ".")
				if err != nil {
					fmt.Println(err)
					return
				}
				err = os.Remove(file)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
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
