%{
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
        Lineno int
    }
    
    func (c Coord) String() (str string) {
        return fmt.Sprintf("<%s:%d>", c.Filename, c.Lineno)
    }
    
    var coordStack = make([]Coord, 0, 100)
    
    func pushCoord(c Coord) {
        coordStack = append(coordStack, c)
    }
    
    func getCoord() (c Coord) {
        return coordStack[len(coordStack) - 1]
    }
    
    func getCoordRef() (c *Coord) {
        return &(coordStack[len(coordStack) - 1])
    }
    
    func popCoord() {
        coordStack = coordStack[:len(coordStack) - 1]
    }
    
    type yylexer struct {
        src *bufio.Reader
        buf []byte
        empty bool
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
%}

%yyc c
%yyn c = y.getc()

%%
    y.buf = y.buf[:0]
    lval.coord = getCoord()

[ \t]+

[\r\n]+
    getCoordRef().Lineno += len(strings.Replace(string(y.buf), "\r\n", "\n", -1))
    return NL

%v[0-7]
    lval.r = V0 + Register(y.buf[2] - '0')
    return REGISTER

%a[0-3]
    lval.r = A0 + Register(y.buf[2] - '0')
    return REGISTER

%q0
    lval.r = Q0
    return REGISTER

%q1
    lval.r = Q1
    return REGISTER

%sp
    lval.r = SP
    return REGISTER

%at
    lval.r = AT
    return REGISTER

[-+]?[0-9]+
    i64, err := strconv.ParseInt(string(y.buf), 10, 0)
    if err != nil {
        log.Fatal(err)
    }
    
    lval.i = int(i64)
    return INTEGER

[a-zA-Z_.][a-zA-Z0-9_.]*
    lval.s = string(y.buf)
    return IDENTIFIER

%%
    y.empty = true
    return int(c)
}
