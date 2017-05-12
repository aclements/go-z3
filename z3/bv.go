// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import (
	"math"
	"math/big"
)

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"

// BVSort returns a bit-vector sort (type) of the given width in bits.
func (ctx *Context) BVSort(bits int) *Sort {
	var csort C.Z3_sort
	ctx.do(func() {
		csort = C.Z3_mk_bv_sort(ctx.c, C.unsigned(bits))
	})
	return wrapSort(ctx, csort)
}

// BVConst returns a bit-vector constant named "name" with the given
// width in bits.
func (ctx *Context) BVConst(name string, bits int) *Expr {
	return ctx.Const(name, ctx.BVSort(bits))
}

// BVSAsBig returns the value of expr as a math/big.Int, interpreting
// expr as a signed two's complement number. expr must be a
// bit-vector. If expr is a bit-vector, but not constant, this returns
// nil.
func (expr *Expr) BVSAsBig() *big.Int {
	v := expr.BVUAsBig()
	if v == nil {
		return nil
	}
	size := expr.Sort().BVSize()
	if v.Bit(size-1) != 0 {
		shift := big.NewInt(1)
		shift.Lsh(shift, uint(size))
		v.Sub(v, shift)
	}
	return v
}

// BVUAsBig is like BVSAsBig, but interprets expr as unsigned.
func (expr *Expr) BVUAsBig() *big.Int {
	if expr.Sort().Kind() != SortBV {
		panic("not a bit-vector")
	}
	if expr.astKind() != C.Z3_NUMERAL_AST {
		return nil
	}
	var str string
	expr.ctx.do(func() {
		cstr := C.Z3_get_numeral_string(expr.ctx.c, expr.c)
		str = C.GoString(cstr)
	})
	var v big.Int
	if _, ok := v.SetString(str, 10); !ok {
		panic("failed to parse numeral string")
	}
	return &v
}

// BVAsInt returns the value of expr as an int64, interpreting expr as
// a two's complement signed number. expr must be a bit-vector. If
// expr is a bit-vector, but not constant, it returns 0, false, false.
// If expr is a constant bit-vector, but its value cannot be
// represented as an int64, it returns 0, true, false.
func (expr *Expr) BVAsInt() (val int64, isConst, ok bool) {
	// Oddly, Z3_get_numeral_int64 interprets the number as
	// unsigned, which makes no sense since the API is strictly
	// less useful than Z3_get_numeral_uint64 and doesn't mirror
	// Z3_mk_int64. So, use Z3_get_numeral_uint64 and sign extend
	// it ourselves.
	uval, isConst, ok := expr.BVAsUint()
	if !isConst {
		return 0, isConst, ok
	}
	size := expr.Sort().BVSize()
	if ok && size < 64 {
		// Fits in an int64 regardless of sign. Sign-extend it.
		return int64(uval) << uint(64-size) >> uint(64-size), true, true
	}
	// size is >= 64, so we have to tread carefully.
	if ok && uval < 1<<63 {
		// Positive and fits in an int64.
		return int64(uval), true, true
	}
	// It may have overflowed uint64 just because of sign bits.
	// Take the slow path.
	bigVal := expr.BVSAsBig()
	if bigVal.Cmp(big.NewInt(math.MaxInt64)) > 0 {
		return 0, true, false
	}
	if bigVal.Cmp(big.NewInt(math.MinInt64)) < 0 {
		return 0, true, false
	}
	return bigVal.Int64(), true, true
}

// BVAsUint is like BVAsInt, but interprets expr as unsigned and fails
// if expr cannot be represented as a uint64.
func (expr *Expr) BVAsUint() (val uint64, isConst, ok bool) {
	if expr.Sort().Kind() != SortBV {
		panic("not a bit-vector")
	}
	if expr.astKind() != C.Z3_NUMERAL_AST {
		return 0, false, false
	}
	var cval C.__uint64
	expr.ctx.do(func() {
		ok = z3ToBool(C.Z3_get_numeral_uint64(expr.ctx.c, expr.c, &cval))
	})
	return uint64(cval), true, ok
}

//go:generate go run genwrap.go -- $GOFILE

// BVNot returns the bit-wise negation of l.
//
// l must have bit-vector sort.
//
//wrap:expr BVNot Z3_mk_bvnot l

// BVAnd returns the bit-wise and of l and r.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVAnd Z3_mk_bvand l r

// BVOr returns the bit-wise or of l and r.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVOr Z3_mk_bvor l r

// BVXor returns the bit-wise xor of l and r.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVXor Z3_mk_bvxor l r

// BVNand returns the bit-wise nand of l and r.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVNand Z3_mk_bvnand l r

// BVNor returns the bit-wise nor of l and r.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVNor Z3_mk_bvnor l r

// BVXnor returns the bit-wise xnor of l and r.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVXnor Z3_mk_bvxnor l r

// BVNeg returns the two's complement negation of l.
//
// l must have bit-vector sort.
//
//wrap:expr BVNeg Z3_mk_bvneg l

// BVAdd returns the two's complement sum of l and r.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVAdd Z3_mk_bvadd l r

// BVSub returns the two's complement subtraction l minus r.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVSub Z3_mk_bvsub l r

// BVMul returns the two's complement product of l and r.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVMul Z3_mk_bvmul l r

// BVUDiv returns the unsigned division of l over r.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVUDiv Z3_mk_bvudiv l r

// BVSDiv returns the two's complement signed division of l over r.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVSDiv Z3_mk_bvsdiv l r

// BVURem returns the unsigned remainder of l divided by r.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVURem Z3_mk_bvurem l r

// BVSRem returns the two's complement signed remainder of l divided by r.
//
// The sign of the result follows the sign of l.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVSRem Z3_mk_bvsrem l r

// BVSMod returns the two's complement signed modulus of l divided by r.
//
// The sign of the result follows the sign of r.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVSMod Z3_mk_bvsmod l r

// BVULT returns the l < r, where l and r are unsigned.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVULT Z3_mk_bvult l r

// BVSLT returns the l < r, where l and r are signed.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVSLT Z3_mk_bvslt l r

// BVULE returns the l <= r, where l and r are unsigned.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVULE Z3_mk_bvule l r

// BVSLE returns the l <= r, where l and r are signed.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVSLE Z3_mk_bvsle l r

// BVUGE returns the l >= r, where l and r are unsigned.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVUGE Z3_mk_bvuge l r

// BVSGE returns the l >= r, where l and r are signed.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVSGE Z3_mk_bvsge l r

// BVUGT returns the l > r, where l and r are unsigned.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVUGT Z3_mk_bvugt l r

// BVSGT returns the l > r, where l and r are signed.
//
// l and r must have the same bit-vector sort.
//
//wrap:expr BVSGT Z3_mk_bvsgt l r

// BVConcat returns concatenation of l and r.
//
// l and r must have bit-vector sort. The result is a bit-vector whose
// length is the sum of the lengths of l and r.
//
//wrap:expr BVConcat Z3_mk_concat l r

// BVExtract returns bits [high, low] (inclusive) of l, where bit 0 is
// the least significant bit.
//
// l must have bit-vector sort.
//
//wrap:expr BVExtract l high:int low:int : Z3_mk_extract high:unsigned low:unsigned l

// BVSignExtend returns l sign-extended to a bit-vector of length m+i,
// where m is the length of l.
//
// l must have bit-vector sort.
//
//wrap:expr BVSignExtend l i:int : Z3_mk_sign_ext i:unsigned l

// BVZeroExtend returns l zero-extended to a bit-vector of length m+i,
// where m is the length of l.
//
// l must have bit-vector sort.
//
//wrap:expr BVZeroExtend l i:int : Z3_mk_zero_ext i:unsigned l

// BVRepeat returns l repeated up to length i.
//
// l must have bit-vector sort.
//
//wrap:expr BVRepeat l i:int : Z3_mk_repeat i:unsigned l

// BVShiftLeft returns l shifted left by i bits.
//
// This is equivalent to l * 2^i.
//
// l and i must have the same bit-vector sort. The result has the same
// sort.
//
//wrap:expr BVShiftLeft Z3_mk_bvshl l i

// BVShiftRightLogical returns l logically shifted right by i bits.
//
// This is equivalent to l / 2^i, where l and i are unsigned.
//
// l and i must have the same bit-vector sort. The result has the same
// sort.
//
//wrap:expr BVShiftRightLogical Z3_mk_bvlshr l i

// BVShiftRightArithmetic returns l arithmetically shifted right by i bits.
//
// This is like BVShiftRightLogical, but the sign of the result is the
// sign of l.
//
// l and i must have the same bit-vector sort. The result has the same
// sort.
//
//wrap:expr BVShiftRightArithmetic Z3_mk_bvashr l i

// BVRotateLeft returns l rotated left by i bits.
//
// l and i must have the same bit-vector sort.
//
//wrap:expr BVRotateLeft Z3_mk_ext_rotate_left l i

// BVRotateRight returns l rotated right by i bits.
//
// l and i must have the same bit-vector sort.
//
//wrap:expr BVRotateRight Z3_mk_ext_rotate_right l i

// Int2BV converts integer l to a bit-vector of width bits.
//
// l must have integer sort.
//
//wrap:expr Int2BV l bits:int : Z3_mk_int2bv bits:unsigned l

// BVS2Int converts signed bit-vector l to an integer.
//
// l must have bit-vector sort.
//
//wrap:expr BVS2Int l : Z3_mk_bv2int l 1:Z3_bool

// BVU2Int converts unsigned bit-vector l to an integer.
//
// l must have bit-vector sort.
//
//wrap:expr BVU2Int l : Z3_mk_bv2int l 0:Z3_bool

// TODO: Z3_mk_bv*_no_{over,under}flow
