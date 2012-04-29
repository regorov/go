package mm3dmodel

type DirtySegment struct {
    index      int
    offsetType uint16
    flags      uint16
    elements   [][]byte
}

// Function Index returns the index of this dirty segment definition into its parent model's list.
func (ds *DirtySegment) Index() (index int) {
    return ds.index
}

func (ds *DirtySegment) OffsetType() (offsetType uint16) { return ds.offsetType }
func (ds *DirtySegment) Flags() (flags uint16)           { return ds.flags }
func (ds *DirtySegment) NElements() (num int)            { return len(ds.elements) }
func (ds *DirtySegment) Element(n int) (element []byte)  { return ds.elements[n] }
