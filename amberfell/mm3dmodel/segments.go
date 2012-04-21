package mm3dmodel

import (
    "bytes"
    "github.com/kierdavis/go/bytereader"
)

func ParseMetadataSegment(model *Model, dataFlags uint16, dataElements [][]byte) {
    /*
    for _, element := range dataElements {
        
    }
    */
}

func ParseVerticesSegment(model *Model, dataFlags uint16, dataElements [][]byte) {
    model.Vertices = make([]*Vertex, len(dataElements))
    
    for i, element := range dataElements {
        br := bytereader.NewByteReader(bytes.NewReader(element), bytereader.LITTLE_ENDIAN)
        model.Vertices[i] = &Vertex{
            flags: br.Unsigned16(),
            x: br.Float32(),
            y: br.Float32(),
            z: br.Float32(),
        }
    }
}

func ParseTrianglesSegment(model *Model, dataFlags uint16, dataElements [][]byte) {
    model.Triangles = make([]*Vertex, len(dataElements))
    
    for i, element := range dataElements {
        triangle := &Triangle{}
        binary.Read(bytes.NewReader(element), binary.LittleEndian, triangle)
        model.Triangles[i] = triangle
    }
}