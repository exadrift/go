# ansi
provides ansi tools including key combination definition and detection, as well as text styling in a way that lends well to wrapping, splitting and rendering, regardless of byte length (unicode).


provides a mechanism to define keyboard key combinations and map to and from ANSI values, in order to provide configurable input processing for terminal based applications

## keystroke api
```
import (
    "github.com/exadrift/go/ansi/keys"
)

// get a KeyCombo object from a human readable key combination string
altEnter := keys.MustParseHumanName("alt+enter")

// now print the ANSI code
fmt.Print(altEnter.Ansi)

// lookup a KeyCombo object by an ANSI code
keyCombo, err := keys.ParseAnsiCode("\x1b")
if err != nil {
    // the ANSI sequence isn't mapped in this library
}
```

## text and styling

```
import (
    "github.com/exadrift/go/ansi/style"
)

// create a piece of styled text
text := T(Red.Fg(), "hello ", Red.Bg(), Black.Fg(), "world")

// apply styling defaults to the text
text = text.WithDefaultStyles(style.Blue.Bg(), style.Red.Fg())

// render ANSI string representation of the text as rows
rows := text.Render()
for _, row := range rows {
    fmt.Println(row)
}

// render can be used to constrain text to grid dimensions, both horizontally and vertically
// below ensures that each row is exactly 10 characters wide and there are a minimum of 10 rows
// rows are comprised of whitespace in order to pad to width and min row constraints
rows := text.Render(WithWidthConstraint(10), WithMinRows(10))
```