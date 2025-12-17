package gui

import (
	"context"
	"fmt"

	"github.com/nuknal/dockercat/pkg/common"
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

type diskUsage struct {
	imagesSize  string
	volumesSize string
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
		Endpoint:      g.client.DaemonHost(),
		Containers:    info.Containers,
		Images:        info.Images,
		MemTotal:      fmt.Sprintf("%dMB", info.MemTotal/1024/1024),
	}, nil
}

func (g *Gui) getDiskUsage() *diskUsage {
	du := &diskUsage{
		imagesSize:  "0 MB",
		volumesSize: "0 MB",
	}

	dus, err := g.client.SystemDf()
	if err != nil {
		g.errChan <- err
		return du
	}

	var totalSizeOfImgs int64
	for _, i := range dus.Images {
		totalSizeOfImgs += i.Size
	}
	du.imagesSize = common.ParseSizeToString(totalSizeOfImgs)

	var totalSizeOfVolumes int64
	for _, v := range dus.Volumes {
		totalSizeOfVolumes += v.UsageData.Size
	}
	du.volumesSize = common.ParseSizeToString(totalSizeOfVolumes)

	return du
}
