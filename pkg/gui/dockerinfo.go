package gui

import (
	"context"
	"fmt"

	"github.com/nuknal/dockercat/pkg/docker"
)

type dockerInfo struct {
	HostName      string
	ServerVersion string
	APIVersion    string
	KernelVersion string
	OSType        string
	Architecture  string
	Endpoint      string
	Containers    int
	Images        int
	MemTotal      string
}

func (g *Gui) getDockerInfo() (*dockerInfo, error) {
	info, err := g.client.Info(context.TODO())
	if err != nil {
		return nil, err
	}

	var apiVersion string
	if v, err := g.client.ServerVersion(context.TODO()); err != nil {
		apiVersion = ""
	} else {
		apiVersion = v.APIVersion
	}

	return &dockerInfo{
		HostName:      info.Name,
		ServerVersion: info.ServerVersion,
		APIVersion:    apiVersion,
		KernelVersion: info.KernelVersion,
		OSType:        info.OSType,
		Architecture:  info.Architecture,
		Endpoint:      docker.Client.DaemonHost(),
		Containers:    info.Containers,
		Images:        info.Images,
		MemTotal:      fmt.Sprintf("%dMB", info.MemTotal/1024/1024),
	}, nil
}
