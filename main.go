package main

import (
	"github.com/nuknal/dockercat/pkg/docker"
	"github.com/nuknal/dockercat/pkg/gui"
	"github.com/nuknal/dockercat/pkg/version"
)

const (
	endpoint   = "unix:///var/run/docker.sock"
	apiVersion = "1.39"
)

type cpu struct {
	Usage cpuUsage `json:"cpu_usage"`
}

type cpuUsage struct {
	Total float64 `json:"total_usage"`
}

func main() {
	version.Show()
	config := docker.NewClientConfig(endpoint, "", "", "", apiVersion)
	d := docker.NewDocker(config)
	gui := gui.New(d)
	gui.Run()
}
