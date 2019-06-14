package lifecycletimer

import "time"

type Handler func()

type Timer struct {
	ticker *time.Ticker

	cStop chan bool

	interval time.Duration
	handlers []Handler
}

func New(interval time.Duration) *Timer {
	return &Timer{
		handlers: make([]Handler, 0),
		interval: interval,
	}
}

func (t *Timer) Handle(h Handler) *Timer {
	t.handlers = append(t.handlers, h)
	return t
}

func (t *Timer) Start() *Timer {
	t.ticker = time.NewTicker(t.interval)

	go func() {
		for {
			select {

			case <-t.cStop:
				return

			case <-t.ticker.C:
				for _, h := range t.handlers {
					if h != nil {
						h()
					}
				}
			}
		}
	}()

	return t
}

func (t *Timer) Stop() {
	go func() {
		t.cStop <- true
	}()
}
