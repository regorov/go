package main

import (
	"fmt"
	"github.com/kierdavis/argparse"
	"github.com/kierdavis/goutil"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

func IsDigit(c rune) (ok bool) {
	return '0' <= c && c <= '9'
}

type TokenType int

const (
	Rune TokenType = iota
	FixedDigits
	AnyDigits
)

type Token struct {
	Type TokenType
	Data int
}

type Pattern []Token

func ParsePattern(s string) (pattern Pattern) {
	pattern = make(Pattern, 0)

	ch := util.DecodeUTF8Iter(s)
	c := <-ch

	for c != 0 {
		if c == '?' {
			n := 0
			for c == '?' {
				n++
				c = <-ch
			}

			pattern = append(pattern, Token{FixedDigits, n})
			continue
		}

		if c == '*' {
			c = <-ch
			pattern = append(pattern, Token{AnyDigits, int(c)})

		} else {
			pattern = append(pattern, Token{Rune, int(c)})
			c = <-ch
		}
	}

	return pattern
}

func (pattern Pattern) NumGroups() (n int) {
	for _, token := range pattern {
		if token.Type == FixedDigits || token.Type == AnyDigits {
			n++
		}
	}

	return n
}

func (pattern Pattern) Make(groups []int) (s string) {
	i := 0

	for _, token := range pattern {
		switch token.Type {
		case Rune:
			s += util.EncodeUTF8Rune(rune(token.Data))

		case FixedDigits:
			format := fmt.Sprintf("%%0%dd", token.Data)
			s += fmt.Sprintf(format, groups[i])
			i++

		case AnyDigits:
			s += fmt.Sprintf("%d", groups[i])
			i++
		}
	}

	return s
}

func (pattern Pattern) Match(s string) (groups []int) {
	groups = make([]int, 0)

	ch := util.DecodeUTF8Iter(s)
	c := <-ch

	for _, token := range pattern {
		switch token.Type {
		case Rune:
			if c != rune(token.Data) {
				return nil
			}
			c = <-ch

		case FixedDigits:
			n := 0

			for i := 0; i < token.Data; i++ {
				if !IsDigit(c) {
					return nil
				}

				n = (n * 10) + (int(c) - '0')
				c = <-ch
			}

			groups = append(groups, n)

		case AnyDigits:
			n := 0

			for c != rune(token.Data) {
				if !IsDigit(c) {
					return nil
				}

				n = (n * 10) + (int(c) - '0')
				c = <-ch
			}

			groups = append(groups, n)
		}
	}

	// ch should be closed now, so the last <-ch should have returned 0
	if c != 0 {
		return nil
	}

	return groups
}

type Action int

const (
	NoAction Action = iota
	List
	Next
)

type Args struct {
	Action  Action
	Pattern string
	Dir     string
	Zeros   bool
}

type File struct {
	Name   string
	Groups []int
}

type Files struct {
	Files     []File
	NumGroups int
}

func (files Files) Len() (n int) {
	return len(files.Files)
}

func (files Files) Less(i, j int) (less bool) {
	for k := 0; k < files.NumGroups; k++ {
		a := files.Files[i].Groups[k]
		b := files.Files[j].Groups[k]

		if a < b {
			return true
		}
		if a > b {
			return false
		}
	}

	return false
}

func (files Files) Swap(i, j int) {
	files.Files[i], files.Files[j] = files.Files[j], files.Files[i]
}

func ListDir(dir string) (ch chan string) {
	ch = make(chan string)

	go func() {
		l, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		for _, fi := range l {
			ch <- filepath.Join(dir, fi.Name())
		}

		close(ch)
	}()

	return ch
}

func main() {
	p := argparse.New("")

	p.Argument("Pattern", 1, argparse.Store, "PATTERN", "The pattern specifying which filenames should be matched. It is similar to a shell pattern: '*' matches any string of digits, '?' matches exactly one digit and any other character matches itself.")
	p.Argument("Dir", 1, argparse.Store, "DIRECTORY", "The directory to search in.")

	p.Option('l', "list", "Action", 0, argparse.StoreConst(List), "", "List all files in the directory that match this pattern, sorted in numerical order.")
	p.Option('n', "next", "Action", 0, argparse.StoreConst(Next), "", "Return the next logical filename in the directory (i.e. a filename that fits the pattern, with the first group set to the maximum value found in the directory, plus one.)")
	p.Option('0', "zeros", "Zeros", 0, argparse.StoreConst(true), "", "Separate files produced by -l/--list with nulls (0-bytes) instead of newlines (for sending output to xargs -0 etc.)")

	args := &Args{}
	err := p.Parse(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	switch args.Action {
	case NoAction:
		p.Usage()
		fmt.Fprintf(os.Stderr, "No action was given (try one of -l, -n)\n")
		os.Exit(2)

	case List:
		pattern := ParsePattern(args.Pattern)
		files := make([]File, 0)

		for filename := range ListDir(args.Dir) {
			match := pattern.Match(filename)
			if match != nil {
				files = append(files, File{filename, match})
			}
		}

		f := Files{files, pattern.NumGroups()}
		sort.Sort(f)

		for _, file := range f.Files {
			if args.Zeros {
				fmt.Print(file.Name + "\x00")
			} else {
				fmt.Println(file.Name)
			}
		}

	case Next:
		n := 0
		pattern := ParsePattern(args.Pattern)
		if pattern.NumGroups() < 1 {
			p.Usage()
			fmt.Fprintf(os.Stderr, "At least 1 group must be specified\n")
			os.Exit(2)
		}

		for filename := range ListDir(args.Dir) {
			match := pattern.Match(filename)
			if match != nil && match[0] > n {
				n = match[0]
			}
		}

		groups := make([]int, pattern.NumGroups())
		groups[0] = n + 1

		fmt.Println(pattern.Make(groups))
	}
}
