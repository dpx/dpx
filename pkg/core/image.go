package core

import "context"

type ImageAdapter interface {
	Find(ctx context.Context, image string) ([]Image, error)
	Pull(ctx context.Context, image string) error
}

type Image struct {
	ID   string
	Name string
}
