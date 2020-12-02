package golas

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"
)

// HandlerFunc func used as lexer state
type HandlerFunc func(*Lexer) HandlerFunc

// Lexer can perform lexical analysis on .las files
type Lexer struct {
	buffer   *bytes.Buffer
	char     rune
	count    map[string]int
	dots     int
	handler  HandlerFunc
	line     int
	position int
	reader   *bufio.Reader
	tokens   chan Token
}

// NewLexer creates a new Lexer
func NewLexer(r io.Reader) Lexer {
	return Lexer{
		reader: bufio.NewReader(r),
		tokens: make(chan Token, 3),
		buffer: &bytes.Buffer{},
		count:  map[string]int{"a": 0, "v": 0, "w": 0, "c": 0},
	}
}

// NextToken reads the next token from our tokens chan
func (l *Lexer) NextToken() Token {
	for {
		select {
		case token := <-l.tokens:
			return token
		default:
			if l.handler == nil {
				// return Token{Type: TEndOfFile}
				panic("lexer not started : lexer.handler is nil : did you forget to call `lexer.Start(...)`?")
			}
			l.handler = l.handler(l)
		}
	}
}

// Start sets our handler, in turn starting lexical analysis
func (l *Lexer) Start(hf HandlerFunc) {
	l.handler = hf
}

// emit places a token on our tokens chan
func (l *Lexer) emit(t TokenType) {
	l.tokens <- Token{t, l.buffer.String()}
	l.buffer.Reset()
}

// overwriteBuffer clears our buffer then writes a string to it
func (l *Lexer) overwriteBuffer(s string) {
	l.buffer.Reset()
	l.buffer.WriteString(s)
}

func (l *Lexer) peekNext() rune {
	r, e := l.reader.Peek(1)
	if e != nil {
		return -1
	}
	rr, _ := utf8.DecodeRune(r)
	return rr
}

func (l *Lexer) read() rune {
	r, _, e := l.reader.ReadRune()
	if e != nil {
		r = -1 // Use -1 to signal EOF
	}
	return r
}

// step consumes the next rune from our reader
func (l *Lexer) step() {
	char := l.read()

	l.position++
	l.buffer.WriteRune(char)
	l.char = char

	switch char {
	case '\n':
		l.line++
		l.position = 0
		l.dots = 0
	case '.':
		l.dots = l.dots + 1
	}
}

// stepUntil reads from current line position until we read one of the specified runes
func (l *Lexer) stepUntil(oneOfChars ...rune) {
Loop:
	for {
		l.step()
		for i := 0; i < len(oneOfChars); i++ {
			if l.char == oneOfChars[i] {
				break Loop
			}
		}
	}
}

// truncate our buffer by 1. If our buffer were a string, this removes the last character
func (l *Lexer) truncate() {
	l.buffer.Truncate(l.buffer.Len() - 1)
}

func (l *Lexer) validateSection(s string) {
	if l.count[s] >= 1 {
		panic(fmt.Errorf("invalid las file : expected 1 got %d of section %s : line %d : position %d", l.count[s]+1, s, l.line+1, l.position))
	}
	l.count[s]++
}
