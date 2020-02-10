// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package classify

import (
	"context"
)

// Classifier defines the interface for a classification service
type Classifier interface {
	Classify(context.Context, []Request) ([]Response, error)
}

// Request stores object information to classify
type Request struct {
	ID   string
	Data interface{}
}

// Response stores classification results
type Response struct {
	ID      string
	Results []Result
	Err     error
}

// Result stores label and probability
type Result struct {
	Label string
	Prob  float32
}
