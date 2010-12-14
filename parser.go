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
	End
	Var
	Begin
	Error  = utf8.RuneError
	Lambda = 'λ'
	LPar   = '('
	RPar   = ')'
	Dot    = '.'
)

const typestring map[int]string = map[int]string{
	End:   "End",
	Error: "Error",
	Var:   "Variable",
	Begin: "Beginning",
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
	return &Lexer{r, false, make([]byte, defLexBufLen), 0, Token{Begin, nil}}
}

func (l *Lexer) Current() Token {
	return l.current
}

func nonVarRune(rune int) bool {
	return !(unicode.IsLetter(rune) || unicode.IsDigit(rune))
}

func (l *Lexer) Next() os.Error {
	if !rdEmpty && l.pos >= len(l.buf)/2 {
		copy(l.buf, l.buf[l.pos:])
		if n, err := io.ReadFull(l.rd, l.buf[len(l.buf)-l.pos:]); n < l.pos {
			if !(err == io.ErrUnexpectedEOF || err == os.EOF) {
				l.current = Token{Error, &err.String()}
			}
			l.rdEmpty = true
			l.buf = l.buf[:len(l.buf)-l.pos+n]
			return err
		}
		l.pos = 0
	}
	if l.pos >= len(l.buf) {
		l.current = Token{End, nil}
		return os.EOF
	}
	var nextrune, nextwidth = utf8.DecodeRune(l.buf[l.pos:])
	if nextwidth < 1 {
		l.current = Token{Error, &"decoding error"}
		return l.current
	}
	switch nextrune {
	case Lambda, LPar, RPar, Dot:
		l.pos += nextwidth
		l.current = Token{nextrune, nil}
		return nil
	default:
		switch {
		case unicode.IsSpace(nextrune):
		case unicode.IsLetter(nextrune):
			var endVar = bytes.IndexFunc(l.buf[l.pos:], nonVarRune)
			for endVar < 0 {
				var newbuf = make([]byte, 2*len(l.buf))
				copy(newbuf, l.buf[l.pos:])
				if n, err := io.ReadFull(l.rd, l.buf[len(l.buf)-l.pos:]); n < l.pos {
					if !(err == io.ErrUnexpectedEOF || err == os.EOF) {
						l.current = Token{Error, &err.String()}
					}
					l.rdEmpty = true
					l.buf = l.buf[:len(l.buf)-l.pos+n]
					return err
				}
				l.pos = 0
				l.buf = newbuf
				endVar = bytes.IndexFunc(l.buf[l.pos:], nonVarRune)
			}
			l.pos = endVar
			l.current = Token{Var, &string(l.buf[l.pos:endVar])}
			return nil
		default:
			l.pos += nextwidth
			l.current = Token{Error, &string(nextrune)}
			return l.current
		}
	}
}
