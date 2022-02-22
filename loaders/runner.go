package loaders

import (
	log "github.com/sirupsen/logrus"
	"github.com/zetsub0u/void_archives/archive"
	"sync"
	"time"
)

type Loader interface {
	GetRecentRefs(since time.Duration) (archive.Refs, error)
}

type Runner struct {
	mutex   sync.Mutex
	archive archive.Archive
	ticker  *time.Ticker
	doneCh  chan struct{}
	stopCh  chan struct{}
	loaders []Loader
}

func NewRunner(archive archive.Archive, periodicity time.Duration, loaders ...Loader) *Runner {
	return &Runner{
		archive: archive,
		ticker:  time.NewTicker(periodicity),
		doneCh:  make(chan struct{}),
		stopCh:  make(chan struct{}),
		loaders: loaders,
	}
}

func (r *Runner) Start() {
	f := func() {
		r.mutex.Lock()
		defer r.mutex.Unlock()

		log.Infof("running loaders...")

		var loaded int

		for _, l := range r.loaders {
			refs, err := l.GetRecentRefs(5 * time.Minute)
			if err != nil {
				log.Errorf("failed getting refs from loader: %v", err)
				continue
			}
			for _, ref := range refs {
				previous, err := r.archive.GetRef(ref.URL)
				if err != nil {
					log.Errorf("failed checking ref exist (%s) on archive: %v", ref.URL, err)
					continue
				}
				// if we got something, it already was loaded, skip adding it
				if previous != nil {
					continue
				}
				if err := r.archive.AddRef(ref); err != nil {
					log.Errorf("failed adding ref (%s) to archive: %v", ref.URL, err)
					continue
				}
				loaded += 1
			}
		}
		log.Infof("done loading %d refs into the archive...", loaded)
	}

	// first tick
	f()
LOOP:
	for {
		select {
		case <-r.ticker.C:
			f()
		case <-r.stopCh:
			break LOOP
		}
	}
	r.doneCh <- struct{}{}
}

func (r *Runner) Stop() {
	log.Infof("stopping loaders...")
	close(r.stopCh)
	<-r.doneCh
	log.Infof("loaders stopped.")
}
