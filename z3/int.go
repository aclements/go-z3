// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"
import "math/big"

// Int is an expression with int sort.
//
// Ints are mathematical integers with infinite precision.
//
// Int implements Expr.
type Int expr

func init() {
	sortWrappers[SortInt] = func(x expr) Expr {
		return Int(x)
	}
}

// IntSort returns the integer sort.
func (ctx *Context) IntSort() *Sort {
	var csort C.Z3_sort
	ctx.do(func() {
		csort = C.Z3_mk_int_sort(ctx.c)
	})
	return wrapSort(ctx, csort, SortInt)
}

// IntConst returns a int constant named "name".
func (ctx *Context) IntConst(name string) Int {
	return ctx.Const(name, ctx.IntSort()).(Int)
}

// AsInt64 returns the value of expr as an int64. If expr is not a
// literal, it returns 0, false, false. If expr is a literal, but its
// value cannot be represented as an int64, it returns 0, true, false.
func (expr Int) AsInt64() (val int64, isLiteral, ok bool) {
	return expr.asInt64()
}

// AsUint64 is like AsInt64, but returns a uint64 and fails if expr
// cannot be represented as a uint64.
func (expr Int) AsUint64() (val uint64, isLiteral, ok bool) {
	return expr.asUint64()
}

// AsBigInt returns the value of expr as a math/big.Int. If expr is
// not a literal, it returns nil, false.
func (expr Int) AsBigInt() (val *big.Int, isConst bool) {
	return expr.asBigInt()
}

//go:generate go run genwrap.go -t Int $GOFILE intreal.go

// Div returns the floor of l / r.
//
// If r is 0, the result is unconstrained.
//
// Note that this differs from Go division: Go rounds toward zero
// (truncated division), whereas this rounds toward -inf.
//
//wrap:expr Div Z3_mk_div l r

// Mod returns modulus of l / r.
//
// The sign of the result follows the sign of r.
//
//wrap:expr Mod Z3_mk_mod l r

// Rem returns remainder of l / r.
//
// The sign of the result follows the sign of l.
//
//wrap:expr Rem Z3_mk_rem l r

// ToReal returns an expression that converts l to sort Real.
//
//wrap:expr ToReal:Real Z3_mk_int2real l

// ToBV converts l to a bit-vector of width bits.
//
//wrap:expr ToBV:BV l bits:int : Z3_mk_int2bv bits:unsigned l
