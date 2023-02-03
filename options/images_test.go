package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildImage_WhenCalled_ReturnsDefaultConfig(t *testing.T) {
	opt := BuildImage()

	// Dockerfile
	assert.Equal(t, "Dockerfile", opt.Dockerfile())

	// Platform
	v, ok := opt.Platform()
	assert.Empty(t, v)
	assert.False(t, ok)
}

func TestWithDockerfile_GivenNonEmptyPath_SetsDockerfile(t *testing.T) {
	const dockerfile = "app/Dockerfile"

	opt := WithDockerfile(dockerfile)

	assert.Equal(t, dockerfile, opt.Dockerfile())
}

func TestBuildImageDockerfile_WithoutSpecifyingDockerfile_ReturnsDefault(t *testing.T) {
	opt := BuildImage()
	assert.Equal(t, "Dockerfile", opt.Dockerfile())
}

func TestAsLinuxAmd64_WhenCalled_SetsPlatform(t *testing.T) {
	opt := AsLinuxAmd64()

	v, ok := opt.Platform()
	assert.Equal(t, "linux/amd64", v)
	assert.True(t, ok)
}
