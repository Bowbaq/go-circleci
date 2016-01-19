package circleci

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var c *Client

const MyToken = ""

func init() {
	c = NewClient(MyToken)
}
func TestArtifacts(t *testing.T) {
	projects, err := c.Projects()
	assert.NoError(t, err)
	for _, project := range projects {
		artifacts, err := project.Artifacts(c, "master")
		assert.NoError(t, err)
		for _, artifact := range artifacts {
			file, err := artifact.Download(c, ".")
			assert.NoError(t, err)
			err = os.Remove(file)
			assert.NoError(t, err)
		}
	}
}
