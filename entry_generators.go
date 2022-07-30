package main

import "encoding/binary"

type EntryGenerator interface {
	// GenerateKey generates a key corresponding to the index.
	// The generated key is assumed to be unique.
	// The same value should be returned for the same key.
	// Guaranteed to be called only by one goroutine at the same time.
	GenerateKey(index int64) interface{}
	// GenerateValue generates a value corresponding to the index.
	// May be called concurrently.
	// Guaranteed to be called only by one goroutine at the same time.
	GenerateValue(index int64) interface{}
}

type IdentityEntryGenerator struct{}

func (g IdentityEntryGenerator) GenerateKey(index int64) interface{} {
	return index
}

func (g IdentityEntryGenerator) GenerateValue(index int64) interface{} {
	return index
}

type SizedEntryGenerator struct {
	key   []byte
	value []byte
}

func NewSizedEntryGenerator(keySize, valueSize int) *SizedEntryGenerator {
	return &SizedEntryGenerator{
		key:   make([]byte, keySize),
		value: make([]byte, valueSize),
	}
}

func (g SizedEntryGenerator) GenerateKey(index int64) interface{} {
	binary.LittleEndian.PutUint64(g.key, uint64(index))
	return string(g.key)
}

func (g SizedEntryGenerator) GenerateValue(index int64) interface{} {
	binary.LittleEndian.PutUint64(g.value, uint64(index))
	return g.value
}

var entryGenerators = map[string]EntryGenerator{
	"identity":      IdentityEntryGenerator{},
	"sized128x1024": NewSizedEntryGenerator(128, 1024),
	"sized128x4096": NewSizedEntryGenerator(128, 4096),
}
