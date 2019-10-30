// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package event

import (
	"context"
)

// Archiver is the interface for the event archive
type Archiver interface {
	SaveEvent(ctx context.Context, e Event) (string, error)
}
