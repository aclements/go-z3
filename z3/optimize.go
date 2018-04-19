package z3

import (
	"fmt"
	"runtime"
	"strconv"
)

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
*/
import "C"

// An Optimize object is a collection of predicates that can be queried for
// satisfying maximums or minimums.
//
// These predicates form a stack that can be manipulated with
// Push/Pop.
type Optimize struct {
	*optimizeImpl
}

type optimizeImpl struct {
	ctx *Context
	o   C.Z3_optimize
}

// NewOptimize returns a new, empty optimize object
func NewOptimize(ctx *Context) *Optimize {
	var impl *optimizeImpl
	ctx.do(func() {
		impl = &optimizeImpl{
			ctx,
			C.Z3_mk_optimize(ctx.c),
		}
	})
	ctx.do(func() {
		C.Z3_optimize_inc_ref(ctx.c, impl.o)
	})
	runtime.SetFinalizer(impl, func(impl *optimizeImpl) {
		impl.ctx.do(func() {
			C.Z3_optimize_dec_ref(impl.ctx.c, impl.o)
		})
	})
	return &Optimize{impl}
}

// Z3 optimize config properties as of 4.6.0
func newOptimizeConfig() *Config {
	return newConfig([]param{
		{"dump_benchmarks", "bool", "dump benchmarks for profiling (default: false)"},
		{"elim_01", "bool", "eliminate 01 variables (default: true)"},
		{"enable_sat", "bool", "enable the new SAT core for propositional constraints (default: true)"},
		{"enable_sls", "bool", "enable SLS tuning during weighted maxsast (default: false)"},
		{"maxres.add_upper_bound_block", "bool", "restrict upper bound with constraint (default: false)"},
		{"maxres.hill_climb", "bool", "give preference for large weight cores (default: true)"},
		{"maxres.max_core_size", "uint", "break batch of generated cores if size reaches this number (default 3)"},
		{"maxres.max_correction_set_size", "uint", "allow generating correction set constraints up to maximal size (default: 3)"},
		{"maxres.max_num_cores", "uint", "maximumal number of cores per round (default: 4294967295)"},
		{"maxres.maximize_assignment", "bool", "find an MSS/MCS to improve current assignment (default: false)"},
		{"maxres.pivot_on_correection", "bool", "reduce soft constraints if the current correction set is smaller than current core (default: true)"},
		{"maxres.wmax", "bool", "use weighted theory solver to constrain upper bounds (default: false)"},
		{"maxsat_engine", "string", "select engine for maxsat: 'core_maxsat', 'wmax', 'axres', 'pd-maxres' (default: maxres)"},
		{"optsmt_engine", "string", "select optimization engine: 'basic', 'farkas', 'symba' (default: basic)"},
		{"pb.compile_equality", "bool", "compile arithmetical equalities into pseudo-Boolean equality (instead of two equalities) (default: false)"},
		{"pb.neat", "bool", "use neat (as opposed to less readable, but faster) pretty printer when displaying context (default: true)"},
		{"priority", "string", "select how to prioritize objectives: 'lex' (lexicographic), 'pareot', or 'box' (default: lex)"},
		{"rlimit", "uint", "resource limit (0 means no limit) (default: 0)"},
		{"timout", "uint", "timeout (in milliseconds) (UINT_MAX and 0 mean no timeout) (default 4294967295)"},
	})
}

// Set an option on the optimize object
func (o *Optimize) Set(name string, val interface{}) {
	config := newOptimizeConfig()
	config.m[name] = val
	cparams := config.toC(o.ctx)
	defer C.Z3_params_dec_ref(o.ctx.c, cparams)
	o.ctx.do(func() {
		C.Z3_optimize_set_params(o.ctx.c, o.o, cparams)
	})
	runtime.KeepAlive(o)
	runtime.KeepAlive(cparams)
}

// Assert adds val as a hard constraint to the set of predicates to be satisfied.
func (o *Optimize) Assert(val Bool) {
	o.ctx.do(func() {
		C.Z3_optimize_assert(o.ctx.c, o.o, val.c)
	})
	runtime.KeepAlive(o)
	runtime.KeepAlive(val)
}

// SoftOptions contains  options for soft constraints
type SoftOptions struct {
	ctx *Context
	weight, id string
}

// SoftOption is a SoftOptions setting closure
type SoftOption func(*SoftOptions)

// Weight sets the weight of a soft constraint.
func Weight(input interface{}) SoftOption {
	return func(args *SoftOptions) {
		switch w := input.(type) {
		case int:
			args.weight = strconv.FormatInt(int64(w), 10)
		case int8:
			args.weight = strconv.FormatInt(int64(w), 10)
		case int16:
			args.weight = strconv.FormatInt(int64(w), 10)
		case int32:
			args.weight = strconv.FormatInt(int64(w), 10)
		case int64:
			args.weight = strconv.FormatInt(w, 10)
		case uint:
			args.weight = strconv.FormatUint(uint64(w), 10)
		case uint8:
			args.weight = strconv.FormatUint(uint64(w), 10)
		case uint16:
			args.weight = strconv.FormatUint(uint64(w), 10)
		case uint32:
			args.weight = strconv.FormatUint(uint64(w), 10)
		case uint64:
			args.weight = strconv.FormatUint(w, 10)
		case float32:
			args.weight = strconv.FormatFloat(float64(w), 'f', -1, 64)
		case float64:
			args.weight = strconv.FormatFloat(w, 'f', -1, 64)
		case string:
			args.weight = w
		default:
			panic(fmt.Sprintf("unhandled type %T: weight should be a string or number", w))
		}
	}
}

// ID sets the ID of a soft constraint.
func ID(id string) SoftOption {
	return func(args *SoftOptions) {
		args.id = id
	}
}

// AssertSoft adds val as a soft constraint to the set of predicates to
// be satisfied.  Optional setters may be used to set a weight or id to
// the predicate.
// Weight represents the penalty for violating a constraint
// ID provides a mechanism to group soft constraints
func (o *Optimize) AssertSoft(val Bool, setters ...SoftOption) {
	args := &SoftOptions{
		ctx: o.ctx,
		weight: "1",
		id: "",
	}

	for _, setter := range setters {
		setter(args)
	}

	symbol := o.ctx.symbol(args.id)
	o.ctx.do(func() {
		C.Z3_optimize_assert_soft(o.ctx.c, o.o, val.c, C.CString(args.weight), symbol)
	})
	runtime.KeepAlive(o)
	runtime.KeepAlive(val)
	runtime.KeepAlive(args)
}

// Push saves the current state of the optimizer so that it can be restored
// with Pop.
func (o *Optimize) Push() {
	o.ctx.do(func() {
		C.Z3_optimize_push(o.ctx.c, o.o)
	})
	runtime.KeepAlive(o)
}

// Pop removes all predicates added since the proceeding Push.
func (o *Optimize) Pop() {
	o.ctx.do(func() {
		C.Z3_optimize_pop(o.ctx.c, o.o)
	})
	runtime.KeepAlive(o)
}

// Objective is an opaque handle that can be used to retrieve the
// upper/lower bounds of an optimization solution.
type Objective struct {
	*Optimize
	handle C.uint
}

// Lower produces the lower bound of the objective.
func (obj *Objective) Lower() Value {
	var ast AST
	obj.ctx.do(func() {
		cast := C.Z3_optimize_get_lower(obj.ctx.c, obj.o, obj.handle)
		ast = wrapAST(obj.ctx, cast)
	})
	return ast.AsValue()
}

// Upper produces the upper bound of the objective.
func (obj *Objective) Upper() Value {
	var ast AST
	obj.ctx.do(func() {
		cast := C.Z3_optimize_get_upper(obj.ctx.c, obj.o, obj.handle)
		ast = wrapAST(obj.ctx, cast)
	})
	return ast.AsValue()
}

// Maximize instructs the optimizer to solve for the maximum value
// of val.
func (o *Optimize) Maximize(val Value) *Objective {
	var handle C.uint
	o.ctx.do(func() {
		handle = C.Z3_optimize_maximize(o.ctx.c, o.o, val.AsAST().c)
	})
	runtime.KeepAlive(o)
	runtime.KeepAlive(val)
	return &Objective{o, handle}
}

// Minimize instructs the optimizer to solve for the minimum value
// of val.
func (o *Optimize) Minimize(val Value) *Objective {
	var handle C.uint
	o.ctx.do(func() {
		handle = C.Z3_optimize_minimize(o.ctx.c, o.o, val.AsAST().c)
	})
	runtime.KeepAlive(o)
	runtime.KeepAlive(val)
	return &Objective{o, handle}
}

// Check instructs the optimizer to check consistency and produce optimum
// values
func (o *Optimize) Check() (sat bool, err error) {
	var res C.Z3_lbool
	o.ctx.do(func() {
		res = C.Z3_optimize_check(o.ctx.c, o.o)
	})
	if res == C.Z3_L_UNDEF {
		// Get the reason
		o.ctx.do(func() {
			cerr := C.Z3_optimize_get_reason_unknown(o.ctx.c, o.o)
			err = &ErrSatUnknown{C.GoString(cerr)}
		})
	}
	runtime.KeepAlive(o)
	return res == C.Z3_L_TRUE, err
}

// Model returns the model for the last Check. Model panics if Check
// has not been called or the last Check did not return true.
func (o *Optimize) Model() *Model {
	var model *Model
	o.ctx.do(func() {
		model = wrapModel(o.ctx, C.Z3_optimize_get_model(o.ctx.c, o.o))
	})
	runtime.KeepAlive(o)
	return model
}

// String returns a string representation of o
func (o *Optimize) String() string {
	var res string
	o.ctx.do(func() {
		res = C.GoString(C.Z3_optimize_to_string(o.ctx.c, o.o))
	})
	runtime.KeepAlive(o)
	return res
}
