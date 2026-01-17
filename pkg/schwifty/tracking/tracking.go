package tracking

import "sync"

type TrackedObject struct {
	IsReferencedByGo bool
	Type             string
}

var trackedObjects map[uintptr]*TrackedObject = make(map[uintptr]*TrackedObject)
var lock sync.Mutex = sync.Mutex{}

func Track(ptr uintptr, objectType string) {
	lock.Lock()
	defer lock.Unlock()
	trackedObjects[ptr] = &TrackedObject{
		IsReferencedByGo: true,
		Type:             objectType,
	}
}

func TrackGC(ptr uintptr) {
	lock.Lock()
	defer lock.Unlock()
	if val, ok := trackedObjects[ptr]; ok {
		val.IsReferencedByGo = false
	}
}

func Untrack(ptr uintptr) {
	lock.Lock()
	defer lock.Unlock()
	delete(trackedObjects, ptr)
}

func Alive() []*TrackedObject {
	lock.Lock()
	defer lock.Unlock()
	var alive []*TrackedObject
	for _, object := range trackedObjects {
		alive = append(alive, object)
	}
	return alive
}
