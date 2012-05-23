//line src/github.com/kierdavis/go/k750/k750asm/parser.y:2
package main

import (
	"fmt"
	"log"
)

//line src/github.com/kierdavis/go/k750/k750asm/parser.y:10
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

//line src/github.com/kierdavis/go/k750/k750asm/parser.y:72

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 21
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 34

var yyAct = []int{

	12, 26, 11, 25, 24, 15, 18, 13, 16, 9,
	23, 17, 21, 16, 15, 23, 13, 16, 16, 22,
	19, 5, 28, 7, 3, 2, 27, 6, 1, 8,
	10, 20, 14, 4,
}
var yyPact = []int{

	14, -1000, 14, -1000, 18, 1, -1000, -1000, -1000, -1000,
	2, -1000, -1000, -1000, -1000, -4, -1000, 10, 6, -1000,
	-7, -9, -11, -1000, -1000, 11, 16, -1000, -1000,
}
var yyPgo = []int{

	0, 33, 32, 31, 2, 30, 29, 0, 28, 25,
	24,
}
var yyR1 = []int{

	0, 8, 9, 9, 10, 1, 1, 6, 6, 5,
	5, 4, 4, 4, 2, 3, 3, 3, 3, 7,
	7,
}
var yyR2 = []int{

	0, 1, 2, 1, 2, 2, 2, 1, 0, 3,
	1, 1, 1, 1, 4, 1, 3, 3, 1, 1,
	1,
}
var yyChk = []int{

	-1000, -8, -9, -10, -1, 7, -10, 5, -6, 8,
	-5, -4, -7, 6, -2, 4, 7, 9, 10, -4,
	-3, 6, -7, 4, 11, 12, 12, -7, 6,
}
var yyDef = []int{

	0, -2, 1, 3, 0, 8, 2, 4, 5, 6,
	7, 10, 11, 12, 13, 19, 20, 0, 0, 9,
	0, 15, 18, 19, 14, 0, 0, 16, 17,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 12, 9, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 8, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 10, 3, 11,
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
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:33
		{
			close(parserOutput)
		}
	case 4:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:38
		{
			parserOutput <- yyS[yypt-1].it
		}
	case 5:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:40
		{
			yyVAL.it = Item(&Instruction{coord: yyS[yypt-1].coord, name: yyS[yypt-1].s, operands: yyS[yypt-0].oL})
		}
	case 6:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:41
		{
			yyVAL.it = Item(&Label{coord: yyS[yypt-1].coord, name: yyS[yypt-1].s})
		}
	case 7:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:43
		{
			yyVAL.oL = yyS[yypt-0].oL
		}
	case 8:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:44
		{
			yyVAL.oL = nil
		}
	case 9:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:46
		{
			yyVAL.oL = append(yyS[yypt-2].oL, yyS[yypt-0].o)
		}
	case 10:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:47
		{
			yyVAL.oL = []Operand{yyS[yypt-0].o}
		}
	case 11:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:49
		{
			yyVAL.o = Operand(&LiteralOperand{coord: yyS[yypt-1].coord, Literal: yyS[yypt-0].l})
		}
	case 12:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:50
		{
			yyVAL.o = Operand(&RegisterOperand{coord: yyS[yypt-1].coord, num: yyS[yypt-0].r})
		}
	case 13:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:51
		{
			yyVAL.o = yyS[yypt-0].o
		}
	case 14:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:54
		{
			size := yyS[yypt-3].i
			if size != 8 && size != 16 && size != 32 {
				log.Fatalf("Invalid memory addressing size: %d (expected 8, 16 or 32)", size)
			}

			yyVAL.o = yyS[yypt-1].o
			yyVAL.o.SetSize(MemSize(size))
		}
	case 15:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:64
		{
			yyVAL.o = Operand(&MemoryOperand{coord: yyS[yypt-1].coord, reg: yyS[yypt-0].r, disp: Zero})
		}
	case 16:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:65
		{
			yyVAL.o = Operand(&MemoryOperand{coord: yyS[yypt-1].coord, reg: yyS[yypt-2].r, disp: yyS[yypt-0].l})
		}
	case 17:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:66
		{
			yyVAL.o = Operand(&MemoryOperand{coord: yyS[yypt-1].coord, reg: yyS[yypt-0].r, disp: yyS[yypt-2].l})
		}
	case 18:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:67
		{
			yyVAL.o = Operand(&MemoryOperand{coord: yyS[yypt-1].coord, reg: NoRegister, disp: yyS[yypt-0].l})
		}
	case 19:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:69
		{
			yyVAL.l = Literal(&ConstantLiteral{coord: yyS[yypt-1].coord, value: uint32(yyS[yypt-0].i)})
		}
	case 20:
		//line src/github.com/kierdavis/go/k750/k750asm/parser.y:70
		{
			yyVAL.l = Literal(&LabelLiteral{coord: yyS[yypt-1].coord, name: yyS[yypt-0].s})
		}
	}
	goto yystack /* stack new state and value */
}
