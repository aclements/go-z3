// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import (
	"math/big"
	"runtime"
	"unsafe"
)

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"

// Real is an expression with real sort.
//
// Real implements Expr.
type Real expr

func init() {
	sortWrappers[SortReal] = func(x expr) Expr {
		return Real(x)
	}
}

// RealSort returns the real sort.
func (ctx *Context) RealSort() *Sort {
	var csort C.Z3_sort
	ctx.do(func() {
		csort = C.Z3_mk_real_sort(ctx.c)
	})
	return wrapSort(ctx, csort, SortReal)
}

// RealConst returns a int constant named "name".
func (ctx *Context) RealConst(name string) Real {
	return ctx.Const(name, ctx.RealSort()).(Real)
}

// FromBigRat returns a real literal whose value is val.
func (ctx *Context) FromBigRat(val *big.Rat) Real {
	cstr := C.CString(val.String())
	defer C.free(unsafe.Pointer(cstr))
	sort := ctx.RealSort()
	var cexpr C.Z3_ast
	ctx.do(func() {
		cexpr = C.Z3_mk_numeral(ctx.c, cstr, sort.c)
	})
	runtime.KeepAlive(sort)
	return Real(wrapExpr(ctx, cexpr))
}

// AsRat returns the value of expr as a numerator and denominator Int
// literals. If expr is not a literal or is not rational, it returns
// false for isLiteralRational. To round an arbitrary real to be
// rational, see method Real.Approx.
func (expr Real) AsRat() (numer, denom Int, isLiteralRational bool) {
	if expr.astKind() != C.Z3_NUMERAL_AST {
		// Algebraic literals do not count as Z3_NUMERAL_AST,
		// so this gets all the cases we need.
		return Int{}, Int{}, false
	}
	var cnumer, cdenom C.Z3_ast
	expr.ctx.do(func() {
		cnumer = C.Z3_get_numerator(expr.ctx.c, expr.c)
	})
	expr.ctx.do(func() {
		cdenom = C.Z3_get_denominator(expr.ctx.c, expr.c)
	})
	runtime.KeepAlive(expr)
	return Int(wrapExpr(expr.ctx, cnumer)), Int(wrapExpr(expr.ctx, cdenom)), true
}

// AsBigRat returns the value of expr as a math/big.Rat. If expr is
// not a literal or is not rational, it returns nil, false.
func (expr Real) AsBigRat() (val *big.Rat, isLiteralRational bool) {
	numer, denom, isLiteralRational := expr.AsRat()
	if !isLiteralRational {
		return nil, false
	}
	var rat big.Rat
	bigNumer, _ := numer.AsBigInt()
	bigDenom, _ := denom.AsBigInt()
	rat.SetFrac(bigNumer, bigDenom)
	return &rat, true
}

// Approx approximates expr as two rational literals, where the
// difference between lower and upper is less than 1/10**precision. If
// expr is not a literal or is not irrational, it returns false for
// isLiteralIrrational.
func (expr Real) Approx(precision int) (lower, upper Real, isLiteralIrrational bool) {
	var isAlgebraicNumber bool
	expr.ctx.do(func() {
		// Despite the name, this really means an *irrational*
		// algebraic number.
		isAlgebraicNumber = z3ToBool(C.Z3_is_algebraic_number(expr.ctx.c, expr.c))
	})
	if !isAlgebraicNumber {
		return Real{}, Real{}, false
	}
	var clower, cupper C.Z3_ast
	expr.ctx.do(func() {
		clower = C.Z3_get_algebraic_number_lower(expr.ctx.c, expr.c, C.unsigned(precision))
	})
	expr.ctx.do(func() {
		cupper = C.Z3_get_algebraic_number_upper(expr.ctx.c, expr.c, C.unsigned(precision))
	})
	runtime.KeepAlive(expr)
	return Real(wrapExpr(expr.ctx, clower)), Real(wrapExpr(expr.ctx, cupper)), true
}

// TODO: AsBigFloat? AsFloat64? AsFloat32? I don't actually know how
// to implement those without potentially double rounding.

//go:generate go run genwrap.go -t Real $GOFILE intreal.go

// Div returns l / r.
//
// If r is 0, the result is unconstrained.
//
//wrap:expr Div Z3_mk_div l r

// ToInt returns the floor of l as sort Int.
//
// Note that this is not truncation. For example, ToInt(-1.3) is -2.
//
//wrap:expr ToInt:Int Z3_mk_real2int l

// IsInt returns an expression that is true if l has no fractional
// part.
//
//wrap:expr IsInt:Bool Z3_mk_is_int l
