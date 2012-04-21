package mm3dmodel

import (
    
)

type Vertex struct {
    flags uint16
    x float32
    y float32
    z float32
}

func (vertex *Vertex) Flags() (flags uint16) {return vertex.flags}
func (vertex *Vertex) X() (x float32) {return vertex.x}
func (vertex *Vertex) Y() (y float32) {return vertex.y}
func (vertex *Vertex) Z() (z float32) {return vertex.z}
func (vertex *Vertex) Hidden() (hidden bool) {return vertex.flags & 0x0001 != 0}
func (vertex *Vertex) Selected() (selected bool) {return vertex.flags & 0x0002 != 0}
func (vertex *Vertex) Free() (free bool) {return vertex.flags & 0x0004 != 0}
func (vertex *Vertex) Associate(model *Model) {}