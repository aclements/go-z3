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
import "runtime"

// BoolSort returns the boolean sort (type).
func (ctx *Context) BoolSort() *Sort {
	var csort C.Z3_sort
	ctx.do(func() {
		csort = C.Z3_mk_bool_sort(ctx.c)
	})
	return wrapSort(ctx, csort)
}

// FromBool returns a boolean expression with value val.
func (ctx *Context) FromBool(val bool) *Expr {
	var cexpr C.Z3_ast
	ctx.do(func() {
		if val {
			cexpr = C.Z3_mk_true(ctx.c)
		} else {
			cexpr = C.Z3_mk_false(ctx.c)
		}
	})
	return wrapExpr(ctx, cexpr)
}

// BoolConst returns a boolean constant named "name".
func (ctx *Context) BoolConst(name string) *Expr {
	return ctx.Const(name, ctx.BoolSort())
}

// AsBool returns the value of expr as a bool. expr must have boolean
// sort. If expr is not constant, AsBool returns false, false.
func (expr *Expr) AsBool() (val bool, isConst bool) {
	var res C.Z3_lbool
	expr.ctx.do(func() {
		res = C.Z3_get_bool_value(expr.ctx.c, expr.c)
	})
	runtime.KeepAlive(expr)
	return res == C.Z3_L_TRUE, res != C.Z3_L_UNDEF
}

//go:generate go run genwrap.go -- $GOFILE

// Eq returns a boolean expression that is true if l and r are equal.
//
// l and r must have the same sort.
//
//wrap:expr Eq Z3_mk_eq l r

// Distinct returns a boolean expression that is true if no two exprs
// are equal.
//
// All expressions in exprs must have the same sort.
//
//wrap:expr Distinct ctx:*Context exprs... : Z3_mk_distinct exprs...

// Not returns the boolean negation of l.
//
// l must have boolean sort.
//
//wrap:expr Not Z3_mk_not l

// IfThenElse returns an expression whose value is cons is cond is
// true, otherwise alt.
//
// cond must have boolean sort. cons and alt must have the same sort.
//
//wrap:expr IfThenElse Z3_mk_ite cond cons alt

// Iff returns an expression that is true if l and r are equal (l
// if-and-only-if r).
//
// l and r must have boolean sort.
//
//wrap:expr Iff Z3_mk_iff l r

// Implies returns an expression that is true if l implies r.
//
// l and r must have boolean sort.
//
//wrap:expr Implies Z3_mk_implies l r

// Xor returns an expression that is true if l xor r.
//
// l and r must have boolean sort.
//
//wrap:expr Xor Z3_mk_xor l r

// And returns an expression that is true if l and all arguments are
// true.
//
// All arguments must have boolean sort.
//
//wrap:expr And Z3_mk_and l r...

// Or returns an expression that is true if l or any argument is true.
//
// All arguments must have boolean sort.
//
//wrap:expr Or Z3_mk_or l r...
