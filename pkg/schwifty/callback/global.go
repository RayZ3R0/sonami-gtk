package callback

import (
	"sync"
	"unsafe"

	"github.com/jwijenbergh/puregotk/v4/glib"
)

type MainThreadCallback func(u uintptr) bool

type MainThreadCallbackEntry struct {
	Callback MainThreadCallback
	Param    uintptr
}

var mainThreadCallbacks = make(map[uintptr]MainThreadCallbackEntry)
var mainThreadCallbacksLock = sync.RWMutex{}

var mainLoopHandler = glib.SourceFunc(func(ptr uintptr) bool {
	mainThreadCallbacksLock.RLock()
	entry, ok := mainThreadCallbacks[ptr]
	if !ok {
		mainThreadCallbacksLock.RUnlock()
		return glib.SOURCE_REMOVE
	}
	mainThreadCallbacksLock.RUnlock()

	shouldContinue := entry.Callback(entry.Param)
	if !shouldContinue {
		mainThreadCallbacksLock.Lock()
		delete(mainThreadCallbacks, ptr)
		mainThreadCallbacksLock.Unlock()
	}
	return shouldContinue
})

func OnMainThread(callback MainThreadCallback, params uintptr) uint {

	mainThreadCallbacksLock.RLock()
	if len(mainThreadCallbacks) >= 4096 {
		mainThreadCallbacksLock.RUnlock()
		return 0
	}
	mainThreadCallbacksLock.RUnlock()

	id := uintptr(unsafe.Pointer(&callback))
	mainThreadCallbacksLock.Lock()
	mainThreadCallbacks[id] = MainThreadCallbackEntry{
		Callback: callback,
		Param:    params,
	}
	mainThreadCallbacksLock.Unlock()

	return glib.IdleAdd(&mainLoopHandler, id)
}
