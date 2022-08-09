package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/hazelcast/hazelcast-go-client"
)

type Future struct {
	i         int
	processor interface{}
	key       interface{}
	//ch        chan interface{}
	f func(res interface{}, err error)
}

type Service struct {
	client *hazelcast.Client
	m      *hazelcast.Map
	wg     *sync.WaitGroup
	ch     chan Future
	state  int32
	mode   Mode
}

func StartNewService(ctx context.Context, config *Config) (*Service, error) {
	config.Client.Serialization.SetIdentifiedDataSerializableFactories(&EntryProcessorFactory{})
	client, err := hazelcast.StartNewClientWithConfig(ctx, config.Client)
	if err != nil {
		return nil, fmt.Errorf("starting Hazelcast client: %w", err)
	}
	m, err := client.GetMap(ctx, "sample-map")
	if err != nil {
		return nil, fmt.Errorf("getting map: %w", err)
	}
	mode := config.Mode()
	var wg *sync.WaitGroup
	var ch chan Future
	ch = make(chan Future, config.Concurrency)
	wg = &sync.WaitGroup{}
	wg.Add(config.Concurrency)
	svc := &Service{
		client: client,
		m:      m,
		wg:     wg,
		ch:     ch,
		mode:   mode,
	}
	if mode == pooledConcurrency {
		for i := 0; i < config.Concurrency; i++ {
			go svc.worker()
		}
	}
	return svc, nil
}

func (s *Service) Stop(ctx context.Context) error {
	if atomic.CompareAndSwapInt32(&s.state, 0, 1) {
		close(s.ch)
		if s.mode == pooledConcurrency || s.mode == allConcurrent {
			s.wg.Wait()
		}
		return s.client.Shutdown(ctx)
	}
	return nil
}

func (s *Service) Do(ctx context.Context, it Future) error {
	if s.mode == noConcurrency {
		return s.do(ctx, it)
	}
	if s.mode == allConcurrent {
		go func() {
			if err := s.do(ctx, it); err != nil {
				log.Printf("ERROR: %s", err.Error())
			}
			s.wg.Done()
		}()
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.ch <- it:
		return nil
	}
}

func (s *Service) worker() {
	for it := range s.ch {
		if err := s.do(context.Background(), it); err != nil {
			log.Print("ERROR: %s", err.Error())
		}
	}
	s.wg.Done()
}

func (s *Service) do(ctx context.Context, it Future) error {
	it.f(s.m.ExecuteOnKey(ctx, it.processor, it.key))
	return nil
}
