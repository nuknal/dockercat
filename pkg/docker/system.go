package docker

import (
	"context"

	"github.com/docker/docker/api/types"
)

// SystemDf requests the current data usage from the daemon
func (d *Docker) SystemDf() (types.DiskUsage, error) {
	df, err := d.DiskUsage(context.TODO())

	return df, err
}
