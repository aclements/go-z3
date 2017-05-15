// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import "runtime"

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
*/
import "C"

// An AST is the abstract syntax tree of an expression, sort, etc.
type AST struct {
	*astImpl
	noEq
}

type astImpl struct {
	ctx *Context
	c   C.Z3_ast
}

func wrapAST(ctx *Context, c C.Z3_ast) AST {
	impl := &astImpl{ctx, c}
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
	runtime.SetFinalizer(impl, func(impl *astImpl) {
		impl.ctx.do(func() {
			C.Z3_dec_ref(impl.ctx.c, impl.c)
		})
	})
	return AST{impl, noEq{}}
}

// Equal returns true if ast and o are identical ASTs.
func (ast AST) Equal(o AST) bool {
	// Sadly, while AST equality is just pointer equality on the
	// underlying C pointers, it's impossible to expose this as Go
	// equality because we need a Go pointer to attach the
	// finalizer to. We can't make *that* pointer 1:1 with the C
	// pointer without making the object permanently live.
	var out bool
	ast.ctx.do(func() {
		out = z3ToBool(C.Z3_is_eq_ast(ast.ctx.c, ast.c, o.c))
	})
	runtime.KeepAlive(ast)
	runtime.KeepAlive(o)
	return out
}

// String returns ast as an S-expression.
func (ast AST) String() string {
	var res string
	ast.ctx.do(func() {
		res = C.GoString(C.Z3_ast_to_string(ast.ctx.c, ast.c))
	})
	runtime.KeepAlive(ast)
	return res
}

// Hash returns a hash of ast. Structurally identical ASTs will have
// the same hash code.
func (ast AST) Hash() uint64 {
	var res uint64
	ast.ctx.do(func() {
		res = uint64(C.Z3_get_ast_hash(ast.ctx.c, ast.c))
	})
	runtime.KeepAlive(ast)
	return res
}

// ID returns the unique identifier for ast. Within a Context, two
// ASTs have the same ID if and only if they are Equal.
func (ast AST) ID() uint64 {
	var res uint64
	ast.ctx.do(func() {
		res = uint64(C.Z3_get_ast_id(ast.ctx.c, ast.c))
	})
	runtime.KeepAlive(ast)
	return res
}

// TODO: AST.AsValue, etc.
