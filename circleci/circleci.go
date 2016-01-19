// Package circleci provided an API client for CircleCI. For more details, see https://circleci.com/docs/api
package circleci

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	libraryVersion = "0.1"
	defaultBaseURL = "https://circleci.com/api/v1/"
	userAgent      = "go-circleci/" + libraryVersion
)

// A Client manages communication with the CircleCI API.
type Client struct {
	// A CircleCI API token to authenticate requests
	token string

	// Base URL for API requests.  Defaults to the public CircleCI API, but can be
	// set to a domain endpoint to use with CircleCI Enterprise.  BaseURL should
	// always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the CircleCI API.
	UserAgent string
}

// Params represent additional query parameters to be passed to NewRequest
type Params map[string]string

// NewClient returns a new CircleCI API client. A token must be provided to
// authenticate API requests. Tokens can be created at https://circleci.com/account/api
func NewClient(token string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{token: token, BaseURL: baseURL, UserAgent: userAgent}
}

// NewRequest creates an API request. A relative URL can be provided in path,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
// Additional query parameters can be provided in params.
func (c *Client) NewRequest(method, path string, params *Params) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	endpoint := c.BaseURL.ResolveReference(rel)

	if params == nil {
		params = &Params{
			"circle-token": c.token,
		}
	} else {
		(*params)["circle-token"] = c.token
	}

	values := endpoint.Query()
	for key, value := range *params {
		values.Set(key, value)
	}

	endpoint.RawQuery = values.Encode()

	req, err := http.NewRequest(method, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

type Branch struct {
	PusherLogins   []string `json:"pusher_logins"`
	LastNonSuccess Build    `json:"last_non_success"`
	LastSuccess    Build    `json:"last_success"`
	RecentBuilds   []Build  `json:"recent_builds"`
	RunningBuilds  []Build  `json:"running_builds"`
}

type Build struct {
	PushedAt    time.Time `json:"pushed_at"`
	VCSRevision string    `json:"vcs_revision"`
	BuildNum    uint      `json:"build_num"`
	Outcome     string
}

type DetailedBuild struct {
	VCSURL          string `json:"vcs_url"`
	BuildURL        string `json:"build_url"`
	BuildNum        uint   `json:"build_num"`
	Branch          string
	VCSRevision     string `json:"vcs_revision"`
	CommitterName   string `json:"committer_name"`
	CommitterEmail  string `json:"committer_email"`
	Subject         string
	Body            string
	Why             string
	DontBuild       string     `json:"dont_build"`
	QueuedAt        *time.Time `json:"queued_at,omitempty"`
	StartTime       *time.Time `json:"start_time,omitempty"`
	StopTime        *time.Time `json:"stop_time,omitempty"`
	BuildTimeMillis uint       `json:"build_time_millis"`
	Username        string
	Reponame        string
	Lifecycle       string
	Outcome         string
	Status          string
	RetryOf         uint `json:"retry_of"`
	Previous        struct {
		Status   string
		BuildNum uint `json:"build_num"`
	}
	Steps []Step
}

type Step struct {
	Name    string
	Actions []Action
}

type Action struct {
	BashCommand   string     `json:"bash_command"`
	RunTimeMillis uint       `json:"run_time_millis"`
	StartTime     *time.Time `json:"start_time,omitempty"`
	EndTime       *time.Time `json:"end_time,omitempty"`
	Name          string
	Command       string
	ExitCode      int `json:"exit_code"`
	Type          string
	Index         uint
	Status        string
}

// RecentBuilds returns a list of the last 30 detailed builds across all the projects
// followed by the authenticated API user. If username and project are specified, only
// builds for that repository (eg. github/github) are returned. If branch is specified,
// only builds for that branch are returned (username and project must be specified).
// If defined, limit and offset control which page of data is returned. They must both be
// positive, and limit is capped at 100. Defaults are limit=30, offset=0
func (c *Client) RecentBuilds(username, project, branch string, limit, offset uint) ([]DetailedBuild, error) {
	var endpoint string
	if username != "" && project != "" {
		if branch != "" {
			endpoint = fmt.Sprintf("project/%s/%s/tree/%s", username, project, branch)
		} else {
			endpoint = fmt.Sprintf("project/%s/%s", username, project)
		}
	} else {
		endpoint = "recent-builds"
	}

	if limit > 100 {
		limit = 100
	}

	req, err := c.NewRequest("GET", endpoint, &Params{
		"limit":  strconv.Itoa(int(limit)),
		"offset": strconv.Itoa(int(offset)),
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var builds []DetailedBuild
	if err := json.NewDecoder(resp.Body).Decode(&builds); err != nil {
		return nil, err
	}

	return builds, nil
}

// BuildDetails returns detailed information about a specific build. This includes runtime and
// outcomes of the different steps of the build
func (c *Client) BuildDetails(username, project string, buildNum uint) (DetailedBuild, error) {
	endpoint := fmt.Sprintf("project/%s/%s/%d", username, project, buildNum)

	req, err := c.NewRequest("GET", endpoint, nil)
	if err != nil {
		return DetailedBuild{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return DetailedBuild{}, err
	}
	defer resp.Body.Close()

	var build DetailedBuild
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return DetailedBuild{}, err
	}

	return build, nil
}
