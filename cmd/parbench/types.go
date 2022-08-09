package main

import (
	"fmt"

	"github.com/hazelcast/hazelcast-go-client/serialization"
)

type EntryProcessor1 struct {
	value string
}

func (e EntryProcessor1) FactoryID() int32 {
	return 38
}

func (e EntryProcessor1) ClassID() int32 {
	return 555
}

func (e EntryProcessor1) WriteData(output serialization.DataOutput) {
	output.WriteString(e.value)
}

func (e *EntryProcessor1) ReadData(input serialization.DataInput) {
	e.value = input.ReadString()
}

type EntryProcessorFactory struct {
}

func (e EntryProcessorFactory) Create(id int32) serialization.IdentifiedDataSerializable {
	if id == 555 {
		return &EntryProcessor1{}
	}
	panic(fmt.Errorf("unknown id: %d", id))
}

func (e EntryProcessorFactory) FactoryID() int32 {
	return 38
}
