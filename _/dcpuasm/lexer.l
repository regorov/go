%{
// Mostly ripped from:
// https://github.com/cznic/golex/blob/master/calc/tokenizer.l

// Compile this with cznic's GoLex, not anneli's!

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
        y.current = 0
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
%}

%yyc c
%yyn c = y.getc()

%%
    y.buf = y.buf[:0]

\0
    return

[ \t]+
    ;

,
    return COMMA

:
    return COLON

\+
    return PLUS

\[
    return LBRACKET

\]
    return RBRACKET

a
    lval.ival = 0
    return REGISTER

b
    lval.ival = 1
    return REGISTER

c
    lval.ival = 2
    return REGISTER

x
    lval.ival = 3
    return REGISTER

y
    lval.ival = 4
    return REGISTER

z
    lval.ival = 5
    return REGISTER

i
    lval.ival = 6
    return REGISTER

j
    lval.ival = 7
    return REGISTER

pc
    return PC

sp
    return SP

o
    return O

push
    return PUSH

peek
    return PEEK

pop
    return POP

[a-zA-Z_.][a-zA-Z0-9_.]*
    lval.sval = string(y.buf)
    return IDENTIFIER

[0-9]+
    v, err := strconv.ParseInt(string(y.buf), 10, 32); die(err)
    lval.ival = int(v)
    return INTEGER

0x[0-9a-fA-F]+
    v, err := strconv.ParseInt(string(y.buf[2:]), 16, 32); die(err)
    lval.ival = int(v)
    return INTEGER

".+"
    lval.sval = string(y.buf[1:len(y.buf)-1])
    return LITSTRING

\n+
    return NEWLINE

%%

    y.empty = true
    panic(fmt.Sprintf("Invalid character: %q", c))
}
