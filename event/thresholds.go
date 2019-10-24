// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package event

const thUnasigned = -1

// Threshold is used for define threshold
type Threshold struct {
	low      float32
	medium   float32
	high     float32
	critical float32
}

// NewThreshold creates a new empty threshold
func NewThreshold() Threshold {
	return Threshold{
		low:      thUnasigned,
		medium:   thUnasigned,
		high:     thUnasigned,
		critical: thUnasigned,
	}
}

// Low sets the low value for threshold
func (t *Threshold) Low(value float32) {
	t.low = value
}

// Medium sets the medium value for threshold
func (t *Threshold) Medium(value float32) {
	t.medium = value
}

// High sets the high value for threshold
func (t *Threshold) High(value float32) {
	t.high = value
}

// Critical sets the critical value for threshold
func (t *Threshold) Critical(value float32) {
	t.critical = value
}

// Level computes level from a value
func (t Threshold) Level(value float32) Level {
	if t.critical != thUnasigned && value >= t.critical {
		return Critical
	}
	if t.high != thUnasigned && value >= t.high {
		return High
	}
	if t.medium != thUnasigned && value >= t.medium {
		return Medium
	}
	if t.low != thUnasigned && value >= t.low {
		return Low
	}
	return Info
}
