package cache

import (
	"fmt"
	"time"

	lru "github.com/bluele/gcache"
	parser "github.com/prantlf/saz-tools/pkg/parser"
	murmur "github.com/spaolacci/murmur3"
)

type Cache struct {
	cache lru.Cache
}

func Create() *Cache {
	return &Cache{lru.New(20).LRU().Build()}
}

func (cache *Cache) Get(key string) ([]parser.Session, bool) {
	entry, err := cache.cache.Get(key)
	if err != nil {
		return nil, false
	}
	return entry.([]parser.Session), true
}

func (cache *Cache) Put(sessions []parser.Session) string {
	key := computeKey(sessions)
	cache.cache.SetWithExpire(key, sessions, time.Minute*5)
	return key
}

func computeKey(sessions []parser.Session) string {
	firstSession := &sessions[0]
	lastSession := &sessions[len(sessions)-1]
	identifier := firstSession.Request.Method + firstSession.Request.URL.String() +
		firstSession.Timers.ClientBeginRequest + firstSession.Timers.ClientDoneResponse +
		firstSession.Request.Method + lastSession.Request.URL.String() +
		firstSession.Timers.ClientBeginRequest + lastSession.Timers.ClientDoneResponse
	first, second := murmur.Sum128([]byte(identifier))
	return fmt.Sprintf("%16x%16x", first, second)
}
