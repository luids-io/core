// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package grpctls

import (
	"errors"
	"os"
	"strings"
)

// ParseURI parses uri strings used in structs, it returns protocol and address
func ParseURI(s string) (proto string, addr string, err error) {
	err = nil
	if strings.HasPrefix(s, "unix://") {
		proto = "unix"
		addr = s[7:]
	} else if strings.HasPrefix(s, "tcp://") {
		proto = "tcp"
		addr = s[6:]
	} else {
		err = errors.New("invalid prefix")
	}
	return
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
