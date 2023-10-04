package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		i := map[string]int{
			"aaa": 100,
			"bbb": 200,
			"ccc": 300,
		}
		c := NewCache(3)
		for k, v := range i {
			c.Set(Key(k), v)
		}
		for k := range i {
			_, ok := c.Get(Key(k))
			require.True(t, ok)
		}

		c.Clear()

		for k := range i {
			_, ok := c.Get(Key(k))
			require.False(t, ok)
		}
	})

	t.Run("push out logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", 100) // last
		c.Set("bbb", 200) // middle
		c.Set("ccc", 300) // first

		c.Set(Key("ddd"), 400) // ddd->first

		_, ok := c.Get("aaa") // should be pushed out
		require.False(t, ok)
	})

	t.Run("least recently used", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", 100) // last
		c.Set("bbb", 200) // middle
		c.Set("ccc", 300) // first

		c.Set("bbb", 500) // bbb->first
		c.Set("aaa", 600) // aaa->first

		c.Get("aaa") // aaa->first
		c.Get("bbb") // bbb->first
		c.Get("aaa") // aaa->first

		c.Set("ddd", 400) // ddd->first

		_, ok := c.Get("ccc") // should be pushed out
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
