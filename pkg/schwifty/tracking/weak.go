package tracking

import (
	"runtime"
	"unsafe"

	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gobject"
)

//go:linkname xWeakRefInit github.com/jwijenbergh/puregotk/v4/gobject.xWeakRefInit
var xWeakRefInit func(uintptr, uintptr)

//go:linkname xWeakRefClear github.com/jwijenbergh/puregotk/v4/gobject.xWeakRefClear
var xWeakRefClear func(uintptr)

//go:linkname xWeakRefGet github.com/jwijenbergh/puregotk/v4/gobject.xWeakRefGet
var xWeakRefGet func(uintptr) uintptr

type WeakRef struct {
	ptr uintptr
}

func (x *WeakRef) Clear() {
	if x.ptr == 0 {
		return
	}

	xWeakRefClear(x.ptr)
	glib.Free(x.ptr)
	x.ptr = 0
}

func (x *WeakRef) Get() *gobject.Object {
	if x.ptr == 0 {
		return nil
	}

	var cls *gobject.Object

	cret := xWeakRefGet(x.ptr)

	if cret == 0 {
		return nil
	}
	cls = &gobject.Object{}
	cls.Ptr = cret
	return cls
}

func (x *WeakRef) Use(cb func(obj *gobject.Object)) bool {
	if x.ptr == 0 {
		return false
	}

	obj := x.Get()
	if obj != nil {
		defer obj.Unref()
		cb(obj)
		return true
	}
	return false
}

func NewWeakRef[T Trackable](obj T) *WeakRef {
	x := &WeakRef{
		ptr: glib.Malloc(uint(unsafe.Sizeof(uintptr(0)) * 2)),
	}
	xWeakRefInit(x.ptr, obj.GoPointer())
	runtime.SetFinalizer(x, func(x *WeakRef) {
		x.Clear()
	})
	return x
}
