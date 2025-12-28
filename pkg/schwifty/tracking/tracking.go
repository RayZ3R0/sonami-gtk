package tracking

import "sync"

type TrackedWidget struct {
	IsReferencedByGo bool
	Type             string
}

var trackedWidgets map[uintptr]*TrackedWidget = make(map[uintptr]*TrackedWidget)
var lock sync.Mutex = sync.Mutex{}

func Track(ptr uintptr, widgetType string) {
	lock.Lock()
	defer lock.Unlock()
	trackedWidgets[ptr] = &TrackedWidget{
		IsReferencedByGo: true,
		Type:             widgetType,
	}
}

func TrackGC(ptr uintptr) {
	lock.Lock()
	defer lock.Unlock()
	if val, ok := trackedWidgets[ptr]; ok {
		val.IsReferencedByGo = false
	}
}

func Untrack(ptr uintptr) {
	lock.Lock()
	defer lock.Unlock()
	delete(trackedWidgets, ptr)
}

func Alive() []*TrackedWidget {
	lock.Lock()
	defer lock.Unlock()
	var alive []*TrackedWidget
	for _, widget := range trackedWidgets {
		alive = append(alive, widget)
	}
	return alive
}
