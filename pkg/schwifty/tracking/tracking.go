package tracking

type TrackedWidget struct {
	IsReferencedByGo bool
	Type             string
}

var trackedWidgets map[uintptr]*TrackedWidget = make(map[uintptr]*TrackedWidget)

func Track(ptr uintptr, widgetType string) {
	trackedWidgets[ptr] = &TrackedWidget{
		IsReferencedByGo: true,
		Type:             widgetType,
	}
}

func TrackGC(ptr uintptr) {
	if val, ok := trackedWidgets[ptr]; ok {
		val.IsReferencedByGo = false
	}
}

func Untrack(ptr uintptr) {
	delete(trackedWidgets, ptr)
}

func Alive() []*TrackedWidget {
	var alive []*TrackedWidget
	for _, widget := range trackedWidgets {
		alive = append(alive, widget)
	}
	return alive
}
