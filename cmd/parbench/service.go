package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/hazelcast/hazelcast-go-client"
)

type item struct {
	i   int
	key interface{}
}

type Service struct {
	client *hazelcast.Client
	m      *hazelcast.Map
	wg     *sync.WaitGroup
	ch     chan item
	state  int32
	mode   Mode
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
	mode := config.Mode()
	var wg *sync.WaitGroup
	var ch chan item
	ch = make(chan item, config.Concurrency)
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

func (s *Service) Do(ctx context.Context, i int, key interface{}) error {
	if s.mode == noConcurrency {
		return s.do(ctx, i, key)
	}
	if s.mode == allConcurrent {
		go func() {
			if err := s.do(ctx, i, key); err != nil {
				log.Print("ERROR: %s", err.Error())
			}
			s.wg.Done()
		}()
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.ch <- item{i: i, key: key}:
		return nil
	}
}

func (s *Service) worker() {
	for it := range s.ch {
		if err := s.do(context.Background(), it.i, it.key); err != nil {
			log.Print("ERROR: %s", err.Error())
		}
	}
	s.wg.Done()
}

func (s *Service) do(ctx context.Context, i int, key interface{}) error {
	var err error
	t := measureTime(func() {
		err = s.m.Set(ctx, key, "value")
	})
	fmt.Printf("%06d\t%12d\n", i, t.Nanoseconds())
	return err
}
