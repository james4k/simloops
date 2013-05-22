package simloops

import (
	"sync"
	"testing"
)

// TODO: revise this..seems heavier on the locks than it should be

type lockstepChanSim struct {
	ents  []*entity
	chans []chan func(*entity)
	wg    sync.WaitGroup
}

func (t *lockstepChanSim) spawn(count int) {
	for i := 0; i < count; i++ {
		e := newEntity()
		c := make(chan func(*entity), 1)
		t.ents = append(t.ents, e)
		t.chans = append(t.chans, c)
		go t.loop(e, c)
	}
}

func (t *lockstepChanSim) tick() {
	t.run((*entity).prethink)
	t.run((*entity).think)
	t.run((*entity).postthink)
}

func (t *lockstepChanSim) done() {
	t.run(nil)
}

func (t *lockstepChanSim) run(fn func(*entity)) {
	t.wg.Add(len(t.ents))
	for _, c := range t.chans {
		c <- fn
	}
	t.wg.Wait()
}

func (t *lockstepChanSim) loop(e *entity, c chan func(*entity)) {
	for {
		fn := <-c
		if fn == nil {
			t.wg.Done()
			break
		}
		fn(e)
		t.wg.Done()
	}
}

func BenchmarkLockstepChan(b *testing.B) {
	bench(&lockstepChanSim{}, b)
}
