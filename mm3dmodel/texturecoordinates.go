package mm3dmodel

type TextureCoordinates struct {
    index         int
    flags         uint16
    triangleIndex uint32
    triangle      *Triangle
    s1            float32
    s2            float32
    s3            float32
    t1            float32
    t2            float32
    t3            float32
}

// Function Index returns the index of this definiton into its parent model's list.
func (textureCoordinates *TextureCoordinates) Index() (index int) {
    return textureCoordinates.index
}

func (textureCoordinates *TextureCoordinates) Flags() (flags uint16) {
    return textureCoordinates.flags
}

func (textureCoordinates *TextureCoordinates) TriangleIndex() (triangleIndex uint32) {
    return textureCoordinates.triangleIndex
}

func (textureCoordinates *TextureCoordinates) Triangle() (triangle *Triangle) {
    return textureCoordinates.triangle
}

func (textureCoordinates *TextureCoordinates) Vertex1Coord() (s float32, t float32) {
    return textureCoordinates.s1, textureCoordinates.t1
}

func (textureCoordinates *TextureCoordinates) Vertex2Coord() (s float32, t float32) {
    return textureCoordinates.s2, textureCoordinates.t2
}

func (textureCoordinates *TextureCoordinates) Vertex3Coord() (s float32, t float32) {
    return textureCoordinates.s3, textureCoordinates.t3
}

func (textureCoordinates *TextureCoordinates) Associate(model *Model) {
    if model.triangles != nil {
        textureCoordinates.triangle = model.triangles[textureCoordinates.triangleIndex]
        textureCoordinates.triangle.textureCoordinates = textureCoordinates
    }
}
