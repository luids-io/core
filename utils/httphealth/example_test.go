package httphealth_test

import (
	"errors"
	"log"
	"math/rand"
	"net"

	"github.com/luids-io/core/utils/httphealth"
)

// service is a supervised object
type service struct{}

// Ping implements httphealth.Pingable interface
func (s service) Ping() error {
	if rand.Intn(10) > 8 {
		return errors.New("error in supervised")
	}
	return nil
}

func Example() {
	lis, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		log.Fatalf("listening: %v", err)
	}
	health := httphealth.New(&service{})
	health.Serve(lis)
}
