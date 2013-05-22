package simloops

import (
	"testing"
)

type nonconcurrentSim struct {
	ents []*entity
}

func (t *nonconcurrentSim) spawn(count int) {
	for i := 0; i < count; i++ {
		t.ents = append(t.ents, newEntity())
	}
}

func (t *nonconcurrentSim) done() {
}

func (t *nonconcurrentSim) tick() {
	for _, e := range t.ents {
		e.prethink()
	}
	for _, e := range t.ents {
		e.think()
	}
	for _, e := range t.ents {
		e.postthink()
	}
}

func BenchmarkNonConcurrent(b *testing.B) {
	bench(&nonconcurrentSim{}, b)
}
