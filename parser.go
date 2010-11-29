package λparse

import "./λeval"
import (
	"os"
	"io"
	"utf8"
	"unicode"
	"fmt"
)

const (
	_ = -iota
	end
	error
	vari
	lambda = 'λ'
	lpar   = '('
	rpar   = ')'
	dot    = '.'
)

const typestring map[int]string = map[int]string{
	end:   "end",
	error: "error",
	vari:  "variable",
}

const defLexBufLen = 256

type Token struct {
	Type    int
	Content *string
}

func (t Token) String() {
	var ts string
	var ok bool
	if ts, ok = typestring[t.Type]; !ok {
		ts = string(t.Type)
	}
	if t.Content == nil {
		return ts
	}
	return ts + "(" + *t.Content + ")"
}

type Lexer struct {
	rd      io.Reader
	rdEmpty bool
	buf     []byte
	pos     int
	current Token
}

func NewLexer(r io.Reader) *Lexer {
	var l = &Lexer{r, make([]byte, defLexBufLen)}
	if err := l.Next(); err != nil {
		//TODO: Figure out a way to handle this.
	}
	return l
}

func (l *Lexer) Current() Token {
	return l.current
}

func (l *Lexer) Next() os.Error {
	if !rdEmpty && l.pos >= len(l.buf)/2 {
		copy(l.buf, l.buf[l.pos:])
		if n, err := io.ReadFull(l.rd, l.buf[len(l.buf)-l.pos:]); n < l.pos {
			if !(err == io.ErrUnexpectedEOF || err == os.EOF) {
				l.current = Token{error}
			}
			l.rdEmpty = true
			l.buf = l.buf[:len(l.buf)-l.pos+n]
			return err
		}
		l.pos = 0
	}
	//XXX Continue here
}
