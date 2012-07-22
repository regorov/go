//line src/github.com/kierdavis/go/k920/k920asmlib/parser.y:3
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

//line src/github.com/kierdavis/go/k920/k920asmlib/parser.y:23
type yySymType struct {
	yys int
	s   string
	i   int64
	f   float64
	R   Register
	BR  BitReg
	O   Operand
	OL  []Operand
}

const EOF = 57346
const NL = 57347
const IDENTIFIER = 57348
const INTEGER = 57349
const FLOAT = 57350
const REGISTER = 57351
const BITREG = 57352

var yyToknames = []string{
	"EOF",
	"NL",
	"IDENTIFIER",
	"INTEGER",
	"FLOAT",
	"REGISTER",
	"BITREG",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line src/github.com/kierdavis/go/k920/k920asmlib/parser.y:95

type lexer struct {
	input        *bufio.Reader
	currentToken []rune
	lastColumn   int
	lineno       int
	column       int
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
		input:  bufio.NewReader(input),
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
	for f(ll.Next()) {
	}

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

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 16
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 24

var yyAct = []int{

	18, 17, 14, 15, 16, 12, 18, 17, 6, 15,
	16, 7, 5, 6, 10, 9, 19, 3, 4, 2,
	8, 1, 11, 13,
}
var yyPact = []int{

	2, -1000, 7, -1000, 10, 9, -6, -1000, -1000, -1000,
	-1000, -1000, -1000, 0, -1000, -1000, -1000, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 2, 23, 22, 21, 19, 17, 18, 12,
}
var yyR1 = []int{

	0, 4, 5, 5, 6, 6, 7, 8, 3, 3,
	2, 2, 1, 1, 1, 1,
}
var yyR2 = []int{

	0, 2, 2, 1, 2, 2, 2, 2, 1, 0,
	2, 1, 1, 1, 1, 1,
}
var yyChk = []int{

	-1000, -4, -5, -6, -7, -8, 6, 4, -6, 5,
	5, -3, 11, -2, -1, 9, 10, 7, 6, -1,
}
var yyDef = []int{

	0, -2, 0, 3, 0, 0, 9, 1, 2, 4,
	5, 6, 7, 8, 11, 12, 13, 14, 15, 10,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 11,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10,
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
		//line src/github.com/kierdavis/go/k920/k920asmlib/parser.y:45
		{
			return 0
		}
	case 6:
		//line src/github.com/kierdavis/go/k920/k920asmlib/parser.y:54
		{
			fullName := yyS[yypt-1].s
			p := strings.Index(fullName, ".")

			var group, name string

			if p < 0 {
				group = "std"
				name = fullName
			} else {
				group = fullName[:p]
				name = fullName[p+1:]
			}

			parserOutput <- &Instruction{Group: group, Name: name, Operands: yyS[yypt-0].OL}
		}
	case 7:
		//line src/github.com/kierdavis/go/k920/k920asmlib/parser.y:72
		{
			name := yyS[yypt-1].s

			if name[0] == '.' {
				name = parserBaseLabel + name
			} else {
				parserBaseLabel = name
			}

			parserOutput <- &Label{Name: name}
		}
	case 8:
		//line src/github.com/kierdavis/go/k920/k920asmlib/parser.y:84
		{
			yyVAL.OL = yyS[yypt-0].OL
		}
	case 9:
		//line src/github.com/kierdavis/go/k920/k920asmlib/parser.y:85
		{
			yyVAL.OL = []Operand{}
		}
	case 10:
		//line src/github.com/kierdavis/go/k920/k920asmlib/parser.y:87
		{
			yyVAL.OL = append(yyS[yypt-1].OL, yyS[yypt-0].O)
		}
	case 11:
		//line src/github.com/kierdavis/go/k920/k920asmlib/parser.y:88
		{
			yyVAL.OL = []Operand{yyS[yypt-0].O}
		}
	case 12:
		//line src/github.com/kierdavis/go/k920/k920asmlib/parser.y:90
		{
			yyVAL.O = yyS[yypt-0].R
		}
	case 13:
		//line src/github.com/kierdavis/go/k920/k920asmlib/parser.y:91
		{
			yyVAL.O = yyS[yypt-0].BR
		}
	case 14:
		//line src/github.com/kierdavis/go/k920/k920asmlib/parser.y:92
		{
			yyVAL.O = Integer(yyS[yypt-0].i)
		}
	case 15:
		//line src/github.com/kierdavis/go/k920/k920asmlib/parser.y:93
		{
			yyVAL.O = LabelRef(yyS[yypt-0].s)
		}
	}
	goto yystack /* stack new state and value */
}
