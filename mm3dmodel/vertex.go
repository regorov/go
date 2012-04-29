package mm3dmodel

type Vertex struct {
    index int
    flags uint16
    x     float32
    y     float32
    z     float32
}

// Function Index returns the index of this definiton into its parent model's list.
func (vertex *Vertex) Index() (index int) {
    return vertex.index
}

func (vertex *Vertex) Flags() (flags uint16)     { return vertex.flags }
func (vertex *Vertex) X() (x float32)            { return vertex.x }
func (vertex *Vertex) Y() (y float32)            { return vertex.y }
func (vertex *Vertex) Z() (z float32)            { return vertex.z }
func (vertex *Vertex) Hidden() (hidden bool)     { return vertex.flags&0x0001 != 0 }
func (vertex *Vertex) Selected() (selected bool) { return vertex.flags&0x0002 != 0 }
func (vertex *Vertex) Free() (free bool)         { return vertex.flags&0x0004 != 0 }
func (vertex *Vertex) Associate(model *Model)    {}
