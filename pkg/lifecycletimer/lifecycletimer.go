// Package lifecycletimer provides functionalities
// to register handlers which will be executed
// asyncronously during program lifetime.
package lifecycletimer

import "time"

// Handler is a function which will be executed
// on each lifecycle elapse.
type Handler func()

// Timer manages the lifecycle ticker
// and provides functions for registering
// handler and starting/stopping the
// lifecycle ticker.
type Timer struct {
	ticker *time.Ticker

	cStop chan bool

	interval time.Duration
	handlers []Handler
}

// New initializes a new instance of
// Timer with the passed duration as
// time between each lifecycle elapse.
// This does not automatically start the
// lifecycle timer.
func New(interval time.Duration) *Timer {
	return &Timer{
		handlers: make([]Handler, 0),
		interval: interval,
	}
}

// Handle registers the passed handler so
// that it will be executed on each lifecycle
// timer elapse.
func (t *Timer) Handle(h Handler) *Timer {
	t.handlers = append(t.handlers, h)
	return t
}

// Start starts the lifecycle timer.
// The first handler execution will
// occur after the first cycle elapsed.
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

// Stop stops the lifecycle timer.
func (t *Timer) Stop() {
	go func() {
		t.cStop <- true
	}()
}
