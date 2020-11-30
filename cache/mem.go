package cache

import (
	"log"
	"sync"
	"time"
)

// MemCache is used by manager to store and retrieve key value data
type MemCache interface {
	Get(key string) (interface{}, bool)
	Remove(size int) []interface{}
	RemoveKey(key string) (interface{}, bool)
	Put(key string, data interface{})
	PutWithExpireTime(key string, data interface{}, expiredTime time.Duration)
	Delete(key string)
	Size() int
	BatchGet(size int) []interface{}
}

type cdata struct {
	d               interface{}
	expireTimePoint time.Time
}

func (d *cdata) expired() bool {
	return !d.expireTimePoint.IsZero() && d.expireTimePoint.After(time.Now())
}

type memCache struct {
	mu                     *sync.Mutex
	keyData                map[string]*cdata
	defaultExpireTime      time.Duration
	defaultClearExpireTime time.Duration
}

func (m *memCache) BatchGet(size int) []interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	var ds []interface{}
	for _, v := range m.keyData {
		ds = append(ds, v.d)
		if len(ds) == size {
			break
		}
	}
	return ds
}

func (m *memCache) RemoveKey(key string) (interface{}, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, ok := m.keyData[key]
	if ok {
		delete(m.keyData, key)
		return data.d, ok
	}
	return nil, ok

}

func (m *memCache) Size() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.keyData)
}

func (m *memCache) Remove(size int) []interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	var ds []interface{}
	for k, v := range m.keyData {
		delete(m.keyData, k)
		ds = append(ds, v.d)
		if len(ds) == size {
			break
		}
	}
	return ds
}

// NewMemCache if defaultExpireTime equals 0, the cache data with defaultExpireTime will not be check whether expire
// if defaultClearExpireTime equals 0, all cache will not check expire
// Mem cache is thread safe
func NewMemCache(defaultExpireTime time.Duration, defaultClearExpireTime time.Duration) MemCache {
	m := &memCache{
		mu:                     &sync.Mutex{},
		keyData:                make(map[string]*cdata),
		defaultClearExpireTime: defaultClearExpireTime,
		defaultExpireTime:      defaultExpireTime,
	}
	go m.clear()
	return m
}

// clear cache in period time set by defaultClearExpireTime
func (m *memCache) clear() {
	if m.defaultClearExpireTime == 0 {
		return
	}
	for {
		time.Sleep(m.defaultClearExpireTime)
		log.Printf("before clear cache size:%d \n", m.Size())
		m.mu.Lock()
		now := time.Now()
		for k, d := range m.keyData {
			if !d.expireTimePoint.IsZero() && d.expireTimePoint.Before(now) {
				log.Printf("clear k:%s", k)
				delete(m.keyData, k)
			}
		}
		m.mu.Unlock()
		log.Printf("after clear cache size:%d \n", m.Size())
	}
}

func (m *memCache) Get(key string) (interface{}, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, ok := m.keyData[key]
	if ok {
		return data.d, ok
	}
	return nil, ok
}

func (m *memCache) Put(key string, d interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cd := &cdata{d: d}
	if m.defaultExpireTime != 0 {
		cd.expireTimePoint = time.Now().Add(m.defaultExpireTime)
	}
	m.keyData[key] = cd

}

func (m *memCache) PutWithExpireTime(key string, data interface{}, expiredTime time.Duration) {
	cd := &cdata{d: data}
	if expiredTime == 0 {
		cd.expireTimePoint = time.Now().Add(expiredTime)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.keyData[key] = cd
}

func (m *memCache) Delete(key string) {
	m.mu.Lock()
	delete(m.keyData, key)
	m.mu.Unlock()
}
