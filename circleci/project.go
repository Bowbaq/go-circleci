package circleci

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Project struct {
	VCSURL   string `json:"vcs_url"`
	Followed bool
	Username string
	Reponame string
	Branches map[string]Branch
}

// Projects returns a list of Project followed by the authenticated API
// user.
func (c *Client) Projects() ([]Project, error) {
	req, err := c.NewRequest("GET", "projects", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *Project) Artifacts(c *Client, branchName string) ([]Artifact, error) {
	branch, ok := p.Branches[branchName]
	if !ok {
		return nil, errors.New("Branch not found")
	}
	build := branch.LastSuccess
	request := fmt.Sprintf("project/%s/%s/%v/artifacts", p.Username, p.Reponame, build.BuildNum)
	log.Printf("Request: %v", request)

	req, err := c.NewRequest("GET", request, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artifacts []Artifact
	if err := json.NewDecoder(resp.Body).Decode(&artifacts); err != nil {
		return nil, err
	}

	return artifacts, nil
}
