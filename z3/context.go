// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import (
	"runtime"
	"sync"
	"unsafe"
)

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"

// Context is an environment for creating symbolic values and checking
// satisfiability.
//
// Nearly all interaction with Z3 is done relative to a Context.
// Values are bound to the Context that created them and cannot be
// combined with Values from other Contexts.
//
// Context is thread-safe. However, most operations block other
// operations (one notable exception is Interrupt). Hence, to do
// things in parallel, it's best to create multiple Contexts.
type Context struct {
	c C.Z3_context

	syms map[string]C.Z3_symbol

	// lock protects AST reference counts and the context's last
	// error. Use Context.do to acquire this around a Z3 operation
	// and panic if the operation has an error status.
	lock sync.Mutex
}

// NewContext returns a new Z3 context with the given configuration.
//
// If cfg is nil, the default configuration is used.
func NewContext(cfg *Config) *Context {
	if cfg == nil {
		cfg = NewConfig()
	}
	ctx := &Context{
		C.Z3_mk_context_rc(cfg.c),
		make(map[string]C.Z3_symbol),
		sync.Mutex{},
	}
	runtime.SetFinalizer(ctx, func(ctx *Context) {
		C.Z3_del_context(ctx.c)
	})
	// Disable the default error handler, which exits the program.
	//
	// TODO: It might simplify our own error handling to install a
	// handler that does a Go panic. It looks like Z3's error
	// handler logic might actually be designed to handle things
	// like longjmp.
	C.Z3_set_error_handler(ctx.c, nil)
	return ctx
}

// TODO: Z3_update_param_value

// Interrupt stops the current solver, simplifier, or tactic being
// executed by ctx.
func (ctx *Context) Interrupt() {
	C.Z3_interrupt(ctx.c)
	runtime.KeepAlive(ctx)
}

// do calls f with a per-context lock held and panics if the operation
// produces an error code.
//
// Unfortunately, we can't just say that Contexts are not thread-safe
// because we can't help but run finalizers asynchronously, which
// means we need to synchronize both reference counts and the
// per-context last error state.
func (ctx *Context) do(f func()) {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()
	f()
	ecode := C.Z3_get_error_code(ctx.c)
	if ecode == C.Z3_OK {
		runtime.KeepAlive(ctx)
		return
	}
	msg := C.Z3_get_error_msg_ex(ctx.c, ecode)
	runtime.KeepAlive(ctx)
	// TODO: Lift the Z3 errors to better Go errors. At least wrap
	// the string in a type and consider using the error code to
	// determine which of different error types to use.
	panic(C.GoString(msg))
}

// symbol interns name as a Z3 symbol.
func (ctx *Context) symbol(name string) C.Z3_symbol {
	if sym, ok := ctx.syms[name]; ok {
		return sym
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var sym C.Z3_symbol
	ctx.do(func() {
		sym = C.Z3_mk_string_symbol(ctx.c, cname)
		ctx.syms[name] = sym
	})
	return sym
}
