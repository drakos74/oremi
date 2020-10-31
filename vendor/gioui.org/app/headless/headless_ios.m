// SPDX-License-Identifier: Unlicense OR MIT

// +build darwin,ios

@import OpenGLES;

#include <CoreFoundation/CoreFoundation.h>
#include "_cgo_export.h"

void gio_headless_releaseContext(CFTypeRef ctxRef) {
	CFBridgingRelease(ctxRef);
}

CFTypeRef gio_headless_newContext(void) {
	EAGLContext *ctx = [[EAGLContext alloc] initWithAPI:kEAGLRenderingAPIOpenGLES3];
	if (ctx == nil) {
		return nil;
	}
	return CFBridgingRetain(ctx);
}

void gio_headless_clearCurrentContext(CFTypeRef ctxRef) {
	[EAGLContext setCurrentContext:nil];
}

void gio_headless_makeCurrentContext(CFTypeRef ctxRef) {
	EAGLContext *ctx = (__bridge EAGLContext *)ctxRef;
	[EAGLContext setCurrentContext:ctx];
}

void gio_headless_prepareContext(CFTypeRef ctxRef) {
}
