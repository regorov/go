package musical

var (
    Major           = []uint8{2, 2, 1, 2, 2, 2, 1}
    Minor           = []uint8{2, 1, 2, 2, 1, 2, 2}
    MelodicMinor    = []uint8{2, 1, 2, 2, 2, 2, 1}
    HarmonicMinor   = []uint8{2, 1, 2, 2, 1, 3, 1}
    PentatonicMajor = []uint8{2, 2, 3, 2, 3}
    BluesMajor      = []uint8{3, 2, 1, 1, 2, 3}
    PentatonicMinor = []uint8{3, 2, 2, 3, 2}
    BluesMinor      = []uint8{3, 2, 1, 1, 3, 2}
    Augmented       = []uint8{3, 1, 3, 1, 3, 1}
    Diminished      = []uint8{2, 1, 2, 1, 2, 1, 2, 1}
    Chromatic       = []uint8{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
    WholeHalf       = []uint8{2, 1, 2, 1, 2, 1, 2, 1}
    HalfWhole       = []uint8{1, 2, 1, 2, 1, 2, 1, 2}
    WholeTone       = []uint8{2, 2, 2, 2, 2, 2}
    AugmentedFifth  = []uint8{2, 2, 1, 2, 1, 1, 2, 1}
    Japanese        = []uint8{1, 4, 2, 1, 4}
    Oriental        = []uint8{1, 3, 1, 1, 3, 1, 2}
    Ionian          = []uint8{2, 2, 1, 2, 2, 2, 1}
    Dorian          = []uint8{2, 1, 2, 2, 2, 1, 2}
    Phrygian        = []uint8{1, 2, 2, 2, 1, 2, 2}
    Lydian          = []uint8{2, 2, 2, 1, 2, 2, 1}
    Mixolydian      = []uint8{2, 2, 1, 2, 2, 1, 2}
    Aeolian         = []uint8{2, 1, 2, 2, 1, 2, 2}
    Locrian         = []uint8{1, 2, 2, 1, 2, 2, 2}
)

type Scale struct {
    Root      Note
    Intervals []uint8
}

func NewScale(root Note, intervals []uint8) (scale Scale) {
    return Scale{
        Root:      root.AtOctave(0),
        Intervals: intervals,
    }
}

func (scale Scale) Get(index int) (note Note) {
    intervals := scale.Intervals
    note = scale.Root

    if index > 0 {
        x := 0

        for i := 0; i < index; i++ {
            note = note.Transpose(int(intervals[x]))
            x = (x + 1) % len(intervals)
        }

    } else {
        x := 0

        for i := 0; i < index; i++ {
            x = (x - 1) % len(intervals)
            note = note.Transpose(-int(intervals[x]))
        }
    }

    return note
}

func (scale Scale) Index(note Note) (index int) {
    intervals := scale.Intervals
    index = 0
    x := scale.Root
    i := 0

    for x < note {
        x = x.Transpose(int(intervals[i]))
        i = (i + 1) % len(intervals)
    }

    if x == note {
        return index
    }

    return -1
}

func (scale Scale) Transpose(note Note, interval int) (res Note) {
    return scale.Get(scale.Index(note) + interval)
}
