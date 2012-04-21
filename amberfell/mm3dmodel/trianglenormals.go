package mm3dmodel

import (
    
)

type TriangleNormals struct {
    flags uint16
    triangleIndex uint32
    triangle *Triangle
    v1x float32
    v1y float32
    v1z float32
    v2x float32
    v2y float32
    v2z float32
    v3x float32
    v3y float32
    v3z float32
}

func (triangleNormals *TriangleNormals) Flags() (flags uint16) {return triangleNormals.flags}
func (triangleNormals *TriangleNormals) TriangleIndex() (triangleIndex uint32) {return triangleNormals.triangleIndex}
func (triangleNormals *TriangleNormals) Triangle() (triangle *Triangle) {return triangleNormals.triangle}
func (triangleNormals *TriangleNormals) Vertex1Normal() (x float32, y float32, z float32) {return triangleNormals.v1x, triangleNormals.v1y, triangleNormals.v1z}
func (triangleNormals *TriangleNormals) Vertex2Normal() (x float32, y float32, z float32) {return triangleNormals.v2x, triangleNormals.v2y, triangleNormals.v2z}
func (triangleNormals *TriangleNormals) Vertex3Normal() (x float32, y float32, z float32) {return triangleNormals.v3x, triangleNormals.v3y, triangleNormals.v3z}

func (triangleNormals *TriangleNormals) Associate(model *Model) {
    if model.triangles != nil {
        triangleNormals.triangle = model.triangles[triangleNormals.triangleIndex]
        triangleNormals.triangle.triangleNormals = triangleNormals
    }
}
