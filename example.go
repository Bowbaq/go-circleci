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

	recent_builds, err := circleci.RecentBuilds("", "", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("All builds", len(recent_builds))

	alp_builds, err := circleci.RecentBuilds("Echo360", "alp", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ALP builds", len(alp_builds))

	master_builds, err := circleci.RecentBuilds("Echo360", "alp", "master")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Master builds", len(master_builds))

	details, err := circleci.BuildDetails("Echo360", "alp", master_builds[0].BuildNum)

	fmt.Println("Details", details)
}
