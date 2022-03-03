package flashdb

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/arriqaaq/hash"
)

const (
	MinimumStartupTime = 500 * time.Millisecond
	MaximumStartupTime = 2 * MinimumStartupTime
)

// Used to put a random delay before start of each shard, so as to not
// let various shards lock at the same time
func startupDelay() time.Duration {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	d, delta := MinimumStartupTime, (MaximumStartupTime - MinimumStartupTime)
	if delta > 0 {
		d += time.Duration(rand.Int63n(int64(delta)))
	}
	return d
}

type evictor interface {
	run(cache *hash.Hash)
	stop()
}

func newSweeperWithStore(s store, sweepTime time.Duration) evictor {
	var swp = &sweeper{
		interval: sweepTime,
		stopC:    make(chan bool),
		store:    s,
	}
	runtime.SetFinalizer(swp, stopSweeper)
	return swp
}

func stopSweeper(c evictor) {
	c.stop()
}

type sweeper struct {
	store    store
	interval time.Duration
	stopC    chan bool
}

func (s *sweeper) run(cache *hash.Hash) {
	<-time.After(startupDelay())
	ticker := time.NewTicker(s.interval)
	for {
		select {
		case <-ticker.C:
			s.store.evict(cache)
		case <-s.stopC:
			ticker.Stop()
			return
		}
	}
}

func (s *sweeper) stop() {
	s.stopC <- true
}
