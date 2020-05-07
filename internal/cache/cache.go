package cache

import (
	"fmt"
	"time"

	lru "github.com/bluele/gcache"
	analyzer "github.com/prantlf/saz-tools/pkg/analyzer"
	parser "github.com/prantlf/saz-tools/pkg/parser"
	murmur "github.com/spaolacci/murmur3"
)

type Entry struct {
	RawSessions  []parser.Session
	FineSessions []analyzer.Session
}

type Cache struct {
	cache lru.Cache
}

func Create() *Cache {
	return &Cache{lru.New(20).LRU().Build()}
}

func (cache *Cache) Get(key string) (Entry, bool) {
	entry, err := cache.cache.Get(key)
	if err != nil {
		return Entry{}, false
	}
	return entry.(Entry), true
}

func (cache *Cache) Put(entry Entry) string {
	key := computeKey(entry)
	cache.cache.SetWithExpire(key, entry, time.Minute*5)
	return key
}

func computeKey(entry Entry) string {
	firstSession := entry.FineSessions[0]
	lastSession := entry.FineSessions[len(entry.FineSessions)-1]
	identifier := firstSession.Request.Method + firstSession.Request.URL.Full +
		firstSession.Timers.ClientBeginRequest + firstSession.Timers.ClientDoneResponse +
		firstSession.Request.Method + lastSession.Request.URL.Full +
		firstSession.Timers.ClientBeginRequest + lastSession.Timers.ClientDoneResponse
	first, second := murmur.Sum128([]byte(identifier))
	return fmt.Sprintf("%16x%16x", first, second)
}
