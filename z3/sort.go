// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import (
	"runtime"
)

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"

// Sort represents the type of the Expr. A Sort can be a basic type
// such as int or bool or a parameterized type such as a 30 bit wide
// bit-vector or an array from int to bool.
type Sort struct {
	*sortImpl
	noEq
}

// sortImpl wraps the underlying C.Z3_sort. This is separate from Sort
// so a finalizer can be attached to this without exposing it to the
// user.
type sortImpl struct {
	ctx  *Context
	c    C.Z3_sort
	kind Kind
}

func wrapSort(ctx *Context, c C.Z3_sort, kind Kind) Sort {
	if kind == KindUnknown {
		ctx.do(func() {
			kind = Kind(C.Z3_get_sort_kind(ctx.c, c))
		})
	}
	impl := &sortImpl{ctx, c, kind}
	ctx.lock.Lock()
	C.Z3_inc_ref(ctx.c, C.Z3_sort_to_ast(ctx.c, c))
	ctx.lock.Unlock()
	runtime.SetFinalizer(impl, func(impl *sortImpl) {
		impl.ctx.do(func() {
			C.Z3_dec_ref(impl.ctx.c, C.Z3_sort_to_ast(impl.ctx.c, impl.c))
		})
	})
	return Sort{impl, noEq{}}
}

// String returns a string representation of s.
func (s Sort) String() string {
	var res string
	s.ctx.do(func() {
		res = C.GoString(C.Z3_sort_to_string(s.ctx.c, s.c))
	})
	runtime.KeepAlive(s)
	return res
}

// Kind returns s's kind.
func (s Sort) Kind() Kind {
	return s.kind
}

// BVSize returns the bit size of a bit-vector sort.
func (s Sort) BVSize() int {
	var size int
	s.ctx.do(func() {
		size = int(C.Z3_get_bv_sort_size(s.ctx.c, s.c))
	})
	runtime.KeepAlive(s)
	return size
}

// Equal returns true if s and o are structurally identical.
func (s Sort) Equal(o Sort) bool {
	var out bool
	s.ctx.do(func() {
		out = z3ToBool(C.Z3_is_eq_sort(s.ctx.c, s.c, o.c))
	})
	runtime.KeepAlive(s)
	runtime.KeepAlive(o)
	return out
}
