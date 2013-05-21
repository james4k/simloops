package simloops

import (
	"sync"
	"testing"
)

type batchedSim struct {
	ents []*entity
}

func (t *batchedSim) spawn(count int) {
	for i := 0; i < count; i++ {
		t.ents = append(t.ents, newEntity())
	}
}

func (t *batchedSim) done() {
}

func (t *batchedSim) tick() {
	const n = 100
	var wg sync.WaitGroup
	end := len(t.ents) - 1
	for i := 0; ; i += n {
		if i > end {
			break
		}
		ii := i + n - 1
		if ii > end {
			ii = end
		}
		batch := t.ents[i:ii]
		wg.Add(1)
		go func() {
			for _, e := range batch {
				e.prethink()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	for i := 0; ; i += n {
		if i > end {
			break
		}
		ii := i + n
		if ii > end {
			ii = end
		}
		batch := t.ents[i:ii]
		wg.Add(1)
		go func() {
			for _, e := range batch {
				e.think()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	for i := 0; ; i += n {
		if i > end {
			break
		}
		ii := i + n
		if ii > end {
			ii = end
		}
		batch := t.ents[i:ii]
		wg.Add(1)
		go func() {
			for _, e := range batch {
				e.postthink()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkBatched(b *testing.B) {
	bench(&batchedSim{}, b)
}
