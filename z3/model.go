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

// A Model is a binding of constants that satisfies a set of formulas.
type Model struct {
	*modelImpl
	noEq
}

type modelImpl struct {
	ctx *Context
	c   C.Z3_model
}

func wrapModel(ctx *Context, c C.Z3_model) *Model {
	impl := &modelImpl{ctx, c}
	ctx.lock.Lock()
	C.Z3_model_inc_ref(ctx.c, c)
	ctx.lock.Unlock()
	runtime.SetFinalizer(impl, func(impl *modelImpl) {
		impl.ctx.do(func() {
			C.Z3_model_dec_ref(impl.ctx.c, impl.c)
		})
	})
	return &Model{impl, noEq{}}
}

// Eval evaluates val using the concrete interpretations of constants
// and functions in model m.
//
// If completion is true, it will assign interpretations for any
// constants or functions that currently don't have an interpretation
// in m. Otherwise, the resulting value may not be concrete.
//
// Eval returns nil if val cannot be evaluated. This can happen if val
// contains a quantifier or is type-incorrect, or if m is a partial
// model (that is, the option MODEL_PARTIAL was set to true).
func (m *Model) Eval(val Value, completion bool) Value {
	var ok bool
	var ast AST
	m.ctx.do(func() {
		var cast C.Z3_ast
		ok = z3ToBool(C.Z3_model_eval(m.ctx.c, m.c, val.impl().c, boolToZ3(completion), &cast))
		if ok {
			ast = wrapAST(m.ctx, cast)
		}
	})
	runtime.KeepAlive(m)
	runtime.KeepAlive(val)
	if !ok {
		return nil
	}
	return ast.AsValue()
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
