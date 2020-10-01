package main

import (
	"os"

	"github.com/nuknal/dockercat/pkg/docker"
	"github.com/nuknal/dockercat/pkg/gui"
	"github.com/sirupsen/logrus"
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
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		logrus.SetOutput(file)
	} else {
		logrus.Info("Failed to log to file, using default stderr")
	}

	config := docker.NewClientConfig(endpoint, "", "", "", apiVersion)
	d := docker.NewDocker(config)
	gui := gui.New(d)
	gui.Run()
}
