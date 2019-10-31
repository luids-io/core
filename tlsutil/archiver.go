// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package tlsutil

import (
	"context"
)

// Archiver is the main interface that must be implemented by storage backends
type Archiver interface {
	SaveConnection(context.Context, *ConnectionData) (string, error)
	SaveCertificate(context.Context, *CertificateData) (string, error)
	//async write
	StoreRecord(*RecordData) error
}
