// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
*/
import "C"

// TODO: Implement simplifier options.

// TODO: Should this be over ASTs?

type futureParams struct{}

// Simplify simplifies expression x.
//
// The second argument provides parameters to the simplifier. For now
// it must be nil, which means to use the default parameters.
//
// The resulting expression will have the same sort as x.
func (ctx *Context) Simplify(x Value, _ *futureParams) Value {
	return wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_simplify(ctx.c, x.impl().c)
	}).lift(KindUnknown)
}
