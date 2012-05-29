package musical

import (
	"fmt"
	"math"
	"strings"
)

var NoteStrings = [...]string{
	"c", "c#", "d", "d#", "e", "f",
	"f#", "g", "g#", "a", "a#", "b",
}

type Note int

func ParseNote(s string) (note Note) {
	s = strings.ToLower(strings.Trim(s, " \t\n\r"))
	last := s[len(s)-1]
	note = 4 * 12

	if last >= '0' && last <= '9' {
		s = s[:len(s)-1]
		note = Note((last - '0') * 12)
	}

	if len(s) == 0 {
		return 0
	}

	switch s[0] {
	case 'c':
		note += 0
	case 'd':
		note += 2
	case 'e':
		note += 4
	case 'f':
		note += 5
	case 'g':
		note += 7
	case 'a':
		note += 9
	case 'b':
		note += 11
	}

	for _, accidental := range s[1:] {
		if accidental == '#' {
			note++
		} else if accidental == 'b' {
			note--
		}
	}

	return note
}

func (note Note) String() (str string) {
	return fmt.Sprintf("%s%d", note.Note(), note.Octave())
}

func (note Note) Note() (s string) {
	return NoteStrings[note%12]
}

func (note Note) Octave() (octave int) {
	return int(note / 12)
}

func (note Note) AtOctave(octave int) (res Note) {
	return (note % 12) + Note(octave*12)
}

func (note Note) Transpose(halfsteps int) (res Note) {
	return note + Note(halfsteps)
}

func (note Note) Frequency() (freq float64) {
	return 16.35159783128741 * math.Pow(2.0, float64(note)/12.0)
}
