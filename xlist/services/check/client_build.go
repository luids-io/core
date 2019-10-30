// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package check

import (
	"errors"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/luisguillenc/grpctls"
	"github.com/luisguillenc/yalogi"
	"google.golang.org/grpc"

	"github.com/luids-io/core/apiservice"
	"github.com/luids-io/core/option"
	"github.com/luids-io/core/xlist"
)

// ClientBuilder returns builder function
func ClientBuilder(opt ...ClientOption) apiservice.BuildFn {
	return func(def apiservice.Definition, logger yalogi.Logger) (apiservice.Service, error) {
		//validates definition
		err := def.Validate()
		if err != nil {
			return nil, err
		}
		dopts := make([]grpc.DialOption, 0)
		if def.Metrics {
			dopts = append(dopts, grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor))
			dopts = append(dopts, grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor))
		}
		//dial grpc
		dial, err := grpctls.Dial(def.Endpoint, def.ClientCfg(), dopts...)
		if err != nil {
			return nil, err
		}
		if len(def.Opts) > 0 {
			// parse and set cache options
			ttl, negativettl, err := parseCacheOpts(def.Opts)
			if err != nil {
				return nil, err
			}
			if ttl > 0 || negativettl > 0 {
				opt = append(opt, SetCache(ttl, negativettl))
			}
		}
		//creates client
		client := NewClient(dial, []xlist.Resource{}, opt...)
		return client, nil
	}
}

func parseCacheOpts(opts map[string]interface{}) (int, int, error) {
	var ttl, negativettl int
	// get ttl
	value, ok, err := option.Int(opts, "ttl")
	if err != nil {
		return 0, 0, err
	}
	if ok {
		if value < 0 {
			return 0, 0, errors.New("invalid 'ttl'")
		}
		ttl = value
	}
	// get negativettl
	value, ok, err = option.Int(opts, "negativettl")
	if err != nil {
		return 0, 0, err
	}
	if ok {
		if value < 0 {
			return 0, 0, errors.New("invalid 'negativettl'")
		}
		negativettl = value
	} else {
		if ttl > 0 {
			negativettl = ttl
		}
	}
	return ttl, negativettl, nil
}

func init() {
	apiservice.RegisterBuilder(ServiceName(), ClientBuilder())
}
