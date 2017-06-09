// Generated by genwrap.go. DO NOT EDIT

package z3

import "runtime"

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"

// Eq returns a Value that is true if l and r are equal.
func (l BV) Eq(r BV) Bool {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_eq(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// NE returns a Value that is true if l and r are not equal.
func (l BV) NE(r BV) Bool {
	return l.ctx.Distinct(l, r)
}

// Not returns the bit-wise negation of l.
func (l BV) Not() BV {
	// Generated from bv.go:117.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvnot(ctx.c, l.c)
	})
	runtime.KeepAlive(l)
	return BV(val)
}

// AllBits returns a 1-bit bit-vector that is the bit-wise "and" of
// all bits.
func (l BV) AllBits() BV {
	// Generated from bv.go:122.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvredand(ctx.c, l.c)
	})
	runtime.KeepAlive(l)
	return BV(val)
}

// AnyBits returns a 1-bit bit-vector that is the bit-wise "or" of all
// bits.
func (l BV) AnyBits() BV {
	// Generated from bv.go:127.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvredor(ctx.c, l.c)
	})
	runtime.KeepAlive(l)
	return BV(val)
}

// And returns the bit-wise and of l and r.
//
// l and r must have the same size.
func (l BV) And(r BV) BV {
	// Generated from bv.go:133.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvand(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return BV(val)
}

// Or returns the bit-wise or of l and r.
//
// l and r must have the same size.
func (l BV) Or(r BV) BV {
	// Generated from bv.go:139.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvor(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return BV(val)
}

// Xor returns the bit-wise xor of l and r.
//
// l and r must have the same size.
func (l BV) Xor(r BV) BV {
	// Generated from bv.go:145.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvxor(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return BV(val)
}

// Nand returns the bit-wise nand of l and r.
//
// l and r must have the same size.
func (l BV) Nand(r BV) BV {
	// Generated from bv.go:151.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvnand(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return BV(val)
}

// Nor returns the bit-wise nor of l and r.
//
// l and r must have the same size.
func (l BV) Nor(r BV) BV {
	// Generated from bv.go:157.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvnor(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return BV(val)
}

// Xnor returns the bit-wise xnor of l and r.
//
// l and r must have the same size.
func (l BV) Xnor(r BV) BV {
	// Generated from bv.go:163.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvxnor(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return BV(val)
}

// Neg returns the two's complement negation of l.
func (l BV) Neg() BV {
	// Generated from bv.go:167.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvneg(ctx.c, l.c)
	})
	runtime.KeepAlive(l)
	return BV(val)
}

// Add returns the two's complement sum of l and r.
//
// l and r must have the same size.
func (l BV) Add(r BV) BV {
	// Generated from bv.go:173.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvadd(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return BV(val)
}

// Sub returns the two's complement subtraction l minus r.
//
// l and r must have the same size.
func (l BV) Sub(r BV) BV {
	// Generated from bv.go:179.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvsub(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return BV(val)
}

// Mul returns the two's complement product of l and r.
//
// l and r must have the same size.
func (l BV) Mul(r BV) BV {
	// Generated from bv.go:185.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvmul(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return BV(val)
}

// UDiv returns the floor of l / r, treating l and r as unsigned.
//
// If r is 0, the result is unconstrained.
//
// l and r must have the same size.
func (l BV) UDiv(r BV) BV {
	// Generated from bv.go:191.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvudiv(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return BV(val)
}

// SDiv returns l / r rounded toward 0, treating l and r as two's
// complement signed numbers.
//
// If r is 0, the result is unconstrained.
//
// l and r must have the same size.
func (l BV) SDiv(r BV) BV {
	// Generated from bv.go:197.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvsdiv(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return BV(val)
}

// URem returns the unsigned remainder of l divided by r.
//
// l and r must have the same size.
func (l BV) URem(r BV) BV {
	// Generated from bv.go:203.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvurem(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return BV(val)
}

// SRem returns the two's complement signed remainder of l divided by r.
//
// The sign of the result follows the sign of l.
//
// l and r must have the same size.
func (l BV) SRem(r BV) BV {
	// Generated from bv.go:211.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvsrem(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return BV(val)
}

// SMod returns the two's complement signed modulus of l divided by r.
//
// The sign of the result follows the sign of r.
//
// l and r must have the same size.
func (l BV) SMod(r BV) BV {
	// Generated from bv.go:219.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvsmod(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return BV(val)
}

// ULT returns the l < r, where l and r are unsigned.
//
// l and r must have the same size.
func (l BV) ULT(r BV) Bool {
	// Generated from bv.go:225.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvult(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// SLT returns the l < r, where l and r are signed.
//
// l and r must have the same size.
func (l BV) SLT(r BV) Bool {
	// Generated from bv.go:231.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvslt(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// ULE returns the l <= r, where l and r are unsigned.
//
// l and r must have the same size.
func (l BV) ULE(r BV) Bool {
	// Generated from bv.go:237.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvule(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// SLE returns the l <= r, where l and r are signed.
//
// l and r must have the same size.
func (l BV) SLE(r BV) Bool {
	// Generated from bv.go:243.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvsle(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// UGE returns the l >= r, where l and r are unsigned.
//
// l and r must have the same size.
func (l BV) UGE(r BV) Bool {
	// Generated from bv.go:249.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvuge(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// SGE returns the l >= r, where l and r are signed.
//
// l and r must have the same size.
func (l BV) SGE(r BV) Bool {
	// Generated from bv.go:255.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvsge(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// UGT returns the l > r, where l and r are unsigned.
//
// l and r must have the same size.
func (l BV) UGT(r BV) Bool {
	// Generated from bv.go:261.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvugt(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// SGT returns the l > r, where l and r are signed.
//
// l and r must have the same size.
func (l BV) SGT(r BV) Bool {
	// Generated from bv.go:267.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvsgt(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// Concat returns concatenation of l and r.
//
// The result is a bit-vector whose length is the sum of the lengths
// of l and r.
func (l BV) Concat(r BV) BV {
	// Generated from bv.go:274.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_concat(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return BV(val)
}

// Extract returns bits [high, low] (inclusive) of l, where bit 0 is
// the least significant bit.
func (l BV) Extract(high int, low int) BV {
	// Generated from bv.go:279.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_extract(ctx.c, C.unsigned(high), C.unsigned(low), l.c)
	})
	runtime.KeepAlive(l)
	return BV(val)
}

// SignExtend returns l sign-extended to a bit-vector of length m+i,
// where m is the length of l.
func (l BV) SignExtend(i int) BV {
	// Generated from bv.go:284.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_sign_ext(ctx.c, C.unsigned(i), l.c)
	})
	runtime.KeepAlive(l)
	return BV(val)
}

// ZeroExtend returns l zero-extended to a bit-vector of length m+i,
// where m is the length of l.
func (l BV) ZeroExtend(i int) BV {
	// Generated from bv.go:289.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_zero_ext(ctx.c, C.unsigned(i), l.c)
	})
	runtime.KeepAlive(l)
	return BV(val)
}

// Repeat returns l repeated up to length i.
func (l BV) Repeat(i int) BV {
	// Generated from bv.go:293.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_repeat(ctx.c, C.unsigned(i), l.c)
	})
	runtime.KeepAlive(l)
	return BV(val)
}

// Lsh returns l shifted left by i bits.
//
// This is equivalent to l * 2^i.
//
// l and i must have the same size. The result has the same sort.
func (l BV) Lsh(i BV) BV {
	// Generated from bv.go:301.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvshl(ctx.c, l.c, i.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(i)
	return BV(val)
}

// URsh returns l logically shifted right by i bits.
//
// This is equivalent to l / 2^i, where l and i are unsigned.
//
// l and i must have the same size. The result has the same sort.
func (l BV) URsh(i BV) BV {
	// Generated from bv.go:309.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvlshr(ctx.c, l.c, i.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(i)
	return BV(val)
}

// SRsh returns l arithmetically shifted right by i bits.
//
// This is like URsh, but the sign of the result is the sign of l.
//
// l and i must have the same size. The result has the same sort.
func (l BV) SRsh(i BV) BV {
	// Generated from bv.go:317.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bvashr(ctx.c, l.c, i.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(i)
	return BV(val)
}

// RotateLeft returns l rotated left by i bits.
//
// l and i must have the same size.
func (l BV) RotateLeft(i BV) BV {
	// Generated from bv.go:323.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_ext_rotate_left(ctx.c, l.c, i.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(i)
	return BV(val)
}

// RotateRight returns l rotated right by i bits.
//
// l and i must have the same size.
func (l BV) RotateRight(i BV) BV {
	// Generated from bv.go:329.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_ext_rotate_right(ctx.c, l.c, i.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(i)
	return BV(val)
}

// SToInt converts signed bit-vector l to an integer.
func (l BV) SToInt() Int {
	// Generated from bv.go:333.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bv2int(ctx.c, l.c, C.Z3_TRUE)
	})
	runtime.KeepAlive(l)
	return Int(val)
}

// UToInt converts unsigned bit-vector l to an integer.
func (l BV) UToInt() Int {
	// Generated from bv.go:337.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_bv2int(ctx.c, l.c, C.Z3_FALSE)
	})
	runtime.KeepAlive(l)
	return Int(val)
}
