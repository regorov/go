%{
    package main
    
    import (
        "bufio"
        "fmt"
        "os"
    )
    
    var parseResult Node
%}

%union {
    ival int
    sval string
    node Node
    nodes []Node
}

%token <sval> COMMA
%token <sval> COLON
%token <sval> IDENTIFIER
%token <sval> LBRACKET
%token <sval> LITSTRING
%token <sval> NEWLINE
%token <sval> O
%token <sval> PC
%token <sval> PEEK
%token <sval> PLUS
%token <sval> POP
%token <sval> PUSH
%token <sval> RBRACKET
%token <sval> SP

%token <ival> REGISTER
%token <ival> INTEGER

%type <node> instruction
%type <node> integer
%type <node> label
%type <node> operand
%type <node> register
%type <node> source
%type <node> toplevel
%type <node> toplevel_nl

%type <nodes> operands
%type <nodes> toplevels

%%

source:
    toplevels
    {
        parseResult = NewAssembly($1)
    }

toplevels:
    toplevels toplevel_nl
    {
        $$ = append($1, $2)
    }
    
|   toplevel_nl
    {
        $$ = []Node{$1}
    }

toplevel_nl:
    toplevel NEWLINE
    {
        $$ = $1
    }

toplevel:
    instruction
|   label
    {
        $$ = $1
    }

instruction:
    IDENTIFIER operands
    {
        $$ = NewInstruction($1, $2)
    }

|   IDENTIFIER
    {
        $$ = NewInstruction($1, []Node{})
    }

operands:
    operands COMMA operand
    {
        $$ = append($1, $3)
    }

|   operand
    {
        $$ = []Node{$1}
    }

operand:
    register
    {
        $$ = NewOperand(O_REG, $1, nil)
    }

|   LBRACKET register RBRACKET
    {
        $$ = NewOperand(O_MEM, $2, nil)
    }

|   LBRACKET register PLUS integer RBRACKET
    {
        $$ = NewOperand(O_MEMDISP, $2, $4)
    }

|   LBRACKET integer PLUS register RBRACKET
    {
        $$ = NewOperand(O_MEMDISP, $4, $2)
    }

|   POP
    {
        $$ = NewOperand(O_POP, nil, nil)
    }

|   PEEK
    {
        $$ = NewOperand(O_PEEK, nil, nil)
    }

|   PUSH
    {
        $$ = NewOperand(O_PUSH, nil, nil)
    }

|   SP
    {
        $$ = NewOperand(O_SP, nil, nil)
    }

|   PC
    {
        $$ = NewOperand(O_PC, nil, nil)
    }

|   O
    {
        $$ = NewOperand(O_O, nil, nil)
    }

|   LBRACKET integer RBRACKET
    {
        $$ = NewOperand(O_MEMIMM, $2, nil)
    }

|   integer
    {
        $$ = NewOperand(O_IMM, $1, nil)
    }

register:
    REGISTER
    {
        $$ = NewRegister($1)
    }

label:
    IDENTIFIER COLON
    {
        $$ = NewLabel($1)
    }

integer:
    INTEGER
    {
        $$ = NewConstant($1)
    }

|   IDENTIFIER
    {
        $$ = NewLabelRef($1)
    }

%%

func parse(reader io.Reader) (asm Node) {
    yyParse(newLexer(bufio.NewReader(reader)))
    
    exitIfErrors()
    
    return parseResult
}
