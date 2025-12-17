package main

import (
	"context"
	"fmt"
	"time"

	"github.com/integrii/flaggy"
	"github.com/nuknal/dockercat/pkg/docker"
	"github.com/nuknal/dockercat/pkg/gui"
	"github.com/nuknal/dockercat/pkg/version"
)

const (
	endpoint   = "unix:///var/run/docker.sock"
	apiVersion = "1.39"
)

var (
	flagHost       string
	flagAPIVersion string
	flagRefresh    string
)

func main() {
	h := endpoint
	if flagHost != "" {
		h = flagHost
	}
	v := apiVersion
	if flagAPIVersion != "" {
		v = flagAPIVersion
	}
	refresh := time.Second * 5
	if flagRefresh != "" {
		if d, err := time.ParseDuration(flagRefresh); err == nil {
			refresh = d
		}
	}

	config := docker.NewClientConfig(h, "", "", "", v)
	d := docker.NewDocker(config)

	if _, err := d.Ping(context.Background()); err != nil {
		fmt.Println(err)
		return
	}

	gui := gui.New(d, refresh)
	gui.Run()
}

func init() {
	flaggy.SetName("DockerCat")
	flaggy.SetDescription("A Terminal UI For Docker")
	flaggy.SetVersion(version.BuildVersion)
	flaggy.DefaultParser.AdditionalHelpPrepend = "http://github.com/nuknal/dockercat"
	flaggy.String(&flagHost, "H", "host", "docker host, e.g. unix:///var/run/docker.sock")
	flaggy.String(&flagAPIVersion, "A", "api-version", "docker api version, e.g. 1.39")
	flaggy.String(&flagRefresh, "R", "refresh", "refresh interval, e.g. 5s")
	flaggy.Parse()
}
