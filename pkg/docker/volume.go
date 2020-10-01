package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// Volumes get volumes
func (d *Docker) Volumes() ([]*types.Volume, error) {
	res, err := d.VolumeList(context.TODO(), filters.Args{})
	if err != nil {
		return nil, err
	}

	return res.Volumes, nil
}

// InspectVolume inspect volume
func (d *Docker) InspectVolume(name string) (types.Volume, error) {
	volume, _, err := d.VolumeInspectWithRaw(context.TODO(), name)
	return volume, err
}

// RemoveVolume remove volume
func (d *Docker) RemoveVolume(name string) error {
	return d.VolumeRemove(context.TODO(), name, false)
}

// PruneVolumes remove unused volume
func (d *Docker) PruneVolumes() error {
	_, err := d.VolumesPrune(context.TODO(), filters.Args{})
	return err
}
