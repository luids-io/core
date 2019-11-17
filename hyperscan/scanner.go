// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package hyperscan

import (
	"context"
)

// BlockScanner is the interface for check regexps in a data block
type BlockScanner interface {
	ScanBlock(ctx context.Context, block []byte) (bool, []string, error)
}

// StreamScanner is the interface for check regexps in a data stream
type StreamScanner interface {
	ScanStream(ctx context.Context, stream <-chan []byte) (<-chan string, error)
}
