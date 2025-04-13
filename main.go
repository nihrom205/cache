/*
Необходимо написать in-memory кэш, который будет по ключу (uuid пользователя) возвращать профиль и список его заказов.

1. У кэша должен быть TTL (2 сек)
2. Кэшем может пользоваться функция(-и), которая работает с заказами (добавляет/обновляет/удаляет).
Если TTL истек, то возвращается nil. При апдейте TTL снова устанавливается 2 сек. Методы должны быть потокобезопасными
3. Должны быть написаны тестовые сценарии использования данного кэша
(базовые структуры не менять)

Доп задание: автоматическая очистка истекших записей

type Profile struct {
	UUID   string
	Name   string
	Orders []*Order
}

type Order struct {
	UUID      string
	Value     any
	CreatedAt time.Time
	UpdatedAt time.Time
}
*/

package main

import (
	"sync"
	"time"
)

type Profile struct {
	UUID   string
	Name   string
	Orders []*Order
}

type Order struct {
	UUID      string
	Value     any
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ItemCache struct {
	profile  *Profile
	CreateAt time.Time
}

type Cache struct {
	cache map[string]*ItemCache
	ttl   time.Duration
	mu    sync.RWMutex
}

func NewCache(ttl time.Duration) *Cache {
	cache := &Cache{
		cache: make(map[string]*ItemCache),
		ttl:   ttl,
		mu:    sync.RWMutex{},
	}
	go cache.cleanUp()
	return cache
}

func (c *Cache) Set(uuid string, profile *Profile) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[uuid] = &ItemCache{
		profile:  profile,
		CreateAt: time.Now(),
	}
}

func (c *Cache) Get(uuid string) (*Profile, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.cache[uuid]
	if !ok {
		return nil, false
	}
	if time.Since(item.CreateAt) > c.ttl {
		return nil, false
	}
	return item.profile, true
}

func (c *Cache) cleanUp() {
	tick := time.NewTicker(c.ttl)
	defer tick.Stop()
	for range tick.C {
		c.mu.Lock()
		for uuid, item := range c.cache {
			if time.Since(item.CreateAt) > c.ttl {
				delete(c.cache, uuid)
			}
		}
		c.mu.Unlock()
	}
}
