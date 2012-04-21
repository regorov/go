package mm3dmodel

import (
    "github.com/banthar/gl"
)

// Type model represents an MM3D model.
type Model struct {
    majorVersion uint8
    minorVersion uint8
    modelFlags uint8
    dirtySegments []*DirtySegment
    
    metadata map[string]string
    vertices []*Vertex
    triangles []*Triangle
    triangleNormals []*TriangleNormals
    textureCoordinates []*TextureCoordinates
}

// Function MajorVersion returns the major version number of the file's encoding.
func (model *Model) MajorVersion() (majorVersion uint8) {return model.majorVersion}

// Function MinorVersion returns the minor version number of the file's encoding.
func (model *Model) MinorVersion() (minorVersion uint8) {return model.minorVersion}
func (model *Model) ModelFlags() (modelFlags uint8) {return model.modelFlags}
func (model *Model) NDirtySegments() (num int) {return len(model.dirtySegments)}
func (model *Model) DirtySegment(n int) (ds *DirtySegment) {return model.dirtySegments[n]}
func (model *Model) Metadata() (metadata map[string]string) {return model.metadata}
func (model *Model) MetadataValue(key string) (value string, ok bool) {value, ok = model.metadata[key]; return}
func (model *Model) NVertices() (num int) {return len(model.vertices)}
func (model *Model) Vertex(n int) (vertex *Vertex) {return model.vertices[n]}
func (model *Model) NTriangles() (num int) {return len(model.triangles)}
func (model *Model) Triangle(n int) (triangle *Triangle) {return model.triangles[n]}
func (model *Model) NTriangleNormals() (num int) {return len(model.triangleNormals)}
func (model *Model) TriangleNormals(n int) (triangleNormals *TriangleNormals) {return model.triangleNormals[n]}
func (model *Model) NTextureCoordinates() (num int) {return len(model.textureCoordinates)}
func (model *Model) TextureCoordinates(n int) (textureCoordinates *TextureCoordinates) {return model.textureCoordinates[n]}

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
