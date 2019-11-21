// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package block

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/luisguillenc/grpctls"
	"github.com/luisguillenc/yalogi"
	"google.golang.org/grpc"

	"github.com/luids-io/core/apiservice"
)

// ClientBuilder returns builder function
func ClientBuilder(opt ...ClientOption) apiservice.BuildFn {
	return func(def apiservice.ServiceDef, logger yalogi.Logger) (apiservice.Service, error) {
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
		//creates client
		client := NewClient(dial, opt...)
		return client, nil
	}
}

func init() {
	apiservice.RegisterBuilder(ServiceName(), ClientBuilder())
}
