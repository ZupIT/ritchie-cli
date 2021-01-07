package repoutil

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestLocalName(t *testing.T) {
	want := formula.RepoName("local-my-repo")
	got := LocalName("my-repo")
	assert.Equal(t, want, got)
}
