package main

import (
	"errors"
)

type Program struct {
	Size
	Cells         map[Index]Cell
	Input, Output chan Value
	Halt          chan int
}

type Value int64

type Cell struct {
	Index
	Value  Value
	Symbol rune
	Notify chan Value
	Type   Type
}

type Type interface {
	Exec(notify chan Value, halt chan int) bool
	Bind(input chan Value) error
	SinkDir() Dir
}

var ErrBinding = errors.New("Error on binding")

type Forward struct {
	SinkTo Dir
	Input  chan Value
}

func (s *Forward) Exec(notify chan Value, halt chan int) bool {
	select {
	case v, ok := <-s.Input:
		if !ok {
			return false
		}
		select {
		case <-halt:
			return false
		case notify <- v:
			return true
		}
	case <-halt:
		return false
	}
}

func (f *Forward) Bind(notify chan Value) error {
	if f.Input != nil {
		return ErrBinding
	}
	f.Input = notify
	return nil
}

func (f *Forward) SinkDir() Dir {
	return f.SinkTo
}

type Constant struct {
	Value
}

func (c *Constant) Exec(notify chan Value, halt chan int) bool {
	select {
	case <-halt:
		return false
	case notify <- c.Value:
		return true
	}
}

func (c *Constant) Bind(chan Value) error {
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

func (o *Oscillator) Exec(notify chan Value, halt chan int) bool {
	select {
	case <-halt:
		return false
	case notify <- o.Function(o.Clock, o.Period):
		o.Clock += 1
		if o.Clock > o.Period {
			o.Clock = 0
		}
		return true
	}
}

func (o *Oscillator) Bind(chan Value) error {
	return nil
}

func (o *Oscillator) SinkDir() Dir {
	return DirsPlane
}
