package simloops

import (
	"sync"
	"testing"
)

type lockstepSim struct {
	ents  []*entity
	wg    sync.WaitGroup
	mus   [2]sync.RWMutex
	mu    *sync.RWMutex
	mui   int
	state func(*entity)
}

func (t *lockstepSim) spawn(count int) {
	t.mui = 0
	t.mus[0].Lock()
	t.mus[1].Lock()
	t.mu = &t.mus[0]
	for i := 0; i < count; i++ {
		e := newEntity()
		t.ents = append(t.ents, e)
		go t.loop(e)
	}
}

func (t *lockstepSim) tick() {
	t.state = (*entity).prethink
	t.wg.Add(len(t.ents))
	t.mu.Unlock()
	t.wg.Wait()
	t.mu.Lock()
	t.mui ^= 1
	t.mu = &t.mus[t.mui]

	t.state = (*entity).think
	t.wg.Add(len(t.ents))
	t.mu.Unlock()
	t.wg.Wait()
	t.mu.Lock()
	t.mui ^= 1
	t.mu = &t.mus[t.mui]

	t.state = (*entity).postthink
	t.wg.Add(len(t.ents))
	t.mu.Unlock()
	t.wg.Wait()
	t.mu.Lock()
	t.mui ^= 1
	t.mu = &t.mus[t.mui]
}

func (t *lockstepSim) done() {
	t.state = nil
	t.wg.Add(len(t.ents))
	t.mu.Unlock()
	t.wg.Wait()
	t.mu.Lock()
	t.mui ^= 1
	t.mu = &t.mus[t.mui]
}

func (t *lockstepSim) loop(e *entity) {
	i := 0
	mu := &t.mus[i]
	for {
		mu.RLock()
		fn := t.state
		mu.RUnlock()
		if fn == nil {
			t.wg.Done()
			break
		}
		fn(e)
		t.wg.Done()
		i ^= 1
		mu = &t.mus[i]
	}
}

func BenchmarkLockstep(b *testing.B) {
	bench(&lockstepSim{}, b)
}
