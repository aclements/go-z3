// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import "testing"

func TestASTEquality(t *testing.T) {
	ctx := NewContext(nil)
	ints := ctx.IntSort()
	x := ctx.Const("x", ints).(Int)
	a := x.Add(ctx.FromInt(2, ints).(Int))
	b := x.Add(ctx.FromInt(1, ints).(Int)).Add(ctx.FromInt(1, ints).(Int))
	as, bs := ctx.Simplify(a, nil).AsAST(), ctx.Simplify(b, nil).AsAST()
	if !as.Equal(bs) {
		t.Errorf("%v != %v", as, bs)
	}
	if h1, h2 := as.Hash(), bs.Hash(); h1 != h2 {
		t.Errorf("hashes differ: %v != %v", h1, h2)
	}
	if i1, i2 := as.ID(), bs.ID(); i1 != i2 {
		t.Errorf("IDs differ: %v != %v", i1, i2)
	}
}
