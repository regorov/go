package mm3dmodel

import (
    
)

type Model struct {
    majorVersion uint8
    minorVersion uint8
    modelFlags uint8
    dirtySegments []*DirtySegment
    
    metadata map[string]string
    vertices []*Vertex
}

func (model *Model) MajorVersion() (majorVersion uint8) {return model.majorVersion}
func (model *Model) MinorVersion() (minorVersion uint8) {return model.minorVersion}
func (model *Model) ModelFlags() (modelFlags uint8) {return model.modelFlags}
func (model *Model) NDirtySegments() (num int) {return len(model.dirtySegments)}
func (model *Model) DirtySegment(n int) (ds *DirtySegment) {return model.dirtySegments[n]}
func (model *Model) Metadata() (metadata map[string]string) {return model.metadata}
func (model *Model) MetadataValue(key string) (value string, ok bool) {x, ok = model.metadata[key]; return}
func (model *Model) NVertices() (num int) {return len(model.vertices)}
func (model *Model) Vertex(n int) (vertex *Vertex) {return model.vertices[n]}

type DirtySegment struct {
    offsetType uint16
    flags uint16
    elements [][]byte
}

func (ds *DirtySegment) OffsetType() (offsetType uint16) {return ds.offsetType}
func (ds *DirtySegment) Flags() (flags uint16) {return ds.flags}
func (ds *DirtySegment) NElements() (num int) {return len(ds.elements)}
func (ds *DirtySegment) Element(n int) (element []byte) {return ds.elements[n]}

type Vertex struct {
    flags uint16
    x float32
    y float32
    z float32
}

func (vertex *Vertex) Flags() (flags uint16) {return ds.flags}
func (vertex *Vertex) X() (x float32) {return ds.x}
func (vertex *Vertex) Y() (y float32) {return ds.y}
func (vertex *Vertex) Z() (z float32) {return ds.z}
func (vertex *Vertex) Hidden() (hidden bool) {return ds.flags & 0x0001 != 0}
func (vertex *Vertex) Selected() (selected bool) {return ds.flags & 0x0002 != 0}
func (vertex *Vertex) Free() (free bool) {return ds.flags & 0x0004 != 0}
