// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import (
	"runtime"
	"strconv"
)

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"

// TODO: This should follow Expr and be safely copyable and use "Sort"
// instead of "*Sort" in the API.

// Sort represents the type of the Expr. A Sort can be a basic type
// such as int or bool or a parameterized type such as a 30 bit wide
// bit-vector or an array from int to bool.
type Sort struct {
	ctx  *Context
	c    C.Z3_sort
	kind SortKind
	noEq
}

// sortWrappers is a map of Expr constructors for each sort kind.
var sortWrappers = make(map[SortKind]func(x expr) Expr)

func wrapSort(ctx *Context, c C.Z3_sort, kind SortKind) *Sort {
	if kind == SortUnknown {
		ctx.do(func() {
			kind = SortKind(C.Z3_get_sort_kind(ctx.c, c))
		})
	}
	sort := &Sort{ctx, c, kind, noEq{}}
	ctx.lock.Lock()
	C.Z3_inc_ref(ctx.c, C.Z3_sort_to_ast(ctx.c, c))
	ctx.lock.Unlock()
	// TODO: Don't put finalizer on a user-accessible type.
	runtime.SetFinalizer(sort, func(sort *Sort) {
		sort.ctx.do(func() {
			C.Z3_dec_ref(sort.ctx.c, C.Z3_sort_to_ast(sort.ctx.c, sort.c))
		})
	})
	return sort
}

// String returns a string representation of sort.
func (sort *Sort) String() string {
	var res string
	sort.ctx.do(func() {
		res = C.GoString(C.Z3_sort_to_string(sort.ctx.c, sort.c))
	})
	runtime.KeepAlive(sort)
	return res
}

// SortKind is a general category of sort, such as int or array.
type SortKind int

const (
	SortUninterpreted = SortKind(C.Z3_UNINTERPRETED_SORT)
	SortBool          = SortKind(C.Z3_BOOL_SORT)
	SortInt           = SortKind(C.Z3_INT_SORT)
	SortReal          = SortKind(C.Z3_REAL_SORT)
	SortBV            = SortKind(C.Z3_BV_SORT)
	SortArray         = SortKind(C.Z3_ARRAY_SORT)
	SortDatatype      = SortKind(C.Z3_DATATYPE_SORT)
	SortRelation      = SortKind(C.Z3_RELATION_SORT)
	SortFiniteDomain  = SortKind(C.Z3_FINITE_DOMAIN_SORT)
	SortFloatingPoint = SortKind(C.Z3_FLOATING_POINT_SORT)
	SortRoundingMode  = SortKind(C.Z3_ROUNDING_MODE_SORT)
	SortUnknown       = SortKind(C.Z3_UNKNOWN_SORT)
)

// String returns k as a string like "SortBool".
func (k SortKind) String() string {
	switch k {
	case SortUninterpreted:
		return "SortUninterpreted"
	case SortBool:
		return "SortBool"
	case SortInt:
		return "SortInt"
	case SortReal:
		return "SortReal"
	case SortBV:
		return "SortBV"
	case SortArray:
		return "SortArray"
	case SortDatatype:
		return "SortDatatype"
	case SortRelation:
		return "SortRelation"
	case SortFiniteDomain:
		return "SortFiniteDomain"
	case SortFloatingPoint:
		return "SortFloatingPoint"
	case SortRoundingMode:
		return "SortRoundingMode"
	case SortUnknown:
		return "SortUnknown"
	}
	return "SortKind(" + strconv.Itoa(int(k)) + ")"
}

// Kind returns s's kind.
func (s *Sort) Kind() SortKind {
	return s.kind
}

// BVSize returns the bit size of a bit-vector sort.
func (s *Sort) BVSize() int {
	var size int
	s.ctx.do(func() {
		size = int(C.Z3_get_bv_sort_size(s.ctx.c, s.c))
	})
	runtime.KeepAlive(s)
	return size
}

// Equal returns true if s and o are structurally identical.
func (s *Sort) Equal(o *Sort) bool {
	var out bool
	s.ctx.do(func() {
		out = z3ToBool(C.Z3_is_eq_sort(s.ctx.c, s.c, o.c))
	})
	runtime.KeepAlive(s)
	runtime.KeepAlive(o)
	return out
}
