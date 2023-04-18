package main

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/james226/dockerclient"
	"github.com/james226/dockerclient/options"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"testing"
)

type IntegrationTestSuite struct {
	suite.Suite
	context   context.Context
	cli       *client.Client
	client    *dockerclient.DockerClient
	container *dockerclient.Container
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &IntegrationTestSuite{})
}

func (s *IntegrationTestSuite) SetupSuite() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	c, err := dockerclient.NewClient()
	if err != nil {
		panic(err)
	}

	network, err := c.Networks.Create(ctx, "simple-api-integration")
	if err != nil {
		panic(err)
	}

	image, err := c.Images.Build(ctx, "simple-api-container", "./")
	if err != nil {
		panic(err)
	}

	opt := options.
		WithName("simple-api-container").
		WithPortBinding(10000, 10000, "tcp")
	container, err := c.Containers.Start(ctx, image, network, opt)
	if err != nil {
		panic(err)
	}

	s.context = ctx
	s.cli = cli
	s.client = c
	s.container = container
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.container.Stop(s.context, true)

	s.cli.Close()
	s.client.Close()
}

func (s *IntegrationTestSuite) TestEndpoint() {
	resp, err := http.Get("http://localhost:10000/test")

	s.Nil(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	s.NoError(err)
	s.Equal("Hello, world!", string(body))
}
