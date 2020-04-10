// Copyright 2018 Luis Guillén Civera <luisguillenc@gmail.com>. All rights reserved.

package packetproc

import (
	"time"

	"github.com/google/gopacket"
)

type (
	//CbPacket defines a callback on packet
	CbPacket func(gopacket.Packet, time.Time) (bool, Verdict, error)
	//CbTick defines callback for tick routines
	CbTick func(time.Time, time.Time) error
	//CbClose defines callback for cleanups
	CbClose func() error
)

// OnPacket stores callbacks by layer
type OnPacket struct {
	Layer    gopacket.LayerType
	Filter   Filter
	Callback CbPacket
}

// Hooks is responsible for packet processor
type Hooks struct {
	layers   []gopacket.LayerType
	onPacket map[gopacket.LayerType][]OnPacket
	sorted   []OnPacket
	onTick   []CbTick
	onClose  []CbClose
}

// NewHooks returns a new hooks collection
func NewHooks() *Hooks {
	return &Hooks{onPacket: make(map[gopacket.LayerType][]OnPacket)}
}

// OnPacket adds a callback function on new packet
func (h *Hooks) OnPacket(layer gopacket.LayerType, f Filter, fn CbPacket) {
	callbacks, ok := h.onPacket[layer]
	if !ok {
		h.layers = append(h.layers, layer)
		callbacks = make([]OnPacket, 0)
	}
	cb := OnPacket{Layer: layer, Filter: f, Callback: fn}
	callbacks = append(callbacks, cb)
	h.sorted = append(h.sorted, cb)
	h.onPacket[layer] = callbacks
}

// OnTick adds a callback function on each tick
func (h *Hooks) OnTick(fn CbTick) {
	h.onTick = append(h.onTick, fn)
}

// OnClose adds a callback function when closes source
func (h *Hooks) OnClose(fn CbClose) {
	h.onClose = append(h.onClose, fn)
}

// Layers return registered layers
func (h *Hooks) Layers() []gopacket.LayerType {
	ret := make([]gopacket.LayerType, len(h.layers), len(h.layers))
	copy(ret, h.layers)
	return ret
}

// PacketHooksByLayer returns on packet hooks by layer
func (h *Hooks) PacketHooksByLayer(layer gopacket.LayerType) []OnPacket {
	stored, ok := h.onPacket[layer]
	if !ok {
		return []OnPacket{}
	}
	ret := make([]OnPacket, len(stored), len(stored))
	copy(ret, stored)
	return ret
}

// PacketHooks returns on packet hooks in order
func (h *Hooks) PacketHooks() []OnPacket {
	ret := make([]OnPacket, len(h.sorted), len(h.sorted))
	copy(ret, h.sorted)
	return ret
}

// TickHooks returns on tick hooks
func (h *Hooks) TickHooks() []CbTick {
	ret := make([]CbTick, len(h.onTick), len(h.onTick))
	copy(ret, h.onTick)
	return ret
}

// CloseHooks returns on close hooks
func (h *Hooks) CloseHooks() []CbClose {
	ret := make([]CbClose, len(h.onClose), len(h.onClose))
	copy(ret, h.onClose)
	return ret
}
