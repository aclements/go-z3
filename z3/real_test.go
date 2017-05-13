// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import (
	"math/big"
	"testing"
)

func TestRealRational(t *testing.T) {
	ctx := NewContext(nil)
	rat := ctx.FromBigRat(big.NewRat(5, 4))
	numer, denom, isLit := rat.AsRat()
	if !isLit {
		t.Errorf("(%s).AsRat() returned false", rat)
	} else {
		val, isLit, ok := numer.AsInt64()
		if !(val == 5 && isLit && ok) {
			t.Errorf("numerator of %s: wanted 5, true, true; got %v, %v, %v", rat, val, isLit, ok)
		}
		val, isLit, ok = denom.AsInt64()
		if !(val == 4 && isLit && ok) {
			t.Errorf("numerator of %s: wanted 4, true, true; got %v, %v, %v", rat, val, isLit, ok)
		}
	}

	_, _, isLit = rat.Approx(10)
	if isLit {
		// rat is a rational, so Approx should fail.
		t.Errorf("(%s).Approx(10) returned true", rat)
	}
}

func TestRealIrrational(t *testing.T) {
	t.Skipf("need a simplifier") // TODO

	ctx := NewContext(nil)
	root2 := ctx.FromInt(2, ctx.IntSort()).(Int).ToReal().Exp(ctx.FromBigRat(big.NewRat(1, 2)))
	_, _, isLit := root2.AsRat()
	if isLit {
		t.Errorf("(%s).AsRat() returned true", root2)
	}

	l, r, isLit := root2.Approx(10)
	if !isLit {
		t.Errorf("(%s).Approx(10) returned false", root2)
	} else {
		// TODO: Check l and r.
		t.Logf("[%v, %v]", l, r)
	}
}
