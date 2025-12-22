package gui

func (w *WidgetImpl[T]) OnDestroy(callback func()) T {
	w.widget.ConnectDestroy(callback)
	return w.real
}

func (w *WidgetImpl[T]) OnRealize(callback func()) T {
	w.widget.ConnectRealize(callback)
	return w.real
}

func (w *WidgetImpl[T]) OnUnrealize(callback func()) T {
	w.widget.ConnectUnrealize(callback)
	return w.real
}
