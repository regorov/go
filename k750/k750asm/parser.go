//line src/github.com/kierdavis/go/k750/k750asm/parser.y:2
package main

import (
	"fmt"
)

//line src/github.com/kierdavis/go/k750/k750asm/parser.y:9
type yySymType struct {
	yys int
	i   int
	r   Register
	s   string
	it  Item
	o   Operand
	oL  []Operand
	l   Literal

	coord Coord
}

const INTEGER = 57346
const NL = 57347
const REGISTER = 57348
const IDENTIFIER = 57349

var yyToknames = []string{
	"INTEGER",
	"NL",
	"REGISTER",
	"IDENTIFIER",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line src/github.com/kierdavis/go/k750/k750asm/parser.y:54

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 15
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 20

var yyAct = []int{

	11, 14, 16, 13, 15, 9, 14, 5, 13, 15,
	7, 3, 2, 1, 6, 12, 8, 17, 10, 4,
}
var yyPact = []int{

	0, -1000, 0, -1000, 5, -3, -1000, -1000, -1000, -1000,
	-7, -1000, -1000, -1000, -1000, -1000, 2, -1000,
}
var yyPgo = []int{

	0, 19, 0, 18, 16, 15, 13, 12, 11,
}
var yyR1 = []int{

	0, 6, 7, 7, 8, 1, 1, 4, 4, 3,
	3, 2, 2, 5, 5,
}
var yyR2 = []int{

	0, 1, 2, 1, 2, 2, 2, 1, 0, 3,
	1, 1, 1, 1, 1,
}
var yyChk = []int{

	-1000, -6, -7, -8, -1, 7, -8, 5, -4, 8,
	-3, -2, -5, 6, 4, 7, 9, -2,
}
var yyDef = []int{

	0, -2, 1, 3, 0, 8, 2, 4, 5, 6,
	7, 10, 11, 12, 13, 14, 0, 9,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 9, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 8,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c > 0 && c <= len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return fmt.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return fmt.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		fmt.Printf("lex %U %s\n", uint(char), yyTokname(c))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		fmt.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				fmt.Printf("%s", yyStatname(yystate))
				fmt.Printf("saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					fmt.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				fmt.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		fmt.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:32
		{
			close(parserOutput)
		}
	case 4:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:37
		{
			parserOutput <- yyS[yypt-1].it
		}
	case 5:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:39
		{
			yyVAL.it = Item(&Instruction{coord: yyS[yypt-1].coord, name: yyS[yypt-1].s, operands: yyS[yypt-0].oL})
		}
	case 6:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:40
		{
			yyVAL.it = Item(&Label{coord: yyS[yypt-1].coord, name: yyS[yypt-1].s})
		}
	case 7:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:42
		{
			yyVAL.oL = yyS[yypt-0].oL
		}
	case 8:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:43
		{
			yyVAL.oL = nil
		}
	case 9:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:45
		{
			yyVAL.oL = append(yyS[yypt-2].oL, yyS[yypt-0].o)
		}
	case 10:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:46
		{
			yyVAL.oL = []Operand{yyS[yypt-0].o}
		}
	case 11:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:48
		{
			yyVAL.o = Operand(&LiteralOperand{Literal: yyS[yypt-0].l})
		}
	case 12:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:49
		{
			yyVAL.o = Operand(&RegisterOperand{num: yyS[yypt-0].r})
		}
	case 13:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:51
		{
			yyVAL.l = Literal(&ConstantLiteral{value: uint32(yyS[yypt-0].i)})
		}
	case 14:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:52
		{
			yyVAL.l = Literal(&LabelLiteral{coord: yyS[yypt-1].coord, name: yyS[yypt-0].s})
		}
	}
	goto yystack /* stack new state and value */
}
