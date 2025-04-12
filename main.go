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

type Cache struct {
	cache map[string]*Profile
	ttl   time.Duration
	mu    sync.RWMutex
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		cache: make(map[string]*Profile),
		ttl:   ttl,
	}
}

func (c *Cache) Add(uuid string, profile Profile) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[uuid] = &profile
	c.clean()
}

func (c *Cache) Get(uuid string) *Profile {
	//Если TTL истек, то возвращается nil. При апдейте TTL снова устанавливается 2 сек. Методы должны быть потокобезопасными
	c.mu.RLock()
	defer c.mu.RUnlock()
	profile, ok := c.cache[uuid]
	if !ok {
		return nil
	}
	for _, order := range profile.Orders {
		if time.Since(order.CreatedAt) > c.ttl && time.Since(order.UpdatedAt) > c.ttl {
			return nil
		}
	}
	return profile
}

func (c *Cache) Update(uuid string, profile Profile) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.cache[uuid]
	if !ok {
		return
	}

	for _, p := range profile.Orders {
		p.UpdatedAt = time.Now()
	}
	c.cache[uuid] = &profile
}

func (c *Cache) clean() {
	// по заданному условию автоматическая очистка истекших записей
	go func() {

		for {
			select {
			case <-time.After(2 * time.Second):
				c.mu.Lock()
				for uuid, profile := range c.cache {
					for _, order := range profile.Orders {
						if time.Since(order.CreatedAt) > c.ttl && time.Since(order.UpdatedAt) > c.ttl {
							delete(c.cache, uuid)
							return
						}
					}
				}
				c.mu.Unlock()
			default:
			}
		}
	}()
}
