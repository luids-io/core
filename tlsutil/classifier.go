// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package tlsutil

import (
	"context"
)

// Classifier for tls protocol data
type Classifier interface {
	// ClassifyConnections must return responses in the same order
	ClassifyConnections(context.Context, []*ConnectionData) ([]ClassifyResponse, error)
}

// ClassifyResponse stores classification results
type ClassifyResponse struct {
	Results []ClassifyResult
	Err     error
}

// ClassifyResult stores label and probability
type ClassifyResult struct {
	Label string
	Prob  float32
}
