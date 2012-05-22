package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Coord struct {
	Filename string
	Lineno   int
}

func (c Coord) String() (str string) {
	return fmt.Sprintf("<%s:%d>", c.Filename, c.Lineno)
}

var coordStack = make([]Coord, 0, 100)

func pushCoord(c Coord) {
	coordStack = append(coordStack, c)
}

func getCoord() (c Coord) {
	return coordStack[len(coordStack)-1]
}

func getCoordRef() (c *Coord) {
	return &(coordStack[len(coordStack)-1])
}

func popCoord() {
	coordStack = coordStack[:len(coordStack)-1]
}

type yylexer struct {
	src     *bufio.Reader
	buf     []byte
	empty   bool
	current byte
}

func newLexer(src *bufio.Reader) (y *yylexer) {
	y = &yylexer{src: src}
	b, err := src.ReadByte()

	if err == nil {
		y.current = b
	}

	return y
}

func (y *yylexer) getc() (c byte) {
	if y.current != 0 {
		y.buf = append(y.buf, y.current)
	}

	y.current = 0
	b, err := y.src.ReadByte()
	if err == nil {
		y.current = b
	}

	return y.current
}

func (y *yylexer) Error(e string) {
	log.Fatal(e)
}

func (y *yylexer) Lex(lval *yySymType) int {
	//var err error

	c := y.current
	if y.empty {
		c = y.getc()
		y.empty = false
	}

yystate0:

	y.buf = y.buf[:0]
	lval.coord = getCoord()

	goto yystart1

	goto yystate1 // silence unused label error
yystate1:
	c = y.getc()
yystart1:
	switch {
	default:
		goto yyabort
	case c == '%':
		goto yystate4
	case c == '+' || c == '-':
		goto yystate15
	case c == '.' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate17
	case c == '\n' || c == '\r':
		goto yystate3
	case c == '\t' || c == ' ':
		goto yystate2
	case c >= '0' && c <= '9':
		goto yystate16
	}

yystate2:
	c = y.getc()
	switch {
	default:
		goto yyrule1
	case c == '\t' || c == ' ':
		goto yystate2
	}

yystate3:
	c = y.getc()
	switch {
	default:
		goto yyrule2
	case c == '\n' || c == '\r':
		goto yystate3
	}

yystate4:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == 'a':
		goto yystate5
	case c == 'q':
		goto yystate8
	case c == 's':
		goto yystate11
	case c == 'v':
		goto yystate13
	}

yystate5:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == 't':
		goto yystate7
	case c >= '0' && c <= '3':
		goto yystate6
	}

yystate6:
	c = y.getc()
	goto yyrule4

yystate7:
	c = y.getc()
	goto yyrule8

yystate8:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == '0':
		goto yystate9
	case c == '1':
		goto yystate10
	}

yystate9:
	c = y.getc()
	goto yyrule5

yystate10:
	c = y.getc()
	goto yyrule6

yystate11:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c == 'p':
		goto yystate12
	}

yystate12:
	c = y.getc()
	goto yyrule7

yystate13:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '7':
		goto yystate14
	}

yystate14:
	c = y.getc()
	goto yyrule3

yystate15:
	c = y.getc()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9':
		goto yystate16
	}

yystate16:
	c = y.getc()
	switch {
	default:
		goto yyrule9
	case c >= '0' && c <= '9':
		goto yystate16
	}

yystate17:
	c = y.getc()
	switch {
	default:
		goto yyrule10
	case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate17
	}

yyrule1: // [ \t]+

	goto yystate0
yyrule2: // [\r\n]+
	{

		getCoordRef().Lineno += len(strings.Replace(string(y.buf), "\r\n", "\n", -1))
		return NL
	}
yyrule3: // %v[0-7]
	{

		lval.r = V0 + Register(y.buf[2]-'0')
		return REGISTER
	}
yyrule4: // %a[0-3]
	{

		lval.r = A0 + Register(y.buf[2]-'0')
		return REGISTER
	}
yyrule5: // %q0
	{

		lval.r = Q0
		return REGISTER
	}
yyrule6: // %q1
	{

		lval.r = Q1
		return REGISTER
	}
yyrule7: // %sp
	{

		lval.r = SP
		return REGISTER
	}
yyrule8: // %at
	{

		lval.r = AT
		return REGISTER
	}
yyrule9: // [-+]?[0-9]+
	{

		i64, err := strconv.ParseInt(string(y.buf), 10, 0)
		if err != nil {
			log.Fatal(err)
		}

		lval.i = int(i64)
		return INTEGER
	}
yyrule10: // [a-zA-Z_.][a-zA-Z0-9_.]*
	{

		lval.s = string(y.buf)
		return IDENTIFIER
	}
	panic("unreachable")

	goto yyabort // silence unused label error

yyabort: // no lexem recognized
	y.empty = true
	return int(c)
}
