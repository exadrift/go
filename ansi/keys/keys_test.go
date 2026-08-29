package keys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHumanNameParser(t *testing.T) {
	combo, err := ParseHumanName("ctrl+alt+pgup")
	assert.NoError(t, err)
	assert.Equal(t, "ctrl+alt+pgup", combo.HumanName)
}

func TestHumanNameParserReorder(t *testing.T) {
	combo, err := ParseHumanName("pgup+alt+ctrl")
	assert.NoError(t, err)
	assert.Equal(t, "ctrl+alt+pgup", combo.HumanName)
}

func TestHumanNameParserSpaces(t *testing.T) {
	combo, err := ParseHumanName("pgup  + alt+ctrl")
	assert.NoError(t, err)
	assert.Equal(t, "ctrl+alt+pgup", combo.HumanName)
}
