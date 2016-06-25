package main

import (
	"errors"
)

type Program struct {
	Size
	Cells map[Index]*Cell
	read  Value
	write *Value
}

type Value int64

type Cell struct {
	Index
	Read, Write Value
	Symbol      rune
	Type        Type
}

type Type interface {
	Exec(Value) Value
	Bind(*Value) error
	SinkDir() Dir
}

var ErrBinding = errors.New("Error on binding")

type Forward struct {
	SinkTo Dir
	Input  *Value
}

func (s *Forward) Exec(Value) Value {
	return *s.Input
}

func (f *Forward) Bind(input *Value) error {
	if f.Input != nil {
		return ErrBinding
	}
	f.Input = input
	return nil
}

func (f *Forward) SinkDir() Dir {
	return f.SinkTo
}

type Constant struct {
	Value
}

func (c *Constant) Exec(Value) Value {
	return c.Value
}

func (c *Constant) Bind(*Value) error {
	return nil
}

func (c *Constant) SinkDir() Dir {
	return DirsPlane
}

type Oscillator struct {
	Clock    uint64
	Period   uint64
	Function func(uint64, uint64) Value
}

func (o *Oscillator) Exec(Value) Value {
	v := o.Function(o.Clock, o.Period)
	o.Clock += 1
	if o.Clock > o.Period {
		o.Clock = 0
	}
	return v
}

func (o *Oscillator) Bind(*Value) error {
	return nil
}

func (o *Oscillator) SinkDir() Dir {
	return DirsPlane
}
