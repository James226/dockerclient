package options

import (
	"fmt"

	"github.com/docker/go-connections/nat"
)

// StartContainerOptions is used to pass optional arguments when starting a container.
type StartContainerOptions struct {
	name        *string
	ports       map[uint16]string
	environment map[string]string
	platform    *string
}

// StartContainer returns a new instance of StartContainerOptions.
func StartContainer() *StartContainerOptions {
	return &StartContainerOptions{
		ports:       map[uint16]string{},
		environment: map[string]string{},
	}
}

// Name is used to retrieve the configured name for the container. If no
// name is configured, an empty string followed by a false value is returned.
func (opt *StartContainerOptions) Name() (string, bool) {
	if opt.name == nil {
		return "", false
	}
	return *opt.name, true
}

// Ports is used to retrieve the Docker port configurations for starting a container.
// An error value is returned if one of the ports configurations is invalid.
func (opt *StartContainerOptions) Ports() (nat.PortSet, map[nat.Port][]nat.PortBinding, error) {
	ports := make([]string, 0)
	for host, container := range opt.ports {
		ports = append(ports, fmt.Sprintf("%d:%s", host, container))
	}
	return nat.ParsePortSpecs(ports)
}

// EnvironmentVariables is used to retrieve the environment variables configuration
// for starting a containter. The variables are in the formar of "foo=bar".
func (opt *StartContainerOptions) EnvironmentVariables() []string {
	arr := make([]string, 0)
	for name, value := range opt.environment {
		arr = append(arr, fmt.Sprintf("%s=%s", name, value))
	}
	return arr
}

// Platform is used to retrieve the platform configuration for the container. If no
// platform is specifed, an empty string followed by a false value is returned.
func (opt *StartContainerOptions) Platform() (string, bool) {
	if opt.platform == nil {
		return "", false
	}
	return *opt.platform, true
}

// WithName is used to configure the name of the container to start.
func (opt *StartContainerOptions) WithName(name string) *StartContainerOptions {
	opt.name = &name
	return opt
}

// WithPortBinding is used to configure a port binding to the container.
func (opt *StartContainerOptions) WithPortBinding(host, container uint16, protocol string) *StartContainerOptions {
	opt.ports[host] = fmt.Sprintf("%d/%s", container, protocol)
	return opt
}

// Expose is used to expose a TCP port on the container. The specified port will
// be exposed on the container and bound to the same port on the host.
func (opt *StartContainerOptions) Expose(port uint16) *StartContainerOptions {
	return opt.WithPortBinding(port, port, "tcp")
}

// WithEnvironmentVariable is used to configure a single environment variable.
func (opt *StartContainerOptions) WithEnvironmentVariable(name, value string) *StartContainerOptions {
	opt.environment[name] = value
	return opt
}

// WithEnvironmentVariables is used to configure a collection of environment variables.
func (opt *StartContainerOptions) WithEnvironmentVariables(values map[string]string) *StartContainerOptions {
	for name, value := range values {
		opt.WithEnvironmentVariable(name, value)
	}
	return opt
}

// AsLinuxAmd64 is used to configure the containers platform as linux/amd64.
func (opt *StartContainerOptions) AsLinuxAmd64() *StartContainerOptions {
	arch := "linux/amd64"
	opt.platform = &arch
	return opt
}

// WithName is used to configure the name of the container to start.
func WithName(name string) *StartContainerOptions {
	return StartContainer().WithName(name)
}

// WithPortBinding is used to configure a port binding to the container.
func WithPortBinding(host, container uint16, protocol string) *StartContainerOptions {
	return StartContainer().WithPortBinding(host, container, protocol)
}

// Expose is used to expose a TCP port on the container. The specified port will
// be exposed on the container and bound to the same port on the host.
func Expose(port uint16) *StartContainerOptions {
	return StartContainer().Expose(port)
}

// WithEnvironmentVariable is used to configure a single environment variable.
func WithEnvironmentVariable(name, value string) *StartContainerOptions {
	return StartContainer().WithEnvironmentVariable(name, value)
}

// WithEnvironmentVariables is used to configure a collection of environment variables.
func WithEnvironmentVariables(values map[string]string) *StartContainerOptions {
	return StartContainer().WithEnvironmentVariables(values)
}
