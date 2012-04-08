
//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:2
    package main
    
    import (
        "bufio"
        "fmt"
        "os"
    )
    
    var parseResult Node

//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:13
type yySymType struct {
	yys int
    ival int
    sval string
    node Node
    nodes []Node
}

const COMMA = 57346
const COLON = 57347
const IDENTIFIER = 57348
const LBRACKET = 57349
const LITSTRING = 57350
const NEWLINE = 57351
const O = 57352
const PC = 57353
const PEEK = 57354
const PLUS = 57355
const POP = 57356
const PUSH = 57357
const RBRACKET = 57358
const SP = 57359
const REGISTER = 57360
const INTEGER = 57361

var yyToknames = []string{
	"COMMA",
	"COLON",
	"IDENTIFIER",
	"LBRACKET",
	"LITSTRING",
	"NEWLINE",
	"O",
	"PC",
	"PEEK",
	"PLUS",
	"POP",
	"PUSH",
	"RBRACKET",
	"SP",
	"REGISTER",
	"INTEGER",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:188


func parse(reader io.Reader) (asm Node) {
    yyParse(newLexer(bufio.NewReader(reader)))
    
    exitIfErrors()
    
    return parseResult
}

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 27
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 58

var yyAct = []int{

	13, 12, 21, 24, 22, 31, 24, 30, 32, 36,
	29, 35, 9, 7, 25, 26, 23, 27, 22, 23,
	3, 2, 10, 8, 4, 1, 6, 28, 5, 11,
	24, 14, 34, 33, 20, 19, 16, 0, 15, 17,
	0, 18, 22, 23, 24, 14, 0, 0, 20, 19,
	16, 0, 15, 17, 0, 18, 22, 23,
}
var yyPact = []int{

	7, -1000, 7, -1000, 3, -1000, -1000, 24, -1000, -1000,
	10, -1000, -1000, -1000, 0, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 38, -6, -8, -1000, -1000,
	-3, -14, -1000, -5, -7, -1000, -1000,
}
var yyPgo = []int{

	0, 28, 2, 26, 1, 0, 25, 24, 20, 22,
	21,
}
var yyR1 = []int{

	0, 6, 10, 10, 8, 7, 7, 1, 1, 9,
	9, 4, 4, 4, 4, 4, 4, 4, 4, 4,
	4, 4, 4, 5, 3, 2, 2,
}
var yyR2 = []int{

	0, 1, 2, 1, 2, 1, 1, 2, 1, 3,
	1, 1, 3, 5, 5, 1, 1, 1, 1, 1,
	1, 3, 1, 1, 2, 1, 1,
}
var yyChk = []int{

	-1000, -6, -10, -8, -7, -1, -3, 6, -8, 9,
	-9, 5, -4, -5, 7, 14, 12, 15, 17, 11,
	10, -2, 18, 19, 6, 4, -5, -2, -4, 16,
	13, 13, 16, -2, -5, 16, 16,
}
var yyDef = []int{

	0, -2, 1, 3, 0, 5, 6, 8, 2, 4,
	7, 24, 10, 11, 0, 15, 16, 17, 18, 19,
	20, 22, 23, 25, 26, 0, 0, 0, 9, 12,
	0, 0, 21, 0, 0, 13, 14,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19,
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
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:54
		{
	        parseResult = NewAssembly(yyS[yypt-0].nodes)
	    }
	case 2:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:60
		{
	        yyVAL.nodes = append(yyS[yypt-1].nodes, yyS[yypt-0].node)
	    }
	case 3:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:65
		{
	        yyVAL.nodes = []Node{yyS[yypt-0].node}
	    }
	case 4:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:71
		{
	        yyVAL.node = yyS[yypt-1].node
	    }
	case 5:
		yyVAL.node = yyS[yypt-0].node
	case 6:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:78
		{
	        yyVAL.node = yyS[yypt-0].node
	    }
	case 7:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:84
		{
	        yyVAL.node = NewInstruction(yyS[yypt-1].sval, yyS[yypt-0].nodes)
	    }
	case 8:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:89
		{
	        yyVAL.node = NewInstruction(yyS[yypt-0].sval, []Node{})
	    }
	case 9:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:95
		{
	        yyVAL.nodes = append(yyS[yypt-2].nodes, yyS[yypt-0].node)
	    }
	case 10:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:100
		{
	        yyVAL.nodes = []Node{yyS[yypt-0].node}
	    }
	case 11:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:106
		{
	        yyVAL.node = NewOperand(O_REG, yyS[yypt-0].node, nil)
	    }
	case 12:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:111
		{
	        yyVAL.node = NewOperand(O_MEM, yyS[yypt-1].node, nil)
	    }
	case 13:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:116
		{
	        yyVAL.node = NewOperand(O_MEMDISP, yyS[yypt-3].node, yyS[yypt-1].node)
	    }
	case 14:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:121
		{
	        yyVAL.node = NewOperand(O_MEMDISP, yyS[yypt-1].node, yyS[yypt-3].node)
	    }
	case 15:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:126
		{
	        yyVAL.node = NewOperand(O_POP, nil, nil)
	    }
	case 16:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:131
		{
	        yyVAL.node = NewOperand(O_PEEK, nil, nil)
	    }
	case 17:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:136
		{
	        yyVAL.node = NewOperand(O_PUSH, nil, nil)
	    }
	case 18:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:141
		{
	        yyVAL.node = NewOperand(O_SP, nil, nil)
	    }
	case 19:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:146
		{
	        yyVAL.node = NewOperand(O_PC, nil, nil)
	    }
	case 20:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:151
		{
	        yyVAL.node = NewOperand(O_O, nil, nil)
	    }
	case 21:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:156
		{
	        yyVAL.node = NewOperand(O_MEMIMM, yyS[yypt-1].node, nil)
	    }
	case 22:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:161
		{
	        yyVAL.node = NewOperand(O_IMM, yyS[yypt-0].node, nil)
	    }
	case 23:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:167
		{
	        yyVAL.node = NewRegister(yyS[yypt-0].ival)
	    }
	case 24:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:173
		{
	        yyVAL.node = NewLabel(yyS[yypt-1].sval)
	    }
	case 25:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:179
		{
	        yyVAL.node = NewConstant(yyS[yypt-0].ival)
	    }
	case 26:
		//line ./src/github.com/kierdavis/go/_/dcpuasm/parser.y:184
		{
	        yyVAL.node = NewLabelRef(yyS[yypt-0].sval)
	    }
	}
	goto yystack /* stack new state and value */
}
