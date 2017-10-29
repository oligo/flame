package stream

import (
	"fmt"
	"hash/fnv"
	"sync"
)

// ChanDemux is a channel demultiplexer
type ChanDemux interface {
	Execute(input <-chan Message)
	SelectChan(index int) <-chan Message
	// Fanout returns the numbers of output channels
	Fanout() int
}

type IndexedChanDemux struct {
	ChanDemux
	out   []chan Message
	index IndexFunc
}

type IndexFunc func(nchannels int, msg Message) int

func (demux *IndexedChanDemux) Execute(intput <-chan Message) {
	nchannels := len(demux.out)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for msg := range intput {
			index := demux.index(nchannels, msg)
			// Emit message
			demux.out[index] <- msg
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		for i := 0; i < nchannels; i++ {
			close(demux.out[i])
		}
	}()
}

func (demux *IndexedChanDemux) SelectChan(index int) <-chan Message {
	return demux.out[index]
}

// Fanout returns the numbers of output channels
func (demux *IndexedChanDemux) Fanout() int {
	return len(demux.out)
}

// NewIndexedChanDemux creates a IndexedChanDemux
func NewIndexedChanDemux(fanout int, index IndexFunc) ChanDemux {
	demux := &IndexedChanDemux{index: index}
	demux.out = make([]chan Message, fanout)
	for i := 0; i < fanout; i++ {
		demux.out[i] = make(chan Message)
	}
	return demux
}

// GroupDemuxIndex implements IndexedFunc
type GroupDemuxIndex struct {
	key string
}

func NewGroupDemuxIndex(key string) *GroupDemuxIndex {
	return &GroupDemuxIndex{key: key}
}

func (g *GroupDemuxIndex) GroupIndex(nchannels int, input Message) int {
	value := input.Get(g.key)
	return int(Hash(value, nchannels))
}

func Hash(value interface{}, module int) uint32 {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%s", value)))
	return h.Sum32() % uint32(module)
}
