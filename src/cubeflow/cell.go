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
	Symbol rune
	Notify chan Value
	Type   Type
}

type Type interface {
	Exec(notify chan Value, halt chan int) bool
	Bind(input chan Value) error
}

var ErrBinding = errors.New("Error on binding")

type Forward struct {
	Input chan Value
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


