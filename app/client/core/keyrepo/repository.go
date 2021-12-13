package keyrepo

import (
	"fmt"
	"sync"
	"time"

	"myst/pkg/crypto"
)

type Repository struct {
	mux             sync.Mutex
	keys            map[string][]byte
	lastHealthCheck time.Time
}

func (r *Repository) Set(id string, key []byte) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.keys[id] = key
}

func (r *Repository) Key(id string) ([]byte, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	key, ok := r.keys[id]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	return key, nil
}

func (r *Repository) Delete(id string) {
	r.mux.Lock()
	defer r.mux.Unlock()

	delete(r.keys, id)
}

func (r *Repository) HealthCheck() {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.lastHealthCheck = time.Now()
}

func (r *Repository) startHealthCheck() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("health check!")
			r.mux.Lock()
			elapsed := time.Since(r.lastHealthCheck)

			if elapsed < time.Minute {
				fmt.Println("healthy...")

				r.mux.Unlock()
				continue
			}

			p := crypto.DefaultArgon2IdParams

			fmt.Println("health check failed, deleting keys...")

			for k := range r.keys {
				// overwrite data
				if b, err := crypto.GenerateRandomBytes(uint(p.KeyLength)); err == nil {
					r.keys[k] = b
				}

				delete(r.keys, k)
			}

			r.mux.Unlock()
		}
	}
}

func New() *Repository {
	r := &Repository{
		keys: map[string][]byte{},
	}

	go r.startHealthCheck()

	r.HealthCheck()
	// TODO: remove simulated health check
	//go func() {
	//	for {
	//		time.Sleep(10 * time.Second)
	//		r.HealthCheck()
	//	}
	//}()

	return r
}
