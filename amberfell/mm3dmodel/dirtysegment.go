package mm3dmodel

import (
    
)

type DirtySegment struct {
    offsetType uint16
    flags uint16
    elements [][]byte
}

func (ds *DirtySegment) OffsetType() (offsetType uint16) {return ds.offsetType}
func (ds *DirtySegment) Flags() (flags uint16) {return ds.flags}
func (ds *DirtySegment) NElements() (num int) {return len(ds.elements)}
func (ds *DirtySegment) Element(n int) (element []byte) {return ds.elements[n]}
