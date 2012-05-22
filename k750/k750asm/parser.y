%{
    package main
    
    import (
        "fmt"
    )
%}

%union {
    i int
    r Register
    s string
    it Item
    o Operand
    oL []Operand
    l Literal
    
    coord Coord
}

%token <i> INTEGER, NL
%token <r> REGISTER
%token <s> IDENTIFIER

%type <it> rawitem
%type <o> operand
%type <oL> operands, opt_operands
%type <l> integer

%%

assembly:           itemlist                    {close(parserOutput)}

itemlist:           itemlist item
                |   item

item:               rawitem NL                  {parserOutput <- $1}

rawitem:            IDENTIFIER opt_operands     {$$ = Item(&Instruction{coord: yyS[yypt-1].coord, name: $1, operands: $2})}
                |   IDENTIFIER ':'              {$$ = Item(&Label{coord: yyS[yypt-1].coord, name: $1})}

opt_operands:       operands                    {$$ = $1}
                |                               {$$ = nil}

operands:           operands ',' operand        {$$ = append($1, $3)}
                |   operand                     {$$ = []Operand{$1}}

operand:            integer                     {$$ = Operand(&LiteralOperand{Literal: $1})}
                |   REGISTER                    {$$ = Operand(&RegisterOperand{num: $1})}

integer:            INTEGER                     {$$ = Literal(&ConstantLiteral{value: uint32($1)})}
                |   IDENTIFIER                  {$$ = Literal(&LabelLiteral{coord: yyS[yypt-1].coord, name: $1})}

%%
