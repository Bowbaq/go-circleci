package circleci

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type Artifact struct {
	Path       string
	PrettyPath string `json:"pretty_path"`
	NodeIndex  int    `json:"node_index"`
	URL        string
}

func (a *Artifact) Download(c *Client, folder string) (string, error) {
	u, err := url.Parse(a.URL)
	if err != nil {
		return "", err
	}
	fileName := filepath.Join(folder, filepath.Base(u.Path))
	log.Printf("Downloading %v as %v", a.URL, fileName)

	out, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer out.Close()
	req, err := c.NewRequest("GET", a.URL, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	return fileName, nil
}
