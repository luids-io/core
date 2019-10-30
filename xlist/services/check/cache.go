// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package check

import (
	"fmt"
	"time"

	cacheimpl "github.com/patrickmn/go-cache"

	"github.com/luids-io/core/xlist"
)

const defaultCacheCleanups = 6 * time.Minute

// cache implements a cache
type cache struct {
	ttl         int
	negativettl int
	cachei      *cacheimpl.Cache
}

// newCache returns a cache
func newCache(ttl, negativettl int, cleanups time.Duration) *cache {
	c := &cache{
		ttl:         ttl,
		negativettl: negativettl,
		cachei:      cacheimpl.New(time.Duration(ttl)*time.Second, cleanups),
	}
	return c
}

// Flush cleas all items from cache
func (c *cache) flush() {
	c.cachei.Flush()
}

func (c *cache) get(name string, resource xlist.Resource) (xlist.Response, bool) {
	key := fmt.Sprintf("%s_%s", resource.String(), name)
	hit, exp, ok := c.cachei.GetWithExpiration(key)
	if ok {
		r := hit.(xlist.Response)
		if r.TTL >= 0 {
			//updates ttl
			ttl := exp.Sub(time.Now()).Seconds()
			if ttl < 0 { //nonsense
				panic("cache missfunction")
			}
			r.TTL = int(ttl)
		}
		return r, true
	}
	return xlist.Response{}, false
}

func (c *cache) set(name string, resource xlist.Resource, r xlist.Response) xlist.Response {
	//if don't cache
	if (r.TTL == xlist.NeverCache) || (!r.Result && c.negativettl == xlist.NeverCache) {
		return r
	}
	//sets cache
	ttl := c.ttl
	if !r.Result && c.negativettl > 0 {
		ttl = c.negativettl
	}
	if r.TTL < ttl { //minor than cachettl
		r.TTL = ttl //sets reponse to cachettl
	}
	if r.TTL > 0 {
		key := fmt.Sprintf("%s_%s", resource.String(), name)
		c.cachei.Set(key, r, time.Duration(r.TTL)*time.Second)
	}
	return r
}
