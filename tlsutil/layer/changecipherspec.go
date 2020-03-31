// Copyright 2018 Luis Guillén Civera <luisguillenc@gmail.com>. All rights reserved.

package layer

import (
	"fmt"
)

// CipherSpecType is the value for cipherspec protocol
type CipherSpecType int8

// Valid values
const (
	CCSChange CipherSpecType = 1
)

func (c CipherSpecType) getDesc() string {
	if c == CCSChange {
		return "change_cipher_spec"
	}
	return "unknown"
}

func (c CipherSpecType) String() string {
	return fmt.Sprintf("%s(%d)", c.getDesc(), c)
}

// IsValid returns true
func (c CipherSpecType) IsValid() bool {
	return c == CCSChange
}
