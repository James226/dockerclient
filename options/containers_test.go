package options

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartContainer_WhenCalled_ReturnsNewInstanceWithDefaultValues(t *testing.T) {
	opt := StartContainer()

	// Name
	v, ok := opt.Name()
	assert.Empty(t, v)
	assert.False(t, ok)

	// Ports
	set, pmap, err := opt.Ports()
	assert.Len(t, set, 0)
	assert.Len(t, pmap, 0)
	assert.Nil(t, err)

	// Environment Variables
	vars := opt.EnvironmentVariables()
	assert.Len(t, vars, 0)

	// Platform
	v, ok = opt.Platform()
	assert.Empty(t, v)
	assert.False(t, ok)
}

func TestWithName_GivenName_SetsName(t *testing.T) {
	const name = "my-container"

	opt := WithName(name)
	assert.Equal(t, name, *opt.name)
}

func TestWithPortBinding_GivenValues_SetsPortBinding(t *testing.T) {
	const host = 80
	const container = 8080
	const proto = "tcp"

	opt := WithPortBinding(host, container, proto)
	v, ok := opt.ports[host]
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf("%d/%s", container, proto), v)
}

func TestExpose_GivenPort_SetsTcpBinding(t *testing.T) {
	const port = 80

	opt := Expose(port)
	v, ok := opt.ports[port]
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf("%d/tcp", port), v)
}

func TestWithEnvironmentVariable_GivenValues_SetsEnvVar(t *testing.T) {
	const name = "my_val"
	const value = "foo"

	opt := WithEnvironmentVariable(name, value)

	v, ok := opt.environment[name]
	assert.True(t, ok)
	assert.Equal(t, v, value)
}

func TestWithEnvironmentVariables_GivenValues_SetsEnvVars(t *testing.T) {
	const name = "my_val"
	const value = "foo"

	opt := WithEnvironmentVariables(map[string]string{
		name: value,
	})

	v, ok := opt.environment[name]
	assert.True(t, ok)
	assert.Equal(t, v, value)
}

func TestStartContainerAsLinuxAmd64_WhenCalled_SetsPlatform(t *testing.T) {
	opt := StartContainer().AsLinuxAmd64()

	assert.Equal(t, "linux/amd64", *opt.platform)
}

func TestStartContainerName_WhenSet_ReturnsConfiguredName(t *testing.T) {
	name := "my-container"
	opt := &StartContainerOptions{
		name: &name,
	}

	v, ok := opt.Name()
	assert.Equal(t, name, v)
	assert.True(t, ok)
}

func TestStartContainerPorts_WherePortsAreConfigured_ReturnsPortSets(t *testing.T) {
	opt := &StartContainerOptions{
		ports: map[uint16]string{
			8080: "80/tcp",
		},
	}

	set, pmap, err := opt.Ports()
	assert.Nil(t, err)
	assert.NotNil(t, set["80/tcp"])

	bindings := pmap["80/tcp"]
	assert.Len(t, bindings, 1)
	assert.Equal(t, "8080", bindings[0].HostPort)
}

func TestStartContainerEnvironmentVariables_WhenSet_ReturnsValues(t *testing.T) {
	opt := &StartContainerOptions{
		environment: map[string]string{
			"foo": "bar",
		},
	}

	values := opt.EnvironmentVariables()
	assert.Len(t, values, 1)
	assert.Equal(t, "foo=bar", values[0])
}

func TestStartContainerPlatform_WhenSet_ReturnsConfiguredPlatform(t *testing.T) {
	platform := "linux/amd64"
	opt := &StartContainerOptions{
		platform: &platform,
	}

	v, ok := opt.Platform()
	assert.Equal(t, platform, v)
	assert.True(t, ok)
}
