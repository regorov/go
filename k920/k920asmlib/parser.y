//PREFIX yy
%{
    package k920asmlib
    
    import (
        "bufio"
        "fmt"
        "github.com/kierdavis/goutil"
        "io"
        "log"
        "strconv"
        "strings"
        "sync"
        
        . "github.com/kierdavis/go/k920"
    )
    
    var parserMutex sync.Mutex
    var parserOutput chan Object
    var parserBaseLabel string
%}

%union{
    s string
    i int64
    f float64
    R Register
    BR BitReg
    O Operand
    OL []Operand
}

%token EOF NL
%token <s> IDENTIFIER
%token <i> INTEGER
%token <f> FLOAT
%token <R> REGISTER
%token <BR> BITREG

%type <O> operand
%type <OL> operands opt_operands

%%

start   : objects EOF                   {return 0}

objects : objects object
        | object

object  : instruction NL
        | label NL

instruction : IDENTIFIER opt_operands
    {
        fullName := $1
        p := strings.Index(fullName, ".")
        
        var group, name string
        
        if p < 0 {
            group = "std"
            name = fullName
        } else {
            group = fullName[:p]
            name = fullName[p+1:]
        }
        
        parserOutput <- &Instruction{Group: group, Name: name, Operands: $2}
    }

label   : IDENTIFIER ':'
    {
        name := $1
        
        if name[0] == '.' {
            name = parserBaseLabel + name
        } else {
            parserBaseLabel = name
        }
        
        parserOutput <- &Label{Name: name}
    }

opt_operands    : operands              {$$ = $1}
                |                       {$$ = []Operand{}}

operands    : operands operand          {$$ = append($1, $2)}
            | operand                   {$$ = []Operand{$1}}

operand : REGISTER                      {$$ = $1}
        | BITREG                        {$$ = $1}
        | INTEGER                       {$$ = Integer($1)}
        | IDENTIFIER                    {$$ = LabelRef($1)}

%%

type lexer struct {
    input *bufio.Reader
    currentToken []rune
    lastColumn int
    lineno int
    column int
}

func isDigit(r rune) (ok bool) {
    return '0' <= r && r <= '9'
}

func isHexDigit(r rune) (ok bool) {
    return isDigit(r) || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F')
}

func isIdentStartChar(r rune) (ok bool) {
    return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || r == '_' || r == '.'
}

func isIdentChar(r rune) (ok bool) {
    return isIdentStartChar(r) || isDigit(r) || r == '+' || r == '-'
}

func isSpace(r rune) (ok bool) {
    return r == ' ' || r == '\t'
}

func isNewline(r rune) (ok bool) {
    return r == '\r' || r == '\n'
}

func newLexer(input io.Reader) (ll *lexer) {
    ll = &lexer{
        input: bufio.NewReader(input),
        lineno: 1,
        column: 1,
    }
    
    return ll
}

func (ll *lexer) Error(s string) {
    log.Printf("Syntax error: %s (at line %d col %d)", s, ll.lineno, ll.column)
}

func (ll *lexer) Lex(lval *yySymType) (tok int) {
    ll.AcceptRun(isSpace)
    ll.Discard()
    
    r := ll.Next()
    
    switch {
    case isIdentStartChar(r):
        ll.AcceptRun(isIdentChar)
        lval.s = ll.GetToken()
        return IDENTIFIER
    
    case isDigit(r), r == '+', r == '-':
        ll.Back()
        ll.AcceptOneOf('+', '-')
        
        digitFunc := isDigit
        if ll.Accept('0') && ll.AcceptOneOf('x', 'X') {
            digitFunc = isHexDigit
        }
        
        tok = INTEGER
        
        ll.AcceptRun(digitFunc)
        if ll.Accept('.') {
            ll.AcceptRun(digitFunc)
            tok = FLOAT
        }
        
        if ll.AcceptOneOf('e', 'E') {
            ll.AcceptOneOf('+', '-')
            ll.AcceptRun(isDigit)
        }
        
        if isIdentChar(ll.Peek()) {
            ll.Error("Malformed numeric literal")
            return ll.Lex(lval)
        }
        
        if tok == FLOAT {
            n, err := strconv.ParseFloat(ll.GetToken(), 32)
            if err != nil {
                ll.Error(err.Error())
                return ll.Lex(lval)
            }
            
            lval.f = n
        
        } else {
            n, err := strconv.ParseInt(ll.GetToken(), 0, 32)
            if err != nil {
                ll.Error(err.Error())
                return ll.Lex(lval)
            }
            
            lval.i = n
        }
        
        return tok
    
    case isNewline(r):
        ll.AcceptRun(isNewline)
        ll.Discard()
        return NL
    
    case r == '%':
        ll.Discard()
        ll.AcceptRun(isIdentChar)
        s := ll.GetToken()
        
        r, ok := Registers[s]
        if ok {
            lval.R = r
            return REGISTER
        }
        
        br, ok := BitRegs[s]
        if ok {
            lval.BR = br
            return BITREG
        }
        
        ll.Error("Invalid register: %" + s)
        return ll.Lex(lval)
    
    default:
        ll.Discard()
        return int(r)
    }
    
    ll.Error("Unreachable!")
    return ll.Lex(lval)
}

func (ll *lexer) Next() (r rune) {
    r, _, err := ll.input.ReadRune()
    if err != nil {
        if err == io.EOF {
            return EOF
        }
        
        ll.Error(err.Error())
    }
    
    ll.currentToken = append(ll.currentToken, r)
    ll.lastColumn = ll.column
    
    if r == '\n' {
        ll.lineno++
        ll.column = 1
    } else {
        ll.column++
    }
    
    return r
}

func (ll *lexer) Back() {
    err := ll.input.UnreadRune()
    if err == nil {
        if ll.currentToken[len(ll.currentToken)-1] == '\n' {
            ll.lineno--
            ll.column = ll.lastColumn
        
        } else {
            ll.column--
        }
        
        ll.currentToken = ll.currentToken[:len(ll.currentToken)-1]
    }
}

func (ll *lexer) Peek() (r rune) {
    r = ll.Next()
    ll.Back()
    return r
}

func (ll *lexer) Accept(r rune) (ok bool) {
    if ll.Next() == r {
        return true
    }
    
    ll.Back()
    return false
}

func (ll *lexer) AcceptOneOf(runes ...rune) (ok bool) {
    c := ll.Next()
    
    for _, r := range runes {
        if r == c {
            return true
        }
    }
    
    ll.Back()
    return false
}

func (ll *lexer) AcceptRun(f func(rune) bool) {
    for f(ll.Next()) {}
    
    ll.Back()
}

func (ll *lexer) GetToken() (s string) {
    s = util.EncodeUTF8(ll.currentToken)
    ll.Discard()
    return s
}

func (ll *lexer) Discard() {
    ll.currentToken = nil
}

func Parse(input io.Reader) (out Pipe) {
    out = make(Pipe)
    
    go func() {
        defer close(out)
        
        parserMutex.Lock()
        defer parserMutex.Unlock()
        
        parserOutput = out
        parserBaseLabel = ""
        
        yyParse(newLexer(input))
        
        parserOutput = nil
        parserBaseLabel = ""
    }()
    
    return out
}
