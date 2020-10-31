// SPDX-License-Identifier: Unlicense OR MIT

package headless

import (
	"gioui.org/app/internal/glimpl"
	"gioui.org/gpu/backend"
	"gioui.org/gpu/gl"

	_ "gioui.org/app/internal/cocoainit"
)

/*
#cgo CFLAGS: -DGL_SILENCE_DEPRECATION -Werror -Wno-deprecated-declarations -fmodules -fobjc-arc -x objective-c

#include <CoreFoundation/CoreFoundation.h>

__attribute__ ((visibility ("hidden"))) CFTypeRef gio_headless_newContext(void);
__attribute__ ((visibility ("hidden"))) void gio_headless_releaseContext(CFTypeRef ctxRef);
__attribute__ ((visibility ("hidden"))) void gio_headless_clearCurrentContext(CFTypeRef ctxRef);
__attribute__ ((visibility ("hidden"))) void gio_headless_makeCurrentContext(CFTypeRef ctxRef);
__attribute__ ((visibility ("hidden"))) void gio_headless_prepareContext(CFTypeRef ctxRef);
*/
import "C"

type nsContext struct {
	c        *glimpl.Functions
	ctx      C.CFTypeRef
	prepared bool
}

func newGLContext() (context, error) {
	ctx := C.gio_headless_newContext()
	return &nsContext{ctx: ctx, c: new(glimpl.Functions)}, nil
}

func (c *nsContext) MakeCurrent() error {
	C.gio_headless_makeCurrentContext(c.ctx)
	if !c.prepared {
		C.gio_headless_prepareContext(c.ctx)
		c.prepared = true
	}
	return nil
}

func (c *nsContext) ReleaseCurrent() {
	C.gio_headless_clearCurrentContext(c.ctx)
}

func (c *nsContext) Backend() (backend.Device, error) {
	return gl.NewBackend(c.c)
}

func (c *nsContext) Functions() *glimpl.Functions {
	return c.c
}

func (d *nsContext) Release() {
	if d.ctx != 0 {
		C.gio_headless_releaseContext(d.ctx)
		d.ctx = 0
	}
}
