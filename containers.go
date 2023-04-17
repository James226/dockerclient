package dockerclient

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/james226/dockerclient/options"
)

type Container struct {
	ID string

	cli *client.Client
}

type ContainerOperations struct {
	cli *client.Client
}

func (c ContainerOperations) Start(ctx context.Context, image *Image, net *Network, opts ...*options.StartContainerOptions) (*Container, error) {
	opt := options.StartContainer()
	if len(opts) > 0 {
		opt = opts[0]
	}
	name, hasName := opt.Name()
	if hasName {
		err := removeContainer(ctx, c.cli, name, false)
		if err != nil && !client.IsErrNotFound(err) {
			return nil, err
		}
	} else {
		name = image.Name
	}
	portSet, portBindings, err := opt.Ports()
	if err != nil {
		return nil, err
	}
	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		AutoRemove:   true,
	}
	if net != nil {
		hostConfig.NetworkMode = container.NetworkMode(net.ID)
	}
	dockerPlatform := (*v1.Platform)(nil)
	platform, hasPlatform := opt.Platform()
	if hasPlatform {
		parts := strings.Split(platform, "/")
		dockerPlatform = &v1.Platform{
			OS:           parts[0],
			Architecture: parts[1],
		}
	}
	resp, err := c.cli.ContainerCreate(ctx, &container.Config{
		Image:        image.Name,
		Hostname:     name,
		ExposedPorts: portSet,
		Env:          opt.EnvironmentVariables(),
		Tty:          false,
	}, hostConfig, nil, dockerPlatform, name)
	if err != nil {
		return nil, err
	}
	err = c.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}
	containerId := resp.ID
	return &Container{
		ID:  containerId,
		cli: c.cli,
	}, nil
}

func (c *Container) Stop(ctx context.Context, logOutput bool) error {
	return stopContainer(ctx, c.cli, c.ID, logOutput)
}

func stopContainer(ctx context.Context, cli *client.Client, containerId string, logOutput bool) error {
	data, err := cli.ContainerInspect(ctx, containerId)
	if client.IsErrNotFound(err) || data.State.Status == "removing" {
		return nil
	}
	// Take logs before the container is stopped as the logs are
	// lost at that point, due to auto removal.
	if logOutput && data.State.Running {
		out, err := cli.ContainerLogs(ctx, containerId, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			log.Print(err)
		}
		stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	}
	err = cli.ContainerStop(ctx, containerId, container.StopOptions{})
	if err != nil {
		return err
	}
	statusCh, errCh := cli.ContainerWait(ctx, containerId, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Print(err)
		}
	case <-statusCh:
	}
	return nil
}

func getContainerId(ctx context.Context, cli *client.Client, containerName string) (string, error) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return "", err
	}

	dockerContainerName := fmt.Sprintf("/%s", containerName)

	for _, cont := range containers {
		for _, name := range cont.Names {
			if name == dockerContainerName {
				return cont.ID, nil
			}
		}
	}

	return "", nil
}

func removeContainer(ctx context.Context, cli *client.Client, container string, logOutput bool) error {
	containerId, err := getContainerId(ctx, cli, container)
	if err != nil {
		return err
	}

	err = stopContainer(ctx, cli, containerId, logOutput)
	if err != nil {
		return err
	}

	err = cli.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{
		RemoveVolumes: true,
	})
	if err != nil {
		return err
	}
	return nil
}
