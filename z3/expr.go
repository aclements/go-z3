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

type noEq struct {
	_ [0]func()
}

type Expr struct {
	ctx *Context
	c   C.Z3_ast
	noEq
}

func wrapExpr(ctx *Context, c C.Z3_ast) *Expr {
	expr := &Expr{ctx, c, noEq{}}
	// Note that, even if c was just returned by an allocation
	// function, we're still responsible for incrementing its
	// reference count. This is weird, but also nice because we
	// can wrap any AST that comes out of the Z3 API, even if
	// we've already wrapped it, and the reference count will
	// protect the underlying object no matter what happens to the
	// Go wrappers.
	ctx.lock.Lock()
	C.Z3_inc_ref(ctx.c, c)
	ctx.lock.Unlock()
	runtime.SetFinalizer(expr, func(expr *Expr) {
		expr.ctx.do(func() {
			C.Z3_dec_ref(expr.ctx.c, expr.c)
		})
	})
	return expr
}

// Const returns a constant named "name" of the given sort. This
// constant will be same as all other constants created with this
// name.
func (ctx *Context) Const(name string, sort *Sort) *Expr {
	sym := ctx.symbol(name)
	var cexpr C.Z3_ast
	ctx.do(func() {
		cexpr = C.Z3_mk_const(ctx.c, sym, sort.c)
	})
	runtime.KeepAlive(sort)
	return wrapExpr(ctx, cexpr)
}

// FreshConst returns a constant that is distinct from all other
// constants. The name will begin with "prefix".
func (ctx *Context) FreshConst(prefix string, sort *Sort) *Expr {
	cprefix := C.CString(prefix)
	defer C.free(unsafe.Pointer(cprefix))
	var cexpr C.Z3_ast
	ctx.do(func() {
		cexpr = C.Z3_mk_fresh_const(ctx.c, cprefix, sort.c)
	})
	runtime.KeepAlive(sort)
	return wrapExpr(ctx, cexpr)
}

// FromBigInt returns a constant expression whose value is val. sort
// must have kind int, real, finite-domain, or bit-vector.
func (ctx *Context) FromBigInt(val *big.Int, sort *Sort) *Expr {
	cstr := C.CString(val.Text(10))
	defer C.free(unsafe.Pointer(cstr))
	var cexpr C.Z3_ast
	ctx.do(func() {
		cexpr = C.Z3_mk_numeral(ctx.c, cstr, sort.c)
	})
	runtime.KeepAlive(sort)
	return wrapExpr(ctx, cexpr)
}

// FromInt returns a constant expression whose value is val. sort must
// have kind int, finite-domain, or bit-vector.
func (ctx *Context) FromInt(val int64, sort *Sort) *Expr {
	var cexpr C.Z3_ast
	ctx.do(func() {
		cexpr = C.Z3_mk_int64(ctx.c, C.__int64(val), sort.c)
	})
	runtime.KeepAlive(sort)
	return wrapExpr(ctx, cexpr)
}

// String returns a string representation of expr.
func (expr *Expr) String() string {
	var res string
	expr.ctx.do(func() {
		res = C.GoString(C.Z3_ast_to_string(expr.ctx.c, expr.c))
	})
	runtime.KeepAlive(expr)
	return res
}

// Equal returns true if expr and o are structurally identical.
func (expr *Expr) Equal(o *Expr) bool {
	var out bool
	expr.ctx.do(func() {
		out = z3ToBool(C.Z3_is_eq_ast(expr.ctx.c, expr.c, o.c))
	})
	runtime.KeepAlive(expr)
	runtime.KeepAlive(o)
	return out
}

// Sort returns expr's sort (that is, its Z3 type).
func (expr *Expr) Sort() *Sort {
	var csort C.Z3_sort
	expr.ctx.do(func() {
		csort = C.Z3_get_sort(expr.ctx.c, expr.c)
	})
	runtime.KeepAlive(expr)
	return wrapSort(expr.ctx, csort)
}

func (expr *Expr) astKind() C.Z3_ast_kind {
	var ckind C.Z3_ast_kind
	expr.ctx.do(func() {
		ckind = C.Z3_get_ast_kind(expr.ctx.c, expr.c)
	})
	runtime.KeepAlive(expr)
	return ckind
}
