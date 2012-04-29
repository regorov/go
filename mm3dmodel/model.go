package mm3dmodel

import (
    "github.com/banthar/gl"
)

// Type model represents an MM3D model.
type Model struct {
    majorVersion  uint8
    minorVersion  uint8
    modelFlags    uint8
    dirtySegments []*DirtySegment

    metadata           map[string]string
    groups             []*Group
    vertices           []*Vertex
    triangles          []*Triangle
    triangleNormals    []*TriangleNormals
    textureCoordinates []*TextureCoordinates
}

// Function MajorVersion returns the major version number of the file's encoding.
func (model *Model) MajorVersion() (majorVersion uint8) {
    return model.majorVersion
}

// Function MinorVersion returns the minor version number of the file's encoding.
func (model *Model) MinorVersion() (minorVersion uint8) {
    return model.minorVersion
}

// Function ModelFlags returns the model-wide flags for the model. As of MM3D version 1.6, there are
// no model-wide flags and so this value is reserved for future use.
func (model *Model) ModelFlags() (modelFlags uint8) {
    return model.modelFlags
}

// Function NDirtySegments returns the number of dirty segments in the file (segments that were not
// recognised by the parser).
func (model *Model) NDirtySegments() (num int) {
    return len(model.dirtySegments)
}

// Function DirtySegment returns the dirty segment descriptor with the specified index.
func (model *Model) DirtySegment(n int) (ds *DirtySegment) {
    return model.dirtySegments[n]
}

// Function Metadata returns the metadata of the model, as a string to string map.
func (model *Model) Metadata() (metadata map[string]string) {
    return model.metadata
}

// Function MetadataValue returns the metadata value with the specified key.
func (model *Model) MetadataValue(key string) (value string, ok bool) {
    value, ok = model.metadata[key]
    return
}

// Function NVertices returns the number of vertices in the model.
func (model *Model) NVertices() (num int) {
    return len(model.vertices)
}

// Function Vertex returns the vertex with the specified index.
func (model *Model) Vertex(n int) (vertex *Vertex) {
    return model.vertices[n]
}

// Function NTriangles returns the number of triangles in the model.
func (model *Model) NTriangles() (num int) {
    return len(model.triangles)
}

// Function Triangle returns the triangle with the specified index.
func (model *Model) Triangle(n int) (triangle *Triangle) {
    return model.triangles[n]
}

// Function NTriangleNormals returns the number of triangle normal definitions in the model.
func (model *Model) NTriangleNormals() (num int) {
    return len(model.triangleNormals)
}

// Function TriangleNormals returns the triangle normal definition with the specified index.
func (model *Model) TriangleNormals(n int) (triangleNormals *TriangleNormals) {
    return model.triangleNormals[n]
}

// Function NTextureCoordinates returns the number of texture coordinates definitions in the model.
func (model *Model) NTextureCoordinates() (num int) {
    return len(model.textureCoordinates)
}

// Function TextureCoordinates returns the texture coordinates definition with the specified index.
func (model *Model) TextureCoordinates(n int) (textureCoordinates *TextureCoordinates) {
    return model.textureCoordinates[n]
}

// Function GLDraw draws the model into the current GL scene. It assumes that the context is already
// set up, including any translation and rotation.
func (model *Model) GLDraw() {
    var a, b, c float32
    var v *Vertex

    // Start drawing triangles
    gl.Begin(gl.TRIANGLES)

    for _, triangle := range model.triangles {
        triangleNormals := triangle.TriangleNormals()
        textureCoordinates := triangle.TextureCoordinates()

        // Set the normal. Assume it's a flat triangle (the normals at all 3 vertices are equal)
        a, b, c = triangleNormals.Vertex1Normal()
        gl.Normal3f(a, b, c)

        // Vertex 1
        a, b = textureCoordinates.Vertex1Coord()
        gl.TexCoord2f(a, b)
        v = triangle.Vertex1()
        gl.Vertex3f(v.X(), v.Y(), v.Z())

        // Vertex 2
        a, b = textureCoordinates.Vertex2Coord()
        gl.TexCoord2f(a, b)
        v = triangle.Vertex2()
        gl.Vertex3f(v.X(), v.Y(), v.Z())

        // Vertex 3
        a, b = textureCoordinates.Vertex3Coord()
        gl.TexCoord2f(a, b)
        v = triangle.Vertex3()
        gl.Vertex3f(v.X(), v.Y(), v.Z())
    }

    // Finish
    gl.End()
}
