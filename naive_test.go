package simloops

import (
	"sync"
	"testing"
)

type naiveSim struct {
	ents []*entity
}

func (t *naiveSim) spawn(count int) {
	for i := 0; i < count; i++ {
		t.ents = append(t.ents, newEntity())
	}
}

func (t *naiveSim) done() {
}

func (t *naiveSim) tick() {
	var wg sync.WaitGroup
	wg.Add(len(t.ents))
	for _, e := range t.ents {
		go func() {
			e.prethink()
			wg.Done()
		}()
	}
	wg.Wait()
	wg.Add(len(t.ents))
	for _, e := range t.ents {
		go func() {
			e.think()
			wg.Done()
		}()
	}
	wg.Wait()
	wg.Add(len(t.ents))
	for _, e := range t.ents {
		go func() {
			e.postthink()
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkNaive(b *testing.B) {
	bench(&naiveSim{}, b)
}
