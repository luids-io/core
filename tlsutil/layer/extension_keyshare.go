// Copyright 2018 Luis Guillén Civera <luisguillenc@gmail.com>. All rights reserved.

package layer

import "fmt"

// KeyShareEntry is a struct that stores a key
type KeyShareEntry struct {
	Group SupportedGroup `json:"group"`
	Key   []byte         `json:"key,omitempty"`
}

// PSKKeyExchangeMode represents..
type PSKKeyExchangeMode uint8

func (e PSKKeyExchangeMode) getDesc() string {
	if e.IsGREASE() {
		return "GREASE"
	}
	n := uint8(e)
	if n == 0 {
		return "psk_ke"
	} else if n == 1 {
		return "psk_dhe_ke"
	} else if n >= 2 && n <= 253 {
		return "unassigned"
	} else if n > 254 {
		return "reserved_private"
	}

	return "unknown"
}

func (e PSKKeyExchangeMode) String() string {
	return fmt.Sprintf("%s(%d)", e.getDesc(), e)
}

// IsGREASE returns true if is a grease value
func (e PSKKeyExchangeMode) IsGREASE() bool {
	return isGREASE8(uint8(e))
}
