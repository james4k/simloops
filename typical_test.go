package simloops

import (
	"testing"
)

type typicalSim struct {
	ents []*entity
}

func (t *typicalSim) spawn(count int) {
	for i := 0; i < count; i++ {
		t.ents = append(t.ents, newEntity())
	}
}

func (t *typicalSim) done() {
}

func (t *typicalSim) tick() {
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

func BenchmarkTypical(b *testing.B) {
	bench(&typicalSim{}, b)
}
