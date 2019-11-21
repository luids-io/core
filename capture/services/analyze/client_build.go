// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package analyze

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/luisguillenc/grpctls"
	"github.com/luisguillenc/yalogi"
	"google.golang.org/grpc"

	"github.com/luids-io/core/apiservice"
)

// ClientBuilder returns builder function for the apiservice
func ClientBuilder(opt ...ClientOption) apiservice.BuildFn {
	return func(def apiservice.ServiceDef, logger yalogi.Logger) (apiservice.Service, error) {
		//validates definition
		err := def.Validate()
		if err != nil {
			return nil, err
		}
		opts := make([]grpc.DialOption, 0)
		if def.Metrics {
			opts = append(opts, grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor))
			opts = append(opts, grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor))
		}
		//dial grpc
		dial, err := grpctls.Dial(def.Endpoint, def.ClientCfg(), opts...)
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
