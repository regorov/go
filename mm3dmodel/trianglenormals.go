package mm3dmodel

// Type TriangleNormals represents the normals for the three vertices of a triangle.
type TriangleNormals struct {
    index         int
    flags         uint16
    triangleIndex uint32
    triangle      *Triangle
    v1x           float32
    v1y           float32
    v1z           float32
    v2x           float32
    v2y           float32
    v2z           float32
    v3x           float32
    v3y           float32
    v3z           float32
}

// Function Index returns the index of this definiton into its parent model's list.
func (triangleNormals *TriangleNormals) Index() (index int) {
    return triangleNormals.index
}

// Function Flags returns the flags for the definition. As of MM3D version 1.6, there are no flags
// for this structure and so this value is reserved for future use.
func (triangleNormals *TriangleNormals) Flags() (flags uint16) {
    return triangleNormals.flags
}

// Function TriangleIndex returns the index into the associated model's triangles of the triangle
// that this structure describes.
func (triangleNormals *TriangleNormals) TriangleIndex() (triangleIndex uint32) {
    return triangleNormals.triangleIndex
}

// Function Triangle returns the triangle that this structure describes.
func (triangleNormals *TriangleNormals) Triangle() (triangle *Triangle) {
    return triangleNormals.triangle
}

// Function Vertex1Normal returns the normal of the first vertex, as a triplet of floats.
func (triangleNormals *TriangleNormals) Vertex1Normal() (x float32, y float32, z float32) {
    return triangleNormals.v1x, triangleNormals.v1y, triangleNormals.v1z
}

// Function Vertex2Normal returns the normal of the second vertex, as a triplet of floats.
func (triangleNormals *TriangleNormals) Vertex2Normal() (x float32, y float32, z float32) {
    return triangleNormals.v2x, triangleNormals.v2y, triangleNormals.v2z
}

// Function Vertex3Normal returns the normal of the third vertex, as a triplet of floats.
func (triangleNormals *TriangleNormals) Vertex3Normal() (x float32, y float32, z float32) {
    return triangleNormals.v3x, triangleNormals.v3y, triangleNormals.v3z
}

// Function Associate updates the links between this structure and others, using the model it was
// parsed from as a reference point.
func (triangleNormals *TriangleNormals) Associate(model *Model) {
    if model.triangles != nil {
        triangleNormals.triangle = model.triangles[triangleNormals.triangleIndex]
        triangleNormals.triangle.triangleNormals = triangleNormals
    }
}
