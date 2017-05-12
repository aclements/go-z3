// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import (
	"runtime"
	"strconv"
	"unsafe"
)

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"

type Config struct {
	c C.Z3_config
	noEq
}

// TODO: Z3 has a lot of configuration things. It has a global
// configuration, a pre-context configuration, context configuration,
// and Z3_params. Maybe we unify them? Maybe these should just be
// maps?

// NewConfig returns a new configuration for creating contexts.
//
// The params argument provides initial configuration settings. It
// must alternate between configuration key and value.
//
// See Z3_mk_config documentation for accepted configuration settings.
func NewConfig(params ...string) *Config {
	if len(params)%2 != 0 {
		panic("len(params) must be even")
	}
	c := &Config{C.Z3_mk_config(), noEq{}}
	runtime.SetFinalizer(c, func(c *Config) {
		C.Z3_del_config(c.c)
	})
	for i := 0; i < len(params); i += 2 {
		c.Set(params[i], params[i+1])
	}
	return c
}

func (c *Config) Set(name, value string) {
	cname, cvalue := C.CString(name), C.CString(value)
	defer C.free(unsafe.Pointer(cname))
	defer C.free(unsafe.Pointer(cvalue))
	C.Z3_set_param_value(c.c, cname, cvalue)
	runtime.KeepAlive(c)
}

func (c *Config) SetBool(name string, value bool) {
	if value {
		c.Set(name, "true")
	} else {
		c.Set(name, "false")
	}
}

func (c *Config) SetInt(name string, value int) {
	c.Set(name, strconv.Itoa(value))
}
