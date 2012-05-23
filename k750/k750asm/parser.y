%{
    package main
    
    import (
        "fmt"
        "log"
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
%type <o> memory_operand, memory_operand_content, operand
%type <oL> operands, opt_operands
%type <l> integer

%%

assembly:               itemlist                                {close(parserOutput)}

itemlist:               itemlist item
                    |   item

item:                   rawitem NL                              {parserOutput <- $1}

rawitem:                IDENTIFIER opt_operands                 {$$ = Item(&Instruction {coord: yyS[yypt-1].coord, name: $1, operands: $2})}
                    |   IDENTIFIER ':'                          {$$ = Item(&Label       {coord: yyS[yypt-1].coord, name: $1})}

opt_operands:           operands                                {$$ = $1}
                    |                                           {$$ = nil}

operands:               operands ',' operand                    {$$ = append($1, $3)}
                    |   operand                                 {$$ = []Operand{$1}}

operand:                integer                                 {$$ = Operand(&LiteralOperand  {coord: yyS[yypt-1].coord, Literal: $1})}
                    |   REGISTER                                {$$ = Operand(&RegisterOperand {coord: yyS[yypt-1].coord, num: $1})}
                    |   memory_operand                          {$$ = $1}

memory_operand:         INTEGER '[' memory_operand_content ']'
    {
        size := $1
        if size != 8 && size != 16 && size != 32 {
            log.Fatalf("Invalid memory addressing size: %d (expected 8, 16 or 32)", size)
        }
        
        $$ = $3
        $$.SetSize(MemSize(size))
    }

memory_operand_content: REGISTER                                {$$ = Operand(&MemoryOperand {coord: yyS[yypt-1].coord, reg: $1, disp: Zero})}
                    |   REGISTER '+' integer                    {$$ = Operand(&MemoryOperand {coord: yyS[yypt-1].coord, reg: $1, disp: $3})}
                    |   integer '+' REGISTER                    {$$ = Operand(&MemoryOperand {coord: yyS[yypt-1].coord, reg: $3, disp: $1})}
                    |   integer                                 {$$ = Operand(&MemoryOperand {coord: yyS[yypt-1].coord, reg: NoRegister, disp: $1})}

integer:                INTEGER                                 {$$ = Literal(&ConstantLiteral {coord: yyS[yypt-1].coord, value: uint32($1)})}
                    |   IDENTIFIER                              {$$ = Literal(&LabelLiteral    {coord: yyS[yypt-1].coord, name: $1})}

%%
