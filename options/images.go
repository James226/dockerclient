package options

// BuildImageOptions is used to pass optional arguments when building an Image.
type BuildImageOptions struct {
	dockerfile string
	platform   *string
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

// AsLinuxAmd64 is used to specify the build architecture as linux/amd64.
func (opt *BuildImageOptions) AsLinuxAmd64() *BuildImageOptions {
	arch := "linux/amd64"
	opt.platform = &arch
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

// Platform is used to get the configured platform argument. If the platform
// has been configured, the value with true is returned, otherwise an empty string
// with false.
func (opt *BuildImageOptions) Platform() (string, bool) {
	if opt.platform == nil {
		return "", false
	}
	return *opt.platform, true
}

// WithDockerfile returns a new instance of BuildImageOptions with the specified Dockerfile.
func WithDockerfile(dockerfile string) *BuildImageOptions {
	return BuildImage().WithDockerfile(dockerfile)
}
