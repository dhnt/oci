package registry

import (
	"sync"
)

// cidCache contains known cid entries.
type cidCache struct {
	// maps repo:tag -> cid
	cids     map[string]string
	location string

	sync.RWMutex
}

func key(repo, ref string) string {
	return repo + ":" + ref
}

func (r *cidCache) Add(repo, reference string, cid string) {
	r.Lock()

	k := key(repo, reference)

	r.cids[k] = cid

	r.Unlock()
}

func (r *cidCache) Get(repo, reference string) (string, bool) {
	r.RLock()

	k := key(repo, reference)

	val, ok := r.cids[k]

	r.RUnlock()
	return val, ok
}

func newCidCache() *cidCache {
	return &cidCache{
		cids: make(map[string]string),
	}
}
