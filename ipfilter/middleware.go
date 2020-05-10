package ipfilter

import (
	"context"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// ServeHTTP implements middleware for http servers.
func (f Filter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if f.Wrapped == nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	ip := net.ParseIP(remoteIP)
	if ip == nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if f.Check(ip) == Accept {
		f.Wrapped.ServeHTTP(w, r)
		return
	}
	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	return
}

// UnaryServerInterceptor implements middleware for grpc servers
func (f Filter) UnaryServerInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Internal error getting peer IP")
	}
	remoteIP, _, _ := net.SplitHostPort(p.Addr.String())
	ip := net.ParseIP(remoteIP)
	if ip == nil {
		return nil, status.Errorf(codes.Internal, "Internal error getting peer IP")
	}
	if f.Check(ip) == Accept {
		return handler(ctx, req)
	}
	return nil, status.Errorf(codes.PermissionDenied, "IP is not allowed")
}

// StreamServerInterceptor implements middleware for grpc servers
func (f Filter) StreamServerInterceptor(srv interface{}, ss grpc.ServerStream,
	info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	p, ok := peer.FromContext(ctx)
	if !ok {
		return status.Errorf(codes.Internal, "Internal error getting peer IP")
	}
	remoteIP, _, _ := net.SplitHostPort(p.Addr.String())
	ip := net.ParseIP(remoteIP)
	if ip == nil {
		return status.Errorf(codes.Internal, "Internal error getting peer IP")
	}
	if f.Check(ip) == Accept {
		return handler(srv, ss)
	}
	return status.Errorf(codes.PermissionDenied, "IP is not allowed")
}
