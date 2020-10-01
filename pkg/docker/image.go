package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
)

// Images get images from
func (d *Docker) Images(opt types.ImageListOptions) ([]types.ImageSummary, error) {
	return d.ImageList(context.TODO(), opt)
}

// InspectImage inspect image
func (d *Docker) InspectImage(name string) (types.ImageInspect, error) {
	img, _, err := d.ImageInspectWithRaw(context.TODO(), name)
	return img, err
}

// RemoveImage remove image
func (d *Docker) RemoveImage(name string) error {
	_, err := d.ImageRemove(context.TODO(), name, types.ImageRemoveOptions{})
	return err
}

// RemoveDanglingImages remove dangling images
func (d *Docker) RemoveDanglingImages() error {
	opt := types.ImageListOptions{
		Filters: filters.NewArgs(filters.Arg("dangling", "true")),
	}

	images, err := d.Images(opt)
	if err != nil {
		return err
	}

	errIDs := []string{}

	for _, image := range images {
		if err := d.RemoveImage(image.ID); err != nil {
			errIDs = append(errIDs, image.ID[7:19])
		}
	}

	if len(errIDs) > 1 {
		return fmt.Errorf("can not remove ids\n%s", errIDs)
	}

	return nil
}

// HistoryImage show parent layers of an image.
func (d *Docker) HistoryImage(name string) ([]image.HistoryResponseItem, error) {
	his, err := d.ImageHistory(context.TODO(), name)
	return his, err
}
