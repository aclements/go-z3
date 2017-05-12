// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package z3 checks the satisfiability of logical formulas.
//
// This package provides bindings for the Z3 SMT solver
// (https://github.com/Z3Prover/z3). Z3 checks satisfiability of
// logical formulas over a wide range of terms, including booleans,
// integers, reals, bit-vectors, and uninterpreted functions. For a
// good introduction to the concepts of SMT and Z3, see the Z3 guide
// (http://rise4fun.com/z3/tutorialcontent/guide).
//
// Currently this package only supports formulas of booleans and
// bit-vectors, though more types are easy to add.
//
// The main entry point to the z3 package is type Context. All
// expressions are created and all solving is done relative to some
// Context, and expressions from different Contexts cannot be mixed.
//
// Expressions are represented by type Expr. Every expression has a
// type, called a "sort" and represented by type Sort. Type checking
// of expressions is done dynamically by Z3 and creating a badly typed
// expression (an expression that is not "well sorted") will result in
// a panic. An expression of boolean sort is called a "formula".
//
// Terms in expressions can be "numerals", "constants", or
// "uninterpreted functions". A numeral is a literal, fixed value like
// "2". A constant is a term like "x", whose value is fixed but
// unspecified. An uninterpreted function is a function whose mapping
// from arguments to results is fixed but unspecified (this is in
// contrast to an "interpreted function" like + whose interpretation
// is specified to be addition). Functions are pure (side-effect-free)
// like mathematical functions, but unlike mathematical functions they
// are always total. A constant can be thought of as a function with
// zero arguments.
//
// Type Solver checks the satisfiability of a set of formulas. If the
// Solver determines that a set of formulas is satisfiable, it can
// construct a Model giving a specific assignment of constants and
// uninterpreted functions that satisfies the set formulas.
package z3

/*
#include <z3.h>
*/
import "C"

func boolToZ3(b bool) C.Z3_bool {
	if b {
		return C.Z3_TRUE
	}
	return C.Z3_FALSE
}

func z3ToBool(b C.Z3_bool) bool {
	return b != C.Z3_FALSE
}
