package signals

var OnDisplayToast = NewSignal[func(string) bool]()
