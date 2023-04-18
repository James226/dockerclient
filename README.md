# dockerclient

[![GoDoc](https://godoc.org/github.com/james226/dockerclient?status.svg)](https://godoc.org/github.com/james226/dockerclient)

## Getting Started

```shell
go install github.com/james226/dockerclient
```

## Examples

Here is a basic example showing how to spin up a container using this package:

```go
package main

import (
	"context"
	"github.com/james226/dockerclient"
	"github.com/james226/dockerclient/options"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := dockerclient.NewClient()
	if err != nil {
		panic(err)
	}

	network, err := c.Networks.Create(ctx, "example-network")
	if err != nil {
		panic(err)
	}

	image, err := c.Images.Build(ctx, "example-container", "./path/to/Dockerfile")
	if err != nil {
		panic(err)
	}

	wg.Add(1)

	opt := options.WithName("example-container").
		WithPortBinding(10000, 10000, "tcp").
		WithEnvironmentVariables(map[string]string{
			"PORT": "10000",
			"FOO":  "BAR",
		})
	container, err := c.Containers.Start(ctx, image, network, opt)
	if err != nil {
		panic(err)
	}

	go func() {
		abort := make(chan os.Signal, 1)
		signal.Notify(abort, syscall.SIGINT, syscall.SIGTERM)
		<-abort

		defer wg.Done()

		container.Stop(ctx, true)
		cancel()
	}()

	wg.Wait()
}
```