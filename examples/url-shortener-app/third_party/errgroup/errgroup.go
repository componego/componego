package errgroup

import (
	"sync"
)

type Group struct {
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

func (g *Group) Wait() error {
	g.wg.Wait()
	return g.err
}

func (g *Group) Go(fn func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := fn(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
			})
		}
	}()
}
