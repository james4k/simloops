package simloops

import (
	"runtime"
	"testing"
)

type simulator interface {
	spawn(int)
	tick()
	done()
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func bench(s simulator, b *testing.B) {
	s.spawn(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//for j := 0; j < 10; j++ {
		s.tick()
		//}
	}
	b.StopTimer()
	s.done()
}
