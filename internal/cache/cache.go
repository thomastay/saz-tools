// Package cache provides a LRU cache for network sessions.
package cache

import (
	"fmt"
	"time"

	"github.com/bluele/gcache"
	"github.com/spaolacci/murmur3"
	"github.com/thomastay/saz-tools/pkg/parser"
)

// Cache can be created by `Create` and will contain the cached sessions.
type Cache struct {
	cache gcache.Cache
}

// Create returns a new instance of a network session cache.
func Create() *Cache {
	return &Cache{gcache.New(20).LRU().Build()}
}

// Get retrieves an network sessions from the cache with the specified `key`.
// If no sessions with the `key` exist, `nil` will be returned. The second
// boolean result will be `true` or `false` depending on the network sessions
// being returned or not.
func (cache *Cache) Get(key string) ([]parser.Session, bool) {
	entry, err := cache.cache.Get(key)
	if err != nil {
		return nil, false
	}
	return entry.([]parser.Session), true
}

// Put stores network sessions to the cache and returns a key with which they
// can be retrieved later.
func (cache *Cache) Put(sessions []parser.Session) (string, error) {
	key := computeKey(sessions)
	err := cache.cache.SetWithExpire(key, sessions, time.Minute*5)
	if err != nil {
		message := fmt.Sprintf("Putting network sessions with the key %s to cache failed.", key)
		return "", fmt.Errorf("%s\n%s", message, err.Error())
	}
	return key, nil
}

func computeKey(sessions []parser.Session) string {
	firstSession := &sessions[0]
	lastSession := &sessions[len(sessions)-1]
	identifier := firstSession.Request.Method + firstSession.Request.URL.String() +
		firstSession.Timers.ClientBeginRequest + firstSession.Timers.ClientDoneResponse +
		firstSession.Request.Method + lastSession.Request.URL.String() +
		firstSession.Timers.ClientBeginRequest + lastSession.Timers.ClientDoneResponse
	first, second := murmur3.Sum128([]byte(identifier))
	return fmt.Sprintf("%16x%16x", first, second)
}
