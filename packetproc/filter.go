// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package packetproc

import (
	"sort"

	"github.com/google/gopacket"
)

// Filter interface is used for filters
type Filter interface {
	Match(p gopacket.Packet) (bool, error)
	Layers() []gopacket.LayerType
}

// EmptyFilter is a filter that returns always match
type EmptyFilter struct {
	Filter
}

// Match implements Filter interface
func (f EmptyFilter) Match(p gopacket.Packet) (bool, error) { return true, nil }

// Layers implements Filter.Layers
func (f EmptyFilter) Layers() []gopacket.LayerType { return []gopacket.LayerType{} }

// FilterContainer is a filter that uses its childs filter
type FilterContainer struct {
	Filter
	Filters []Filter
}

// Match implements Filter interface
func (f FilterContainer) Match(p gopacket.Packet) (bool, error) {
	for _, filter := range f.Filters {
		match, err := filter.Match(p)
		if err != nil {
			return match, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

// Layers implements Filter.Layers
func (f FilterContainer) Layers() []gopacket.LayerType {
	layers := make(map[gopacket.LayerType]bool)
	for _, child := range f.Filters {
		for _, layer := range child.Layers() {
			layers[layer] = true
		}
	}
	if len(layers) > 0 {
		ret := make([]gopacket.LayerType, 0, len(layers))
		for key := range layers {
			ret = append(ret, key)
		}
		sort.Slice(ret, func(i, j int) bool { return i < j })
		return ret
	}
	return []gopacket.LayerType{}
}
