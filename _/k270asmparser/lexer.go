package k270asmparser

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
    "regexp"
    "strconv"
    "strings"
)

var R_COLON = regexp.MustCompile("^:")
var R_COMMA = regexp.MustCompile("^,")
var R_IDENTIFIER = regexp.MustCompile("^[a-zA-Z_.][a-zA-Z0-9_]*")
var R_LITSTRING = regexp.MustCompile("^\".+\"|'.+'")
var R_NEWLINE = regexp.MustCompile("^\n+")
var R_DIRECTIVE = regexp.MustCompile("^@[a-zA-Z0-9]+")
var R_DECINT = regexp.MustCompile("^[0-9]+\b")
var R_HEXINT = regexp.MustCompile("^0x[0-9a-fA-F]+\b")

type Lexer struct {
    Reader *bufio.Reader
    Buffer []byte
    IsEof bool
}

func NewLexer(reader *bufio.Reader) (lexer *Lexer) {
    lexer = new(Lexer)
    lexer.Reader = reader
    lexer.Buffer = make([]byte, 0, 1024)
    lexer.IsEof = false
    return lexer
}

func (lexer *Lexer) Lex(lval *yySymType) (tok int) {
    if len(lexer.Buffer) < 256 && !lexer.IsEof {
        newbuf := make([]byte, 1024, 1024)
        copy(newbuf, lexer.Buffer)
        n, err := lexer.Reader.Read(newbuf[len(lexer.Buffer):])
        
        if err == io.EOF {
            lexer.IsEof = true
        } else if err != nil {
            panic(err)
        }
        
        lexer.Buffer = newbuf
    }
    
    lexer.Buffer = bytes.TrimLeft(lexer.Buffer, " \t")
    
    var match []byte
    
    match = R_COLON.Find(lexer.Buffer)
    if match != nil {lval.sval = string(match); return COLON}
    
    match = R_COMMA.Find(lexer.Buffer)
    if match != nil {lval.sval = string(match); return COLON}
    
    match = R_IDENTIFIER.Find(lexer.Buffer)
    if match != nil {
        l := strings.ToLower(string(match))
        
        if l == "z" {
            lval.ival = 0
        } else if l == "q" {
            lval.ival = 1
        } else if l == "k0" {
            lval.ival = 2
        } else if l == "k1" {
            lval.ival = 3
        } else if l == "a0" {
            lval.ival = 4
        } else if l == "a1" {
            lval.ival = 5
        } else if l == "a2" {
            lval.ival = 6
        } else if l == "a3" {
            lval.ival = 7
        } else if l == "v0" {
            lval.ival = 8
        } else if l == "v1" {
            lval.ival = 9
        } else if l == "v2" {
            lval.ival = 10
        } else if l == "v3" {
            lval.ival = 11
        } else if l == "v4" {
            lval.ival = 12
        } else if l == "v5" {
            lval.ival = 13
        } else if l == "v6" {
            lval.ival = 14
        } else if l == "v7" {
            lval.ival = 15
        } else {
            lval.sval = string(match)
            return IDENTIFIER
        }
        
        return REGISTER
    }
    
    match = R_LITSTRING.Find(lexer.Buffer)
    if match != nil {lval.sval = string(match[1:len(match) - 1]); return LITSTRING}
    
    match = R_NEWLINE.Find(lexer.Buffer)
    if match != nil {return NEWLINE}
    
    match = R_DIRECTIVE.Find(lexer.Buffer)
    if match != nil {
        lval.sval = string(match[1:])
        name := strings.ToLower(lval.sval)
        
        if name == "u8" {
            return U8
        } else if name == "u16" {
            return U16
        } else if name == "u32" {
            return U32
        } else if name == "u64" {
            return U64
        } else if name == "s8" {
            return S8
        } else if name == "s16" {
            return S16
        } else if name == "s32" {
            return S32
        } else if name == "s64" {
            return S64
        } else if name == "string" {
            return STRING
        } else {
            panic(fmt.Sprintf("Invalid directive: %s", name))
        }
    }
    
    match = R_DECINT.Find(lexer.Buffer)
    if match != nil {
        v, err := strconv.ParseInt(string(match), 10, 64)
        if err != nil {panic(err)}
        lval.ival = int(v)
        return INTEGER
    }
    
    match = R_HEXINT.Find(lexer.Buffer)
    if match != nil {
        v, err := strconv.ParseInt(string(match), 16, 64)
        if err != nil {panic(err)}
        lval.ival = int(v)
        return INTEGER
    }
    
    panic(fmt.Sprintf("Invalid text: %s...", lexer.Buffer[:10]))
    
    return -1
}

func (lexer *Lexer) Error(e string) {
    panic(e)
}
