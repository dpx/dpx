package mock

import (
	"context"

	"github.com/dpx/dpx/pkg/core"
)

type ImageAdapter struct {
	FindFn func(ctx context.Context, image string) ([]core.Image, error)
	PullFn func(ctx context.Context, image string) error
}

func (i *ImageAdapter) Find(ctx context.Context, image string) ([]core.Image, error) {
	return i.FindFn(ctx, image)
}

func (i *ImageAdapter) Pull(ctx context.Context, image string) error {
	return i.PullFn(ctx, image)
}

type ImageFn func(*ImageAdapter) *ImageAdapter

func NewImageAdapter(fns ...ImageFn) *ImageAdapter {
	i := &ImageAdapter{
		FindFn: func(ctx context.Context, image string) ([]core.Image, error) {
			return []core.Image{}, nil
		},
	}

	for _, f := range fns {
		f(i)
	}

	return i
}

func ImageAdapterWithFind(image string) ImageFn {
	return func(i *ImageAdapter) *ImageAdapter {
		i.FindFn = func(ctx context.Context, img string) ([]core.Image, error) {
			return []core.Image{
				{Name: image},
			}, nil
		}

		return i
	}
}

func ImageAdapterWithPull(image string) ImageFn {
	return func(i *ImageAdapter) *ImageAdapter {
		i.PullFn = func(ctx context.Context, image string) error {
			return nil
		}

		return i
	}
}
