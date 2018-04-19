package z3

import (
	"testing"
)

func TestOptimize(t *testing.T) {
	ctx := NewContext(nil)
	opt := NewOptimize(ctx)
	opt.Set("priority", "pareto")

	x := ctx.IntConst("x")
	y := ctx.IntConst("y")
	zero := ctx.FromInt(0, ctx.IntSort()).(Int)
	ten := ctx.FromInt(10, ctx.IntSort()).(Int)
	eleven := ctx.FromInt(11, ctx.IntSort()).(Int)

	opt.Assert(ten.GE(x).And(x.GE(zero)))
	opt.Assert(ten.GE(y).And(y.GE(zero)))
	opt.Assert(x.Add(y).LE(eleven))

	h1 := opt.Maximize(x)
	h2 := opt.Maximize(y)

	const TotalSolutions = 10
	var solutions int
	for {
		if sat, err := opt.Check(); sat {
			t.Log("x: ", h1.Lower(), ", y: ", h2.Lower())
			solutions++
		} else if err != nil {
			t.Fatalf("error: %s", err)
		} else if solutions > TotalSolutions {
			t.Fatalf("Too many solutions found (expected %d, found %d)\n",
				TotalSolutions, solutions)
		} else {
			break
		}
	}
}

// Based on an example from the z3 optimization tutorial at
// https://rise4fun.com/Z3/tutorial/optimization
func TestOptimizeSoft(t *testing.T) {
	ctx := NewContext(nil)
	opt := NewOptimize(ctx)

	a := ctx.BoolConst("a")
	b := ctx.BoolConst("b")
	c := ctx.BoolConst("c")

	opt.AssertSoft(a, Weight(1), ID("A"))
	opt.AssertSoft(b, Weight(2), ID("B"))
	opt.AssertSoft(c, Weight(3), ID("A"))
	opt.Assert(a.Eq(c))
	opt.Assert(a.And(b).Not())

	if sat, err := opt.Check(); sat {
		model := opt.Model()
		if val, _ := model.Eval(c, false).(Bool).AsBool(); !val {
			t.Fatal("c has wrong value")
		}
		if val, _ := model.Eval(b, false).(Bool).AsBool(); val {
			t.Fatal("b has wrong value")
		}
		if val, _ := model.Eval(a, false).(Bool).AsBool(); !val {
			t.Fatal("a has wrong value")
		}
	} else if err != nil {
		t.Fatalf("error: %s", err)
	}
}
