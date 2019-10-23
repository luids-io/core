// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

// Package reason provides some utils for encoding data into response.Reason.
//
// This package is a work in progress and makes no API stability promises.
package reason

// Clean removes policy and other stuff from a reason string
func Clean(reason string) string {
	return cleanScore(cleanPolicy(reason))
}
