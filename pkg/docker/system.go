package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// SystemDf requests the current data usage from the daemon
func (d *Docker) SystemDf() (types.DiskUsage, error) {
	df, err := d.DiskUsage(context.TODO())

	return df, err
}

func (d *Docker) PruneAll() error {
	if _, err := d.ContainersPrune(context.TODO(), filters.Args{}); err != nil {
		return err
	}
	if _, err := d.ImagesPrune(context.TODO(), filters.Args{}); err != nil {
		return err
	}
	if _, err := d.NetworksPrune(context.TODO(), filters.Args{}); err != nil {
		return err
	}
	if _, err := d.VolumesPrune(context.TODO(), filters.Args{}); err != nil {
		return err
	}
	return nil
}
