package image

import (
	"context"
	"os"

	"github.com/docker/cli/cli/streams"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/dpx/dpx/pkg/core"
)

type ImageAdapter struct {
	docker *client.Client
}

func (i *ImageAdapter) Find(ctx context.Context, image string) ([]core.Image, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("reference", image)

	r, err := i.docker.ImageList(ctx, types.ImageListOptions{
		Filters: filterArgs,
	})

	if err != nil {
		return []core.Image{}, err
	}

	return mapImage(r), nil
}

func (i *ImageAdapter) Pull(ctx context.Context, image string) error {
	reader, err := i.docker.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	defer reader.Close()

	out := streams.NewOut(os.Stdout)
	jsonmessage.DisplayJSONMessagesToStream(reader, out, nil)

	return nil
}

func mapImage(images []types.ImageSummary) []core.Image {
	r := []core.Image{}

	for _, i := range images {
		r = append(r, core.Image{ID: i.ID})
	}

	return r
}

type ImageFn func(*ImageAdapter)

func New(fns ...ImageFn) *ImageAdapter {
	i := &ImageAdapter{}

	for _, fn := range fns {
		fn(i)
	}

	return i
}

func WithDocker(d *client.Client) ImageFn {
	return func(ca *ImageAdapter) {
		ca.docker = d
	}
}
