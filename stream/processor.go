package stream

import (
	"sync"
)

type Processor struct {
	Name    string
	cfg     Config
	Handler ProcessHandler
	outChan chan Message
	demux   ChanDemux
}

type ProcessHandler interface {
	Prepare(out chan Message)
	Execute(input Message)
	Cleanup()
}

func (p *Processor) Execute(input <-chan Message) <-chan Message {
	var wg sync.WaitGroup
	numTasks := p.demux.Fanout()
	wg.Add(numTasks)

	p.Handler.Prepare()

	doWork := func(chanId int, inChan <-chan Message) {
		for msg := range inChan {
			p.Handler.Execute(msg)
		}
		wg.Done()
	}

	go func() {
		p.demux.Execute(input)
		for i := 0; i < numTasks; i++ {
			go doWork(i, p.demux.SelectChan(i))
		}
	}()

	go func() {
		wg.Wait()
		p.Handler.Cleanup()
		close(p.outChan)
	}()

	return p.outChan
}

func NewProcessor(name string, process ProcessHandler, demux ChanDemux) Processor {
	out := make(chan Message)
	return Processor{Name: name, Handler: process, outChan: out, demux: demux}
}
