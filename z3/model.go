// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import "runtime"

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"

type Model struct {
	ctx *Context
	c   C.Z3_model
	noEq
}

func wrapModel(ctx *Context, c C.Z3_model) *Model {
	model := &Model{ctx, c, noEq{}}
	ctx.lock.Lock()
	C.Z3_model_inc_ref(ctx.c, c)
	ctx.lock.Unlock()
	// TODO: Don't attach finalizer to a user-visible pointer.
	runtime.SetFinalizer(model, func(model *Model) {
		model.ctx.do(func() {
			C.Z3_model_dec_ref(model.ctx.c, model.c)
		})
	})
	return model
}

// Eval evaluates expr using the values in model m.
//
// If completion is true, it will assign interpretations for any
// constants or functions that currently don't have an interpretation
// in m.
//
// Eval returns nil is expr cannot be evaluated. This can happen if
// expr contains a quantifier or is type-incorrect, or if m is a
// partial model (that is, the option MODEL_PARTIAL was set to true).
func (m *Model) Eval(expr Expr, completion bool) Expr {
	var out C.Z3_ast
	var ok bool
	m.ctx.do(func() {
		ok = z3ToBool(C.Z3_model_eval(m.ctx.c, m.c, expr.impl().c, boolToZ3(completion), &out))
	})
	runtime.KeepAlive(m)
	runtime.KeepAlive(expr)
	if !ok {
		return nil
	}
	return wrapExpr(m.ctx, out).lift(KindUnknown)
}

// String returns a string representation of m.
func (m *Model) String() string {
	var res string
	m.ctx.do(func() {
		res = C.GoString(C.Z3_model_to_string(m.ctx.c, m.c))
	})
	runtime.KeepAlive(m)
	return res
}
