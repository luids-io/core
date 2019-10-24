// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package resolvcheck

import (
	"fmt"
	"net"
	"time"

	cacheimpl "github.com/patrickmn/go-cache"

	"github.com/luids-io/core/dnsutil"
)

const defaultCacheCleanups = 5 * time.Minute

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

func (c *cache) get(client, resolved net.IP, name string) (dnsutil.ResolvResponse, bool) {
	key := fmt.Sprintf("%s_%s_%s", client.String(), resolved.String(), name)
	hit, ok := c.cachei.Get(key)
	if ok {
		r := hit.(dnsutil.ResolvResponse)
		return r, true
	}
	return dnsutil.ResolvResponse{}, false
}

func (c *cache) set(client, resolved net.IP, name string, r dnsutil.ResolvResponse) dnsutil.ResolvResponse {
	//if don't cache
	if !r.Result && c.negativettl == 0 {
		return r
	}
	//sets cache
	ttl := c.ttl
	if !r.Result && c.negativettl > 0 {
		ttl = c.negativettl
	}
	if ttl > 0 {
		key := fmt.Sprintf("%s_%s_%s", client.String(), resolved.String(), name)
		c.cachei.Set(key, r, time.Duration(ttl)*time.Second)
	}
	return r
}
