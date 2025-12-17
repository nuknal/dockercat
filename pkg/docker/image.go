package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

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

func (d *Docker) PullImage(ref string) error {
	rc, err := d.ImagePull(context.TODO(), ref, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer rc.Close()
	_, _ = io.Copy(io.Discard, rc)
	return nil
}

func (d *Docker) TagImage(source, target string) error {
	return d.ImageTag(context.TODO(), source, target)
}

func (d *Docker) PruneImages() error {
	_, err := d.ImagesPrune(context.TODO(), filters.Args{})
	return err
}

func (d *Docker) PushImage(ref, username, password string) error {
	authConfig := map[string]string{
		"username": username,
		"password": password,
	}
	authBytes, _ := json.Marshal(authConfig)
	auth := base64.StdEncoding.EncodeToString(authBytes)
	rc, err := d.ImagePush(context.TODO(), ref, types.ImagePushOptions{RegistryAuth: auth})
	if err != nil {
		return err
	}
	defer rc.Close()
	_, _ = io.Copy(io.Discard, rc)
	return nil
}
