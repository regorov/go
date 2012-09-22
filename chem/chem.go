package chem

import (
	"fmt"
	"math"
)

func periodGroup(num int) (period int, group int) {
	if num <= 2 {
		return 1, num
	}

	if num <= 10 {
		return 2, num - 2
	}

	if num <= 18 {
		return 3, num - 10
	}

	if num <= 36 {
		return 4, num - 18
	}

	if num <= 54 {
		return 5, num - 36
	}

	if num <= 86 {
		return 6, num - 54
	}

	return 7, num - 86
}

func digitPrefix(n int) (s string) {
	switch n {
	case 1:
		return "hen"
	case 2:
		return "do"
	case 3:
		return "tri"
	case 4:
		return "tetra"
	case 5:
		return "penta"
	case 6:
		return "hexa"
	case 7:
		return "hepta"
	case 8:
		return "octa"
	case 9:
		return "nona"
	}

	return ""
}

func hydrocarbonPrefix(n int) (s string) {
	switch n {
	case 1:
		return "meth"
	case 2:
		return "eth"
	case 3:
		return "prop"
	case 4:
		return "but"
	case 5:
		return "pent"
	case 6:
		return "hex"
	case 7:
		return "hept"
	case 8:
		return "oct"
	case 9:
		return "non"
	case 11:
		return "undec"
	case 20:
		return "icos"
	case 21:
		return "henicos"
	}

	if n <= 19 {
		return digitPrefix(n%10) + "dec"
	}

	if n <= 29 {
		return digitPrefix(n%10) + "cos"
	}

	if n <= 39 {
		return digitPrefix(n%10) + "triacont"
	}

	n -= 40
	return digitPrefix(n%10) + digitPrefix(n/10) + "cont"
}

type Element struct {
	Name   string
	Symbol string
	Number int
	Mass   float64
}

func (e Element) String() (s string) {
	return e.Symbol
}

func (e Element) RoundMass() (mass float64) {
	m := math.Mod(e.Mass, 1.0)

	if m < 0.45 {
		return e.Mass - m
	}

	if m < 0.55 {
		return e.Mass - m + 0.5
	}

	return e.Mass - m + 1.0
}

func (e Element) Period() (period int) {
	period, _ = periodGroup(e.Number)
	return period
}

func (e Element) Group() (group int) {
	_, group = periodGroup(e.Number)
	return group
}

type ElementCount struct {
	Element *Element
	Count   int
}

type Compound []ElementCount

func ParseCompound(s string) (c Compound) {
	c = make(Compound, 0)

	i := 0
	baseCount := 0

	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		baseCount = (baseCount * 10) + int(s[i]-'0')
		i++
	}

	if baseCount == 0 {
		baseCount = 1
	}

	for i < len(s) {
		char := s[i]
		i++

		if char < 'A' || char > 'Z' {
			continue
		}

		symbol := string(char)

		if i < len(s) && s[i] >= 'a' && s[i] <= 'z' {
			symbol += string(s[i])
			i++
		}

		element := ElementsBySymbol[symbol]
		count := 0

		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			count = (count * 10) + int(s[i]-'0')
			i++
		}

		if count == 0 {
			count = 1
		}

		c = append(c, ElementCount{element, count * baseCount})
	}

	return c
}

func (c Compound) String() (s string) {
	for _, ec := range c {
		if ec.Count != 0 {
			s += ec.Element.String()

			if ec.Count != 1 {
				s += fmt.Sprintf("%d", ec.Count)
			}
		}
	}

	return s
}

func (c Compound) Mass() (mass float64) {
	for _, ec := range c {
		mass += ec.Element.Mass * float64(ec.Count)
	}

	return mass
}

func (c Compound) RoundMass() (mass float64) {
	for _, ec := range c {
		mass += ec.Element.RoundMass() * float64(ec.Count)
	}

	return mass
}

func (c Compound) Name() (name string) {
	// Attempt to find the name of the compound.

	var (
		carbon   = ElementsBySymbol["C"]
		hydrogen = ElementsBySymbol["H"]
		oxygen   = ElementsBySymbol["O"]
	)

	counts := make(map[*Element]int)
	for _, ec := range c {
		counts[ec.Element] += ec.Count
	}

	// Is it empty?
	if len(counts) == 0 {
		return "N/A"
	}

	// Is it purely one element?
	if len(counts) == 1 {
		for element, _ := range counts {
			return element.Name
		}
	}

	// Some common ones
	if len(c) == 2 {
		if c[0] == (ElementCount{hydrogen, 2}) && c[1] == (ElementCount{oxygen, 1}) {
			return "water"
		}
		if c[0] == (ElementCount{carbon, 1}) && c[1] == (ElementCount{oxygen, 2}) {
			return "carbon dioxide"
		}
	}

	// Is it a hydrocarbon?
	if len(counts) == 2 && counts[hydrogen] > 0 && counts[carbon] > 0 {
		h := counts[hydrogen]
		c := counts[carbon]

		// Alkane?
		if h == (2*c)+2 {
			return hydrocarbonPrefix(c) + "ane"
		}

		// Alkene?
		if h == 2*c {
			return hydrocarbonPrefix(c) + "ene"
		}

		// Alkyne?
		if h == (2*c)-2 {
			return hydrocarbonPrefix(c) + "yne"
		}
	}

	return c.String()
}

var ElementsByNumber = map[int]*Element{
	1:  &Element{"hydrogen", "H", 1, 1.008},
	2:  &Element{"helium", "He", 2, 4.0026},
	3:  &Element{"lithium", "Li", 3, 6.94},
	4:  &Element{"berylium", "Be", 4, 9.0122},
	5:  &Element{"boron", "B", 5, 10.81},
	6:  &Element{"carbon", "C", 6, 12.011},
	7:  &Element{"nitrogen", "N", 7, 14.007},
	8:  &Element{"oxygen", "O", 8, 15.999},
	9:  &Element{"fluorine", "F", 9, 18.998},
	10: &Element{"neon", "Ne", 10, 20.180},
	11: &Element{"sodium", "Na", 11, 22.990},
	12: &Element{"magnesium", "Mg", 12, 24.305},
	13: &Element{"aluminium", "Al", 13, 26.982},
	14: &Element{"silicon", "Si", 14, 28.085},
	15: &Element{"phosphorus", "P", 15, 30.974},
	16: &Element{"sulfur", "S", 16, 32.06},
	17: &Element{"chlorine", "Cl", 17, 35.45},
	18: &Element{"argon", "Ar", 18, 39.948},
	19: &Element{"potassium", "K", 19, 39.098},
	20: &Element{"calcium", "Ca", 20, 40.078},
	21: &Element{"scandium", "Sc", 21, 44.956},
	22: &Element{"titanium", "Ti", 22, 47.867},
	23: &Element{"vanadium", "V", 23, 50.942},
	24: &Element{"chromium", "Cr", 24, 51.996},
	25: &Element{"manganese", "Mn", 25, 54.938},
	26: &Element{"iron", "Fe", 26, 55.845},
	27: &Element{"cobalt", "Co", 27, 58.933},
	28: &Element{"nickel", "Ni", 28, 58.693},
	29: &Element{"copper", "Cu", 29, 63.546},
	30: &Element{"zinc", "Zn", 30, 65.38},
	31: &Element{"gallium", "Ga", 31, 69.723},
	32: &Element{"germanium", "Ge", 32, 72.63},
	33: &Element{"arsenic", "As", 33, 74.922},
	34: &Element{"selenium", "Se", 34, 78.96},
	35: &Element{"bromine", "Br", 35, 79.904},
	36: &Element{"krypton", "Kr", 36, 83.798},
}

var ElementsByName = make(map[string]*Element)
var ElementsBySymbol = make(map[string]*Element)

func init() {
	for _, element := range ElementsByNumber {
		ElementsByName[element.Name] = element
		ElementsBySymbol[element.Symbol] = element
	}
}
