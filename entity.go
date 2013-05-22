package simloops

import (
	"crypto/rand"
	"hash/crc32"
	"hash/fnv"
)

type entity struct {
	data []byte
	work []byte
}

func newEntity() *entity {
	e := &entity{
		data: make([]byte, 512),
		work: make([]byte, 96),
	}
	_, err := rand.Reader.Read(e.data)
	if err != nil {
		panic(err)
	}
	return e
}

func (e *entity) prethink() {
	m := crc32.NewIEEE()
	m.Write(e.data)
	e.work = m.Sum(e.work[:0])
}

func (e *entity) think() {
	m := crc32.NewIEEE()
	m.Write(e.data)
	e.work = m.Sum(e.work)
}

func (e *entity) postthink() {
	m := fnv.New32()
	m.Write(e.data)
	e.work = m.Sum(e.work)
}
