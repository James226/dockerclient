package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithDockerfile_GivenNonEmptyPath_SetsDockerfile(t *testing.T) {
	const dockerfile = "app/Dockerfile"

	opt := WithDockerfile(dockerfile)

	assert.Equal(t, dockerfile, opt.Dockerfile())
}

func TestBuildImageDockerfile_WithoutSpecifyingDockerfile_ReturnsDefault(t *testing.T) {
	opt := BuildImage()
	assert.Equal(t, "Dockerfile", opt.Dockerfile())
}
