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
	Type   Type
	Read   Value
	Write  Value
	Symbol rune
}

type Type interface {
	Exec(Value) Value
	Bind(*Value) error
	OfferDir() Dir
	RequestDir() Dir
}

var ErrBinding = errors.New("error on binding")
var ErrTooManyBinings = errors.New("too many bindings")

type Forward struct {
	SourceDir Dir
	SinkDir   Dir
	Input     *Value
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

func (f *Forward) RequestDir() Dir {
	return f.SourceDir
}

func (f *Forward) OfferDir() Dir {
	return f.SinkDir
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

func (c *Constant) OfferDir() Dir {
	return DirsPlane
}

func (c *Constant) RequestDir() Dir {
	return DirNone
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

func (o *Oscillator) OfferDir() Dir {
	return DirsPlane
}

func (o *Oscillator) RequestDir() Dir {
	return DirNone
}

type BinaryOp struct {
	A, B     *Value
	Function func(Value, Value) Value
}

func (b *BinaryOp) Exec(Value) Value {
	return b.Function(*b.A, *b.B)
}

func (b *BinaryOp) Bind(input *Value) error {
	if b.A == nil {
		b.A = input
		return nil
	}
	if b.B == nil {
		b.B = input
		return nil
	}
	return ErrTooManyBinings
}

func (b *BinaryOp) OfferDir() Dir {
	return DirsPlane
}

func (b *BinaryOp) RequestDir() Dir {
	return DirsPlane
}

type UnaryOp struct {
	V   *Value
	Function func(Value) Value
}

func (u *UnaryOp) Exec(Value) Value {
	return u.Function(*u.V)
}

func (u *UnaryOp) Bind(input *Value) error {
	if u.V == nil {
		u.V = input
		return nil
	}
	return ErrTooManyBinings
}

func (u *UnaryOp) OfferDir() Dir {
	return DirsPlane
}

func (u *UnaryOp) RequestDir() Dir {
	return DirsPlane
}

type Pulse struct {
	Value
}

func (p *Pulse) Exec(v Value) Value {
	if (p.Value == 0) {
		return 0
	}
	v = p.Value
	p.Value = 0
	return v
}

func (*Pulse) Bind(*Value) error {
	return nil
}

func (p *Pulse) OfferDir() Dir {
	return DirsPlane
}

func (p *Pulse) RequestDir() Dir {
	return DirNone
}
