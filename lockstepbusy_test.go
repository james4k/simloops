package simloops

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type lockstepBusySim struct {
	ents   []*entity
	wg     sync.WaitGroup
	notify uint32
	state  func(*entity)
}

func (t *lockstepBusySim) spawn(count int) {
	t.wg.Add(count)
	for i := 0; i < count; i++ {
		e := newEntity()
		t.ents = append(t.ents, e)
		go t.loop(e)
	}
	t.wg.Wait()
}

func (t *lockstepBusySim) tick() {
	t.run((*entity).prethink)
	t.run((*entity).think)
	t.run((*entity).postthink)
}

func (t *lockstepBusySim) done() {
	t.run(nil)
}

func (t *lockstepBusySim) run(fn func(*entity)) {
	t.state = fn
	t.wg.Add(len(t.ents))
	atomic.AddUint32(&t.notify, 1)
	t.wg.Wait()
}

func (t *lockstepBusySim) loop(e *entity) {
	lastnotify := t.notify
	t.wg.Done()
	runtime.Gosched()
	for {
		notify := atomic.LoadUint32(&t.notify)
		backoff := 0
		for notify == lastnotify {
			if backoff < 5 {
				runtime.Gosched()
			} else {
				time.Sleep(1000)
			}
			backoff++
			notify = atomic.LoadUint32(&t.notify)
		}
		lastnotify = notify
		fn := t.state
		if fn == nil {
			t.wg.Done()
			break
		}
		fn(e)
		t.wg.Done()
	}
}

func BenchmarkLockstepBusy(b *testing.B) {
	bench(&lockstepBusySim{}, b)
}
