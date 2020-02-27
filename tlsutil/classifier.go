// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package tlsutil

import (
	"context"

	"github.com/luids-io/core/brain/classify"
)

// Classifier for tls protocol data
type Classifier interface {
	ClassifyConnections(context.Context, []*ConnectionData) ([]classify.Response, error)
}
