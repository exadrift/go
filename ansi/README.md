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
What problem are we solving (use cases)

- represent a single line of text and 