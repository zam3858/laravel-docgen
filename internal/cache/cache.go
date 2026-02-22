package cache

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Entry struct {
	ModTime time.Time `json:"mod_time"`
}

type Cache struct {
	path    string
	entries map[string]Entry
	mu      sync.RWMutex
}

func New(path string) (*Cache, error) {
	c := &Cache{path: path, entries: map[string]Entry{}}

	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return c, nil
		}
		return nil, err
	}

	if len(b) == 0 {
		return c, nil
	}

	if err := json.Unmarshal(b, &c.entries); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Cache) Unchanged(path string, modTime time.Time) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[path]
	return ok && entry.ModTime.Equal(modTime)
}

func (c *Cache) Touch(path string, modTime time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[path] = Entry{ModTime: modTime}
}

func (c *Cache) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	b, err := json.MarshalIndent(c.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, b, 0o644)
}
