package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/hazelcast/hazelcast-go-client"
)

type Service struct {
	client *hazelcast.Client
	m      *hazelcast.Map
	wg     *sync.WaitGroup
	ch     chan interface{}
	state  int32
}

func StartNewService(ctx context.Context, config *Config) (*Service, error) {
	client, err := hazelcast.StartNewClientWithConfig(ctx, config.Client)
	if err != nil {
		return nil, fmt.Errorf("starting Hazelcast client: %w", err)
	}
	m, err := client.GetMap(ctx, "sample-map")
	if err != nil {
		return nil, fmt.Errorf("getting map: %w", err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(config.OperationCount)
	svc := &Service{
		client: client,
		m:      m,
		wg:     wg,
		ch:     make(chan interface{}, config.Concurrency),
	}
	for i := 0; i < config.Concurrency; i++ {
		go svc.worker()
	}
	return svc, nil
}

func (s *Service) Stop(ctx context.Context) error {
	if atomic.CompareAndSwapInt32(&s.state, 0, 1) {
		close(s.ch)
		s.wg.Wait()
		return s.client.Shutdown(ctx)
	}
	return nil
}

func (s *Service) Do(ctx context.Context, key interface{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.ch <- key:
		return nil
	}
}

func (s *Service) worker() {
	for k := range s.ch {
		if err := s.m.Set(context.Background(), k, "value"); err != nil {
			log.Print("ERROR: %s", err.Error())
		}
	}
	s.wg.Done()
}
