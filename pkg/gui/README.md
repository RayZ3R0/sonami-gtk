# codeberg.org/dergs/tidalwave/pkg/gui

This package is a DX-focused GUI library for developing desktop applications using GTK.

It is not required for a developer to know how the internals work, unless they want to
contribute to wrapping more widgets.

By default, GTK has to ways to be constructed: Either by using the .ui files, which are 
XML files that describe the layout of the UI, or by using the GTK API directly.
In the first case, the DX is restricted by the readability of the XML files, which tend to be 
rather verbose and hard to parse. In the second case, the DX is restricted by the fact that
objects are created, then set using imperative calls to methods, then assigned to their
parents.

This package is inspired by the way one builds UIs in Swift, using SwiftUI.
It provides a set of functions that 
can be used as parameters of one another, as well as a set of methods that can be chained 
to modify the properties of the widgets.

## Usage

When you want to create a widget tree, you start by declaring the root using its corresponding
function. In the parameters of the function, you can then pass any children that you want to 
affect to the root. Then, you can modify your widget by calling methods on the resulting 
object.

## Example
### Basic Usage
Let's imagine you want to put a GTK.Label inside a GTK.Box, and have it take the whole width.
You can achieve this by using the following code:

```go
import (
    "github.com/dergs/tidalwave/pkg/gui"
)

func MyTree() {
	return gui.Box(
		gui.Text("Hello World").
			HExpand(true)
	)
}
```

### Wrappers
Some GTK widgets may not yet be wrapped. Ideally, you can simply wrap them, by taking inspiration from
the existing wrapped widgets. However, if your widget is a single use, you can use the wrapper function 
to call base widget methods directly.

For example, let's imagine that the scales were not wrapped. You could use the wrapper function
on it:

```go
import (
	"github.com/dergs/tidalwave/pkg/gui"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func MyTree() {
	scale := gtk.NewScale(gtk.OrientationHorizontal, nil)
	scale.SetRange(0, 100)
	scale.SetValue(50)
	
	return gui.Box(
		gui.Wrapper(scale).
			HExpand(true)
	)
}
```
