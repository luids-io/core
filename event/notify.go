// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package event

import (
	"context"
	"errors"
)

// Notifier is the interface that the notifiers must satisfy
type Notifier interface {
	// NotifyEvent returns notification request ID
	NotifyEvent(ctx context.Context, e Event) (string, error)
}

// NotifyBuffer interface must be used for event buffering implementations
type NotifyBuffer interface {
	PushEvent(e Event) error
}

//default buffer instance
var instance NotifyBuffer

// SetBuffer sets the default buffer instance
func SetBuffer(b NotifyBuffer) {
	instance = b
}

// Notify notifies using the default buffer instance
func Notify(e Event) error {
	if instance != nil {
		return instance.PushEvent(e)
	}
	return errors.New("buffer not available")
}
