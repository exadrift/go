package tui

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstrain(t *testing.T) {
	assert.Equal(t, "hel..", Constrain("hello world", 5))
}

func TestTokenSplitter(t *testing.T) {
	text := fmt.Sprintf("%shello world%shello", UpArrow, DownArrow)
	tokens := SplitAtAnsiTokens(text)
	assert.Len(t, tokens, 4)
	assert.Equal(t, TokenTypeAnsiCode, tokens[0].TokenType)
	assert.Equal(t, UpArrow, tokens[0].Text)
	assert.Equal(t, TokenTypeText, tokens[1].TokenType)
	assert.Equal(t, "hello world", tokens[1].Text)
	assert.Equal(t, TokenTypeAnsiCode, tokens[2].TokenType)
	assert.Equal(t, DownArrow, tokens[2].Text)
	assert.Equal(t, TokenTypeText, tokens[3].TokenType)
	assert.Equal(t, "hello", tokens[3].Text)
}

func TestTokenSplitter2(t *testing.T) {
	text := fmt.Sprintf("hello world%shello", DownArrow)
	tokens := SplitAtAnsiTokens(text)
	assert.Len(t, tokens, 3)
	assert.Equal(t, TokenTypeText, tokens[0].TokenType)
	assert.Equal(t, "hello world", tokens[0].Text)
	assert.Equal(t, TokenTypeAnsiCode, tokens[1].TokenType)
	assert.Equal(t, DownArrow, tokens[1].Text)
	assert.Equal(t, TokenTypeText, tokens[2].TokenType)
	assert.Equal(t, "hello", tokens[2].Text)
}

func TestTokenSplitter3(t *testing.T) {
	text := fmt.Sprintf("hello world%s", DownArrow)
	tokens := SplitAtAnsiTokens(text)
	assert.Len(t, tokens, 2)
	assert.Equal(t, TokenTypeText, tokens[0].TokenType)
	assert.Equal(t, "hello world", tokens[0].Text)
	assert.Equal(t, TokenTypeAnsiCode, tokens[1].TokenType)
	assert.Equal(t, DownArrow, tokens[1].Text)
}

func TestTokenSplitter4(t *testing.T) {
	text := "hello world"
	tokens := SplitAtAnsiTokens(text)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TokenTypeText, tokens[0].TokenType)
	assert.Equal(t, "hello world", tokens[0].Text)
}

func TestWrapTextBasic(t *testing.T) {
	text := "this is a short piece of text"
	lines := WrapTextBasic(text, 5)
	assert.Equal(t, 6, len(lines))
	assert.Equal(t, "this ", lines[0])
	assert.Equal(t, "is a ", lines[1])
	assert.Equal(t, "short", lines[2])
	assert.Equal(t, " piec", lines[3])
	assert.Equal(t, "e of ", lines[4])
	assert.Equal(t, "text", lines[5])
}

func TestWrapTextBasicWithNewline(t *testing.T) {
	text := "this is a short piece of\na\nthing"
	lines := WrapTextBasic(text, 5)
	assert.Equal(t, 7, len(lines))
	assert.Equal(t, "this ", lines[0])
	assert.Equal(t, "is a ", lines[1])
	assert.Equal(t, "short", lines[2])
	assert.Equal(t, " piec", lines[3])
	assert.Equal(t, "e of", lines[4])
	assert.Equal(t, "a", lines[5])
	assert.Equal(t, "thing", lines[6])
}
