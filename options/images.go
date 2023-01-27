package options

// BuildImageOptions is used to pass optional arguments when building an Image.
type BuildImageOptions struct {
	dockerfile string
}

// BuildImage returns a new instance of BuildImageOptions.
func BuildImage() *BuildImageOptions {
	return &BuildImageOptions{}
}

// WithDockerfile is used to specify the path to the Dockerfile.
func (opt *BuildImageOptions) WithDockerfile(dockerfile string) *BuildImageOptions {
	opt.dockerfile = dockerfile
	return opt
}

// Dockerfile is used to get the configured Dockerfile path. If no Dockerfile
// has been configured, "Dockerfile" will be returned as default.
func (opt *BuildImageOptions) Dockerfile() string {
	if opt.dockerfile == "" {
		return "Dockerfile"
	}
	return opt.dockerfile
}

// WithDockerfile returns a new instance of BuildImageOptions with the specified Dockerfile.
func WithDockerfile(dockerfile string) *BuildImageOptions {
	return BuildImage().WithDockerfile(dockerfile)
}
