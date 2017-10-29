package stream

type Collector struct {
	Name    string
	cfg     Config
	Handler CollectHandler
}

//type CollectFunc func(name string, out chan Message)

type CollectHandler interface {
	Prepare(out chan Message)
	Execute()
	Cleanup()
}

func (c *Collector) Execute() <-chan Message {
	done := make(chan bool)

	out := make(chan Message)
	go func() {
		c.Handler.Prepare(out)
		c.Handler.Execute()
		done <- true
	}()

	go func() {
		<-done
		c.Handler.Cleanup()
		close(done)
		close(out)
	}()

	return out
}

func NewCollector(name string, cfg Config, collectHandler CollectHandler) Collector {
	return Collector{Name: name, cfg: cfg, Handler: collectHandler}
}
