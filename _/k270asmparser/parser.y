%{
    package k270asmparser
    
    import (
        "fmt"
    )
%}

%union {
    node Node
    nodes []Node
    sval string
    ival int
}

%type <node> instruction
%type <node> integer
%type <node> integer_const
%type <node> label
%type <node> labelref
%type <node> operand
%type <node> register
%type <node> register_pair
%type <node> source
%type <node> signed_structdata
%type <node> stringdata
%type <node> structdata
%type <node> toplevel
%type <node> toplevel_nl
%type <node> unsigned_structdata

%type <nodes> operand_list
%type <nodes> toplevels

%token <sval> COMMA
%token <sval> COLON
%token <sval> IDENTIFIER
%token <sval> LITSTRING
%token <sval> NEWLINE
%token <sval> S8
%token <sval> S16
%token <sval> S32
%token <sval> S64
%token <sval> STRING
%token <sval> U8
%token <sval> U16
%token <sval> U32
%token <sval> U64

%token <ival> INTEGER
%token <ival> REGISTER

%%

source:
    toplevels
    {
        $$ = Node(NewAssembly($1))
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
|   structdata
|   stringdata
    {
        $$ = $1
    }

instruction:
    IDENTIFIER operand_list
    {
        $$ = Node(NewInstruction($1, $2))
    }

|   IDENTIFIER
    {
        $$ = Node(NewInstruction($1, nil))
    }

operand_list:
    operand_list COMMA operand
    {
        $$ = append($1, $3)
    }

|   operand
    {
        $$ = []Node{$1}
    }

operand:
    register
|   register_pair
|   integer
    {
        $$ = $1
    }

register:
    REGISTER
    {
        $$ = Node(NewRegister($1))
    }

register_pair:
    REGISTER COLON REGISTER
    {
        if !($1 % 2 == 0 && $3 - $1 == 1) {
            yyerror("Invalid register pair combination")
        }
        
        $$ = Node(NewRegisterPair($1))
    }

label:
    IDENTIFIER COLON
    {
        $$ = Node(NewLabel($1))
    }

structdata:
    unsigned_structdata
|   signed_structdata
    {
        $$ = $1
    }

unsigned_structdata:
    U8 integer
    {
        $$ = Node(NewU8Data($2))
    }
    
    U16 integer
    {
        $$ = Node(NewU16Data($2))
    }
    
    U32 integer
    {
        $$ = Node(NewU32Data($2))
    }
    
    U64 integer
    {
        $$ = Node(NewU64Data($2))
    }

signed_structdata:
    S8 integer
    {
        $$ = Node(NewS8Data($2))
    }
    
    S16 integer
    {
        $$ = Node(NewS16Data($2))
    }
    
    S32 integer
    {
        $$ = Node(NewS32Data($2))
    }
    
    S64 integer
    {
        $$ = Node(NewS64Data($2))
    }

stringdata:
    STRING LITSTRING
    {
        $$ = Node(NewStringData($2))
    }

integer:
    integer_const
|   labelref
    {
        $$ = $1
    }

integer_const:
    INTEGER
    {
        $$ = Node(NewInteger($1))
    }

labelref:
    IDENTIFIER
    {
        $$ = Node(NewLabelRef($1))
    }

%%
