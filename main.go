package main

import (
	"context"
	"fmt"

	"github.com/integrii/flaggy"
	"github.com/nuknal/dockercat/pkg/docker"
	"github.com/nuknal/dockercat/pkg/gui"
	"github.com/nuknal/dockercat/pkg/version"
)

const (
	endpoint   = "unix:///var/run/docker.sock"
	apiVersion = "1.39"
)

func main() {
	config := docker.NewClientConfig(endpoint, "", "", "", apiVersion)
	d := docker.NewDocker(config)

	if _, err := d.Ping(context.Background()); err != nil {
		fmt.Println(err)
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
