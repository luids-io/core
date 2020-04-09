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

// Hooks is responsible for packet processor
type Hooks struct {
	layers   []gopacket.LayerType
	onPacket map[gopacket.LayerType][]onPacket
	onTick   []CbTick
	onClose  []CbClose
}

type onPacket struct {
	f  Filter
	fn CbPacket
}

// NewHooks returns a new hooks collection
func NewHooks() *Hooks {
	return &Hooks{
		layers:   make([]gopacket.LayerType, 0),
		onPacket: make(map[gopacket.LayerType][]onPacket),
		onTick:   make([]CbTick, 0),
		onClose:  make([]CbClose, 0),
	}
}

// OnPacket adds a callback function on new packet
func (h *Hooks) OnPacket(layer gopacket.LayerType, f Filter, fn CbPacket) {
	callbacks, ok := h.onPacket[layer]
	if !ok {
		h.layers = append(h.layers, layer)
		callbacks = make([]onPacket, 0)
	}
	callbacks = append(callbacks, onPacket{f: f, fn: fn})
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

// HooksRunner executes Hooks
type HooksRunner struct {
	hooks *Hooks
}

// NewHooksRunner returns a HooksRunner
func NewHooksRunner(h *Hooks) *HooksRunner {
	return &HooksRunner{hooks: h}
}

// Layers returns registered layers
func (h *HooksRunner) Layers() []gopacket.LayerType {
	return h.hooks.layers
}

// Packet executes all registered onPacket hooks for the layerType
// passed in a secuencial way. If some of the hooks returns true, then
// the execution stops and returns true.
func (h *HooksRunner) Packet(layer gopacket.LayerType, packet gopacket.Packet, ts time.Time) (bool, Verdict, []error) {
	callbacks, ok := h.hooks.onPacket[layer]
	if ok {
		var v Verdict
		var stop bool
		errs := make([]error, 0, len(callbacks))
		for _, cb := range callbacks {
			var err error
			if match, _ := cb.f.Match(packet); match {
				stop, v, err = cb.fn(packet, ts)
				if err != nil {
					errs = append(errs, err)
				}
				if stop {
					return true, v, errs
				}
			}
		}
		return false, v, errs
	}
	return false, Default, nil
}

// Tick executes onTick registered hooks. It pass the last timestamp.
func (h *HooksRunner) Tick(lastTick, lastPacket time.Time) []error {
	errs := make([]error, 0, len(h.hooks.onTick))
	for _, cb := range h.hooks.onTick {
		err := cb(lastTick, lastPacket)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// Close executes on close registered hooks.
func (h *HooksRunner) Close() []error {
	errs := make([]error, 0, len(h.hooks.onClose))
	for _, cb := range h.hooks.onClose {
		err := cb()
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
