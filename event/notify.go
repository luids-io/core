// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package event

import (
	"context"
	"errors"
)

// Notifier is the interface that the notifiers must satisfy
type Notifier interface {
	// Notify returns notification request ID
	Notify(ctx context.Context, e Event) (string, error)
}

//default buffer instance
var instance *Buffer

// SetBuffer sets the default buffer instance
func SetBuffer(b *Buffer) {
	instance = b
}

// Notify notifies using the default buffer instance
func Notify(e Event) error {
	if instance != nil {
		return instance.Notify(e)
	}
	return errors.New("buffer not available")
}
