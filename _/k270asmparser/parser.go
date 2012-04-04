
//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:2
    package k270asmparser
    
    import (
        "fmt"
    )

//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:9
type yySymType struct {
	yys int
    node Node
    nodes []Node
    sval string
    ival int
}

const COMMA = 57346
const COLON = 57347
const IDENTIFIER = 57348
const LITSTRING = 57349
const NEWLINE = 57350
const S8 = 57351
const S16 = 57352
const S32 = 57353
const S64 = 57354
const STRING = 57355
const U8 = 57356
const U16 = 57357
const U32 = 57358
const U64 = 57359
const INTEGER = 57360
const REGISTER = 57361

var yyToknames = []string{
	"COMMA",
	"COLON",
	"IDENTIFIER",
	"LITSTRING",
	"NEWLINE",
	"S8",
	"S16",
	"S32",
	"S64",
	"STRING",
	"U8",
	"U16",
	"U32",
	"U64",
	"INTEGER",
	"REGISTER",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:213


//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 34
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 57

var yyAct = []int{

	22, 19, 36, 18, 27, 49, 27, 43, 27, 37,
	32, 50, 44, 38, 29, 30, 26, 23, 26, 23,
	26, 9, 16, 28, 14, 31, 48, 3, 12, 13,
	15, 42, 34, 35, 47, 41, 33, 2, 39, 40,
	17, 10, 4, 7, 45, 46, 8, 11, 1, 21,
	51, 52, 20, 25, 6, 24, 5,
}
var yyPact = []int{

	15, -1000, 15, -1000, 14, -1000, -1000, -1000, -1000, -2,
	-1000, -1000, 16, 2, 2, -1000, -1000, 21, -1000, -1000,
	-1000, -1000, -1000, 5, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 0, -17, -6, 3, -1000, -1000, 2, 2, -1000,
	-1000, -9, 1, 2, 2, -1000, -1000, -12, -1, 2,
	2, -1000, -1000,
}
var yyPgo = []int{

	0, 56, 0, 55, 54, 53, 1, 52, 49, 48,
	47, 46, 43, 42, 27, 41, 40, 37, 36, 35,
	34, 32, 31, 26,
}
var yyR1 = []int{

	0, 9, 17, 17, 14, 13, 13, 13, 13, 1,
	1, 16, 16, 6, 6, 6, 7, 8, 4, 12,
	12, 18, 19, 20, 15, 21, 22, 23, 10, 11,
	2, 2, 3, 5,
}
var yyR2 = []int{

	0, 1, 2, 1, 2, 1, 1, 1, 1, 2,
	1, 3, 1, 1, 1, 1, 1, 3, 2, 1,
	1, 0, 0, 0, 11, 0, 0, 0, 11, 2,
	1, 1, 1, 1,
}
var yyChk = []int{

	-1000, -9, -17, -14, -13, -1, -4, -12, -11, 6,
	-15, -10, 13, 14, 9, -14, 8, -16, 5, -6,
	-7, -8, -2, 19, -3, -5, 18, 6, 7, -2,
	-2, 4, 5, -18, -21, -6, 19, 15, 10, -2,
	-2, -19, -22, 16, 11, -2, -2, -20, -23, 17,
	12, -2, -2,
}
var yyDef = []int{

	0, -2, 1, 3, 0, 5, 6, 7, 8, 10,
	19, 20, 0, 0, 0, 2, 4, 9, 18, 12,
	13, 14, 15, 16, 30, 31, 32, 33, 29, 21,
	25, 0, 0, 0, 0, 11, 17, 0, 0, 22,
	26, 0, 0, 0, 0, 23, 27, 0, 0, 0,
	0, 24, 28,
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
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:57
		{
	        yyVAL.node = Node(NewAssembly(yyS[yypt-0].nodes))
	    }
	case 2:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:63
		{
	        yyVAL.nodes = append(yyS[yypt-1].nodes, yyS[yypt-0].node)
	    }
	case 3:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:68
		{
	        yyVAL.nodes = []Node{yyS[yypt-0].node}
	    }
	case 4:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:74
		{
	        yyVAL.node = yyS[yypt-1].node
	    }
	case 5:
		yyVAL.node = yyS[yypt-0].node
	case 6:
		yyVAL.node = yyS[yypt-0].node
	case 7:
		yyVAL.node = yyS[yypt-0].node
	case 8:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:83
		{
	        yyVAL.node = yyS[yypt-0].node
	    }
	case 9:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:89
		{
	        yyVAL.node = Node(NewInstruction(yyS[yypt-1].sval, yyS[yypt-0].nodes))
	    }
	case 10:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:94
		{
	        yyVAL.node = Node(NewInstruction(yyS[yypt-0].sval, nil))
	    }
	case 11:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:100
		{
	        yyVAL.nodes = append(yyS[yypt-2].nodes, yyS[yypt-0].node)
	    }
	case 12:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:105
		{
	        yyVAL.nodes = []Node{yyS[yypt-0].node}
	    }
	case 13:
		yyVAL.node = yyS[yypt-0].node
	case 14:
		yyVAL.node = yyS[yypt-0].node
	case 15:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:113
		{
	        yyVAL.node = yyS[yypt-0].node
	    }
	case 16:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:119
		{
	        yyVAL.node = Node(NewRegister(yyS[yypt-0].ival))
	    }
	case 17:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:125
		{
	        if !(yyS[yypt-2].ival % 2 == 0 && yyS[yypt-0].ival - yyS[yypt-2].ival == 1) {
	            yyerror("Invalid register pair combination")
	        }
	        
	        yyVAL.node = Node(NewRegisterPair(yyS[yypt-2].ival))
	    }
	case 18:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:135
		{
	        yyVAL.node = Node(NewLabel(yyS[yypt-1].sval))
	    }
	case 19:
		yyVAL.node = yyS[yypt-0].node
	case 20:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:142
		{
	        yyVAL.node = yyS[yypt-0].node
	    }
	case 21:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:148
		{
	        yyVAL.node = Node(NewU8Data(yyS[yypt-0].node))
	    }
	case 22:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:153
		{
	        yyVAL.node = Node(NewU16Data(yyS[yypt-3].node))
	    }
	case 23:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:158
		{
	        yyVAL.node = Node(NewU32Data(yyS[yypt-6].node))
	    }
	case 24:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:163
		{
	        yyVAL.node = Node(NewU64Data(yyS[yypt-9].node))
	    }
	case 25:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:169
		{
	        yyVAL.node = Node(NewS8Data(yyS[yypt-0].node))
	    }
	case 26:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:174
		{
	        yyVAL.node = Node(NewS16Data(yyS[yypt-3].node))
	    }
	case 27:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:179
		{
	        yyVAL.node = Node(NewS32Data(yyS[yypt-6].node))
	    }
	case 28:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:184
		{
	        yyVAL.node = Node(NewS64Data(yyS[yypt-9].node))
	    }
	case 29:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:190
		{
	        yyVAL.node = Node(NewStringData(yyS[yypt-0].sval))
	    }
	case 30:
		yyVAL.node = yyS[yypt-0].node
	case 31:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:197
		{
	        yyVAL.node = yyS[yypt-0].node
	    }
	case 32:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:203
		{
	        yyVAL.node = Node(NewInteger(yyS[yypt-0].ival))
	    }
	case 33:
		//line ./src/github.com/kierdavis/go/_/k270asmparser/parser.y:209
		{
	        yyVAL.node = Node(NewLabelRef(yyS[yypt-0].sval))
	    }
	}
	goto yystack /* stack new state and value */
}
