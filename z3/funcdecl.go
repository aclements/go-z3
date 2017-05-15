// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import (
	"runtime"
	"unsafe"
)

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"

// FuncDecl is a function declaration.
//
// A FuncDecl can represent either a interpreted function like "+" or
// an uninterpreted function created by Context.FuncDecl.
//
// A FuncDecl can be used in an expression by applying it to a set of
// arguments.
type FuncDecl struct {
	*funcDeclImpl
	noEq
}

type funcDeclImpl struct {
	ctx *Context
	c   C.Z3_func_decl
}

func wrapFuncDecl(ctx *Context, c C.Z3_func_decl) FuncDecl {
	impl := &funcDeclImpl{ctx, c}
	ctx.lock.Lock()
	C.Z3_inc_ref(ctx.c, C.Z3_func_decl_to_ast(ctx.c, c))
	ctx.lock.Unlock()
	runtime.SetFinalizer(impl, func(impl *funcDeclImpl) {
		impl.ctx.do(func() {
			C.Z3_dec_ref(impl.ctx.c, C.Z3_func_decl_to_ast(impl.ctx.c, impl.c))
		})
	})
	runtime.KeepAlive(ctx)
	return FuncDecl{impl, noEq{}}
}

// FuncDecl creates an uninterpreted function named "name".
//
// In contrast with an interpreted function like "+", an uninterpreted
// function is only assigned an interpretation in a particular model,
// and different models may assign different interpretations.
func (ctx *Context) FuncDecl(name string, domain []Sort, range_ Sort) FuncDecl {
	sym := ctx.symbol(name)
	cdomain := make([]C.Z3_sort, len(domain))
	for i, sort := range domain {
		cdomain[i] = sort.c
	}
	var cfuncdecl C.Z3_func_decl
	ctx.do(func() {
		cfuncdecl = C.Z3_mk_func_decl(ctx.c, sym, C.uint(len(cdomain)), &cdomain[0], range_.c)
	})
	runtime.KeepAlive(&domain[0])
	runtime.KeepAlive(range_)
	return wrapFuncDecl(ctx, cfuncdecl)
}

// FreshFuncDecl creates a fresh uninterpreted function distinct from
// all other functions.
func (ctx *Context) FreshFuncDecl(prefix string, domain []Sort, range_ Sort) FuncDecl {
	cprefix := C.CString(prefix)
	defer C.free(unsafe.Pointer(cprefix))
	cdomain := make([]C.Z3_sort, len(domain))
	for i, sort := range domain {
		cdomain[i] = sort.c
	}
	var cfuncdecl C.Z3_func_decl
	ctx.do(func() {
		cfuncdecl = C.Z3_mk_fresh_func_decl(ctx.c, cprefix, C.uint(len(cdomain)), &cdomain[0], range_.c)
	})
	runtime.KeepAlive(&domain[0])
	runtime.KeepAlive(range_)
	return wrapFuncDecl(ctx, cfuncdecl)
}

// String returns a string representation of f.
func (f FuncDecl) String() string {
	var res string
	f.ctx.do(func() {
		res = C.GoString(C.Z3_func_decl_to_string(f.ctx.c, f.c))
	})
	runtime.KeepAlive(f)
	return res
}

// AsAST returns the AST representation of f.
func (f FuncDecl) AsAST() AST {
	var cast C.Z3_ast
	f.ctx.do(func() {
		cast = C.Z3_func_decl_to_ast(f.ctx.c, f.c)
	})
	runtime.KeepAlive(f)
	return wrapAST(f.ctx, cast)
}

// Apply creates an expression that applies function f to arguments
// args.
//
// The sorts of args must be the domain of f. The sort of the
// resulting expression will be f's range.
func (f FuncDecl) Apply(args ...Expr) Expr {
	cargs := make([]C.Z3_ast, len(args))
	for i, arg := range args {
		cargs[i] = arg.impl().c
	}
	var cexpr C.Z3_ast
	f.ctx.do(func() {
		cexpr = C.Z3_mk_app(f.ctx.c, f.c, C.uint(len(cargs)), &cargs[0])
	})
	runtime.KeepAlive(&cargs[0])
	return wrapExpr(f.ctx, cexpr).lift(KindUnknown)
}

// TODO: Lots of accessors
