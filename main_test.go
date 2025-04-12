package main

import (
	"testing"
	"time"
)

func TestName(t *testing.T) {
	cache := NewCache(2 * time.Second)

	cache.Add("uuid1", Profile{
		UUID: "uuid1",
		Name: "John Doe",
		Orders: []*Order{
			{UUID: "order1", Value: 100, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	})
	profile := cache.Get("uuid1")
	if profile == nil {
		t.Errorf("Profile not found")
	}
	if profile.Name != "John Doe" {
		t.Errorf("Profile name is incorrect")
	}
}

func TestName2(t *testing.T) {
	cache := NewCache(2 * time.Second)

	cache.Add("uuid1", Profile{
		UUID: "uuid1",
		Name: "John Doe",
		Orders: []*Order{
			{UUID: "order1", Value: 100, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	})

	time.Sleep(3 * time.Second)

	profile := cache.Get("uuid1")
	if profile != nil {
		t.Errorf("Profile not found")
	}
}

func TestName3(t *testing.T) {
	cache := NewCache(2 * time.Second)

	cache.Add("uuid1", Profile{
		UUID: "uuid1",
		Name: "John Doe",
		Orders: []*Order{
			{UUID: "order1", Value: 100, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	})

	time.Sleep(10 * time.Second)

	profile := cache.Get("uuid1")
	if profile != nil {
		t.Errorf("Profile not found")
	}
}
