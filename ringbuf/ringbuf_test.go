package ringbuf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemRetrieval(t *testing.T) {
	rb := New[string](10)
	rb.Push("1. hello")
	rb.Push("2. there")
	rb.Push("3. this")
	rb.Push("4. is")
	rb.Push("5. a")
	rb.Push("6. history")
	rb.Push("7. buffer")
	rb.Push("8. and")
	rb.Push("9. we")
	rb.Push("10. want")
	rb.Push("11. to")
	rb.Push("12. make")
	rb.Push("13. sure")
	rb.Push("14. it")
	rb.Push("15. works")

	assert.Equal(t, "15. works", rb.Item(0))
	assert.Equal(t, "14. it", rb.Item(-1))
	assert.Equal(t, "13. sure", rb.Item(-2))
	assert.Equal(t, "12. make", rb.Item(-3))
	assert.Equal(t, "11. to", rb.Item(-4))
	assert.Equal(t, "10. want", rb.Item(-5))
	assert.Equal(t, "9. we", rb.Item(-6))
	assert.Equal(t, 10, rb.Length())
}
