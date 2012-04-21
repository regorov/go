package mm3dmodel

import (
    
)

type Triangle struct {
    flags uint16
    v1index uint32
    v2index uint32
    v3index uint32
    v1 *Vertex
    v2 *Vertex
    v3 *Vertex
    triangleNormals *TriangleNormals
    textureCoordinates *TextureCoordinates
}

func (triangle *Triangle) Flags() (flags uint16) {return triangle.flags}
func (triangle *Triangle) VertexIndex1() (vertexIndex1 uint32) {return triangle.v1index}
func (triangle *Triangle) VertexIndex2() (vertexIndex2 uint32) {return triangle.v2index}
func (triangle *Triangle) VertexIndex3() (vertexIndex3 uint32) {return triangle.v3index}
func (triangle *Triangle) Vertex1() (vertex1 *Vertex) {return triangle.v1}
func (triangle *Triangle) Vertex2() (vertex2 *Vertex) {return triangle.v2}
func (triangle *Triangle) Vertex3() (vertex3 *Vertex) {return triangle.v3}
func (triangle *Triangle) Hidden() (hidden bool) {return triangle.flags & 0x0001 != 0}
func (triangle *Triangle) Selected() (selected bool) {return triangle.flags & 0x0002 != 0}
func (triangle *Triangle) TriangleNormals() (triangleNormals *TriangleNormals) {return triangle.triangleNormals}
func (triangle *Triangle) TextureCoordinates() (textureCoordinates *TextureCoordinates) {return triangle.textureCoordinates}

func (triangle *Triangle) Associate(model *Model) {
    if model.vertices != nil {
        triangle.v1 = model.vertices[triangle.v1index]
        triangle.v2 = model.vertices[triangle.v2index]
        triangle.v3 = model.vertices[triangle.v3index]
    }
}
