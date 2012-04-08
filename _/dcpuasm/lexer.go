// Mostly ripped from:
// https://github.com/cznic/golex/blob/master/calc/tokenizer.l

package main

import (
    "bufio"
    "fmt"
    "io"
    "strconv"
)

type yylexer struct {
    src *bufio.Reader
    buf []byte
    empty bool
    current byte
}

func die(err error) {
    if err != nil {
        panic(err)
    }
}

func newLexer(src *bufio.Reader) (y *yylexer) {
    y = new(yylexer)
    y.src = src
    
    b, err := src.ReadByte(); die(err)
    y.current = b
    
    return y
}

func (y *yylexer) getc() (b byte) {
    if y.current != 0 {
        y.buf = append(y.buf, y.current)
    }
    
    b, err := y.src.ReadByte()
    if err == io.EOF {
        y.current = 0xFF
    } else if err != nil {
        die(err)
    } else {
        y.current = b
    }
    
    return b
}

func (y *yylexer) Error(e string) {
    panic(e)
}

func (y *yylexer) Lex(lval *yySymType) (tok int) {
    c := y.current
    if y.empty {
        c = y.getc()
        y.empty = false
    }


yystate0:

    y.buf = y.buf[:0]

goto yystart1

goto yystate1 // silence unused label error
yystate1:
c = y.getc()
yystart1:
switch {
default:
goto yyabort
case c == '+':
goto yystate4
case c == ',':
goto yystate5
case c == '.':
goto yystate6
case c == '0':
goto yystate9
case c == ':':
goto yystate13
case c == '[':
goto yystate14
case c == '\n':
goto yystate3
case c == '\t' || c == ' ':
goto yystate2
case c == ']':
goto yystate15
case c == 'a':
goto yystate16
case c == 'b':
goto yystate17
case c == 'c':
goto yystate18
case c == 'i':
goto yystate19
case c == 'j':
goto yystate20
case c == 'o':
goto yystate21
case c == 'p':
goto yystate22
case c == 's':
goto yystate32
case c == 'x':
goto yystate34
case c == 'y':
goto yystate35
case c == 'z':
goto yystate36
case c == 'ÿ':
goto yystate37
case c >= '1' && c <= '9':
goto yystate10
case c >= 'A' && c <= 'Z' || c == '_' || c >= 'd' && c <= 'h' || c >= 'k' && c <= 'n' || c == 'q' || c == 'r' || c >= 't' && c <= 'w':
goto yystate8
}

yystate2:
c = y.getc()
switch {
default:
goto yyrule2
case c == '\t' || c == ' ':
goto yystate2
}

yystate3:
c = y.getc()
switch {
default:
goto yyrule26
case c == '\n':
goto yystate3
}

yystate4:
c = y.getc()
goto yyrule5

yystate5:
c = y.getc()
goto yyrule3

yystate6:
c = y.getc()
switch {
default:
goto yyrule22
case c == '+':
goto yystate7
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
goto yystate8
}

yystate7:
c = y.getc()
goto yyrule25

yystate8:
c = y.getc()
switch {
default:
goto yyrule22
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
goto yystate8
}

yystate9:
c = y.getc()
switch {
default:
goto yyrule23
case c == 'x':
goto yystate11
case c >= '0' && c <= '9':
goto yystate10
}

yystate10:
c = y.getc()
switch {
default:
goto yyrule23
case c >= '0' && c <= '9':
goto yystate10
}

yystate11:
c = y.getc()
switch {
default:
goto yyabort
case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
goto yystate12
}

yystate12:
c = y.getc()
switch {
default:
goto yyrule24
case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
goto yystate12
}

yystate13:
c = y.getc()
goto yyrule4

yystate14:
c = y.getc()
goto yyrule6

yystate15:
c = y.getc()
goto yyrule7

yystate16:
c = y.getc()
switch {
default:
goto yyrule8
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
goto yystate8
}

yystate17:
c = y.getc()
switch {
default:
goto yyrule9
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
goto yystate8
}

yystate18:
c = y.getc()
switch {
default:
goto yyrule10
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
goto yystate8
}

yystate19:
c = y.getc()
switch {
default:
goto yyrule14
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
goto yystate8
}

yystate20:
c = y.getc()
switch {
default:
goto yyrule15
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
goto yystate8
}

yystate21:
c = y.getc()
switch {
default:
goto yyrule18
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
goto yystate8
}

yystate22:
c = y.getc()
switch {
default:
goto yyrule22
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c == 'd' || c >= 'f' && c <= 'n' || c >= 'p' && c <= 't' || c >= 'v' && c <= 'z':
goto yystate8
case c == 'c':
goto yystate23
case c == 'e':
goto yystate24
case c == 'o':
goto yystate27
case c == 'u':
goto yystate29
}

yystate23:
c = y.getc()
switch {
default:
goto yyrule16
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
goto yystate8
}

yystate24:
c = y.getc()
switch {
default:
goto yyrule22
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
goto yystate8
case c == 'e':
goto yystate25
}

yystate25:
c = y.getc()
switch {
default:
goto yyrule22
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
goto yystate8
case c == 'k':
goto yystate26
}

yystate26:
c = y.getc()
switch {
default:
goto yyrule20
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
goto yystate8
}

yystate27:
c = y.getc()
switch {
default:
goto yyrule22
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
goto yystate8
case c == 'p':
goto yystate28
}

yystate28:
c = y.getc()
switch {
default:
goto yyrule21
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
goto yystate8
}

yystate29:
c = y.getc()
switch {
default:
goto yyrule22
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
goto yystate8
case c == 's':
goto yystate30
}

yystate30:
c = y.getc()
switch {
default:
goto yyrule22
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
goto yystate8
case c == 'h':
goto yystate31
}

yystate31:
c = y.getc()
switch {
default:
goto yyrule19
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
goto yystate8
}

yystate32:
c = y.getc()
switch {
default:
goto yyrule22
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
goto yystate8
case c == 'p':
goto yystate33
}

yystate33:
c = y.getc()
switch {
default:
goto yyrule17
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
goto yystate8
}

yystate34:
c = y.getc()
switch {
default:
goto yyrule11
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
goto yystate8
}

yystate35:
c = y.getc()
switch {
default:
goto yyrule12
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
goto yystate8
}

yystate36:
c = y.getc()
switch {
default:
goto yyrule13
case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
goto yystate8
}

yystate37:
c = y.getc()
goto yyrule1

yyrule1: // \xff
{

    return // EOF code (see getc() function)
}
yyrule2: // [ \t]+
{

    ;
goto yystate0
}
yyrule3: // ,
{

    return COMMA
}
yyrule4: // :
{

    return COLON
}
yyrule5: // \+
{

    return PLUS
}
yyrule6: // \[
{

    return LBRACKET
}
yyrule7: // \]
{

    return RBRACKET
}
yyrule8: // a
{

    lval.ival = 0
    return REGISTER
}
yyrule9: // b
{

    lval.ival = 1
    return REGISTER
}
yyrule10: // c
{

    lval.ival = 2
    return REGISTER
}
yyrule11: // x
{

    lval.ival = 3
    return REGISTER
}
yyrule12: // y
{

    lval.ival = 4
    return REGISTER
}
yyrule13: // z
{

    lval.ival = 5
    return REGISTER
}
yyrule14: // i
{

    lval.ival = 6
    return REGISTER
}
yyrule15: // j
{

    lval.ival = 7
    return REGISTER
}
yyrule16: // pc
{

    return PC
}
yyrule17: // sp
{

    return SP
}
yyrule18: // o
{

    return O
}
yyrule19: // push
{

    return PUSH
}
yyrule20: // peek
{

    return PEEK
}
yyrule21: // pop
{

    return POP
}
yyrule22: // [a-zA-Z_.][a-zA-Z0-9_.]*
{

    lval.sval = string(y.buf)
    return IDENTIFIER
}
yyrule23: // [0-9]+
{

    v, err := strconv.ParseInt(string(y.buf), 10, 32); die(err)
    lval.ival = int(v)
    return INTEGER
}
yyrule24: // 0x[0-9a-fA-F]+
{

    v, err := strconv.ParseInt(string(y.buf[2:]), 16, 32); die(err)
    lval.ival = int(v)
    return INTEGER
}
yyrule25: // ".+"
{

    lval.sval = string(y.buf[1:len(y.buf)-1])
    return LITSTRING
}
yyrule26: // \n+
{

    return NEWLINE
}
panic("unreachable")

goto yyabort // silence unused label error

yyabort: // no lexem recognized

    y.empty = true
    error(fmt.Sprintf("Invalid character: %q", c))
}
