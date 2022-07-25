package memory

import (
	"auth/internal/model/principal"
	"testing"
	"time"
)

func TestEntry_Expired(t *testing.T) {
	e := entry{
		p:        principal.Principal{},
		expireAt: time.Now().Add(1 * time.Second),
	}

	time.Sleep(900 * time.Millisecond)

	t.Log(e.expired())
}

func TestInmemory_Get(t *testing.T) {
	var c InmemoryCache
	p, err := c.Get("1")
	t.Log(p, err)
	t.Log(c.HitCount(), c.MissCount())
}

func TestInmemory_EntryCount(t *testing.T) {
	var c InmemoryCache
	c.expireIn = 3 * time.Second
	cancel := c.StartCleaning(1 * time.Second)
	var p principal.Principal
	err := c.Set("1", p)
	err = c.Set("2", p)
	err = c.Set("3", p)
	t.Log(err)
	t.Log(c.EntryCount())
	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()
	for {
		t.Log(c.EntryCount())
		time.Sleep(500 * time.Millisecond)
	}
}
