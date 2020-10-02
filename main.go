package main

import (
	"fmt"

	"github.com/integrii/flaggy"
	"github.com/nuknal/dockercat/pkg/docker"
	"github.com/nuknal/dockercat/pkg/gui"
	"github.com/nuknal/dockercat/pkg/version"
	_ "github.com/nuknal/dockercat/pkg/version"
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
	config := docker.NewClientConfig(endpoint, "", "", "", apiVersion)
	d := docker.NewDocker(config)

	if err := d.ErrorConnectionFailed(); err != nil {
		fmt.Println(err.Error())
		return
	}

	gui := gui.New(d)
	gui.Run()
}

func init() {
	flaggy.SetName("DockerCat")
	flaggy.SetDescription("A Terminal UI For Docker")
	flaggy.SetVersion(version.BuildVersion)
	flaggy.DefaultParser.AdditionalHelpPrepend = "http://github.com/nuknal/dockercat"
	flaggy.Parse()
}
