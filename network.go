package dockerclient

import (
	"context"

	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type Network struct {
	ID string
}

type NetworkOperations struct {
	cli *client.Client
}

func (n NetworkOperations) Create(ctx context.Context, name string) (*Network, error) {
	net, err := getNetwork(ctx, n.cli, name)
	if err != nil {
		return nil, err
	}

	if net != nil {
		return net, nil
	}

	newNetwork, err := n.cli.NetworkCreate(ctx, name, network.CreateOptions{
		Attachable: true,
	})
	if err != nil {
		panic(err)
	}

	return &Network{ID: newNetwork.ID}, nil
}

func getNetwork(ctx context.Context, cli *client.Client, name string) (*Network, error) {
	networks, err := cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, network := range networks {
		if network.Name == name {
			return &Network{ID: network.ID}, nil
		}
	}

	return nil, nil
}
