// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package event

import (
	"context"
)

// Forwarder is the interface for event forwarding
type Forwarder interface {
	ForwardEvent(ctx context.Context, e Event) error
}
