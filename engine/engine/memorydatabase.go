package engine

import (
	"github.com/felzix/huyilla/types"
	"strings"
	"sync"
)

type MemoryDatabase struct {
	types.Database
	sync.Mutex

	stuff map[string][]byte
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		stuff: makeStuff(),
	}
}

func (db *MemoryDatabase) Get(key string, thing types.Serializable) error {
	db.Lock()
	defer db.Unlock()
	blob, ok := db.stuff[key]
	if !ok {
		return types.NewThingNotFoundError(key)
	}

	return thing.Unmarshal(blob)
}

func (db *MemoryDatabase) GetByPrefix(prefix string) <-chan string {
	c := make(chan string)
	go func() {
		for key := range db.stuff {
			if strings.HasPrefix(key, prefix) {
				c <- key
			}
		}
		close(c)
	}()
	return c
}

func (db *MemoryDatabase) Set(key string, thing types.Serializable) error {
	db.Lock()
	defer db.Unlock()
	blob, err := thing.Marshal()
	if err != nil {
		return err
	}
	db.stuff[key] = blob
	return nil
}

func (db *MemoryDatabase) Has(key string) bool {
	db.Lock()
	defer db.Unlock()
	_, ok := db.stuff[key]
	return ok
}

func (db *MemoryDatabase) End(key string) error {
	db.Lock()
	defer db.Unlock()
	delete(db.stuff, key)
	return nil
}

func (db *MemoryDatabase) EndAll() error {
	db.Lock()
	defer db.Unlock()
	db.stuff = makeStuff()
	return nil
}

func makeStuff() map[string][]byte {
	return make(map[string][]byte, 0)
}
