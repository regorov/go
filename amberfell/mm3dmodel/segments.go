package mm3dmodel

import (
    "bytes"
    "encoding/binary"
)

func ParseMetadataSegment(model *Model, dataFlags uint16, dataElements [][]byte) {
    /*
    for _, element := range dataElements {
        
    }
    */
}

type parsedVertex struct {
    Flags uint16
    X float32
    Y float32
    Z float32
}

func ParseVerticesSegment(model *Model, dataFlags uint16, dataElements [][]byte) {
    model.vertices = make([]*Vertex, len(dataElements))
    
    for i, element := range dataElements {
        p := &parsedVertex{}
        binary.Read(bytes.NewReader(element), binary.LittleEndian, p)
        
        vertex := &Vertex{
            flags: p.Flags,
            x: p.X,
            y: p.Y,
            z: p.Z,
        }
        
        model.vertices[i] = vertex
    }
}

type parsedTriangle struct {
    Flags uint16
    V1Index uint32
    V2Index uint32
    V3Index uint32
}

func ParseTrianglesSegment(model *Model, dataFlags uint16, dataElements [][]byte) {
    model.triangles = make([]*Triangle, len(dataElements))
    
    for i, element := range dataElements {
        p := &parsedTriangle{}
        binary.Read(bytes.NewReader(element), binary.LittleEndian, p)
        
        triangle := &Triangle{
            flags: p.Flags,
            v1index: p.V1Index,
            v2index: p.V2Index,
            v3index: p.V3Index,
        }
        
        triangle.Associate(model)
        model.triangles[i] = triangle
    }
}

type parsedTriangleNormals struct {
    Flags uint16
    Index uint32
    V1Normal [3]float32
    V2Normal [3]float32
    V3Normal [3]float32
}

func ParseTriangleNormalsSegment(model *Model, dataFlags uint16, dataElements [][]byte) {
    model.triangleNormals = make([]*TriangleNormals, len(dataElements))
    
    for i, element := range dataElements {
        p := &parsedTriangleNormals{}
        binary.Read(bytes.NewReader(element), binary.LittleEndian, p)
        
        triangleNormals := &TriangleNormals{
            flags: p.Flags,
            triangleIndex: p.Index,
            v1x: p.V1Normal[0],
            v1y: p.V1Normal[1],
            v1z: p.V1Normal[2],
            v2x: p.V2Normal[0],
            v2y: p.V2Normal[1],
            v2z: p.V2Normal[2],
            v3x: p.V3Normal[0],
            v3y: p.V3Normal[1],
            v3z: p.V3Normal[2],
        }
        
        triangleNormals.Associate(model)
        model.triangleNormals[i] = triangleNormals
    }
}

type parsedTextureCoordinates struct {
    Flags uint16
    Index uint32
    S [3]float32
    T [3]float32
}

func ParseTextureCoordinatesSegment(model *Model, dataFlags uint16, dataElements [][]byte) {
    model.textureCoordinates = make([]*TextureCoordinates, len(dataElements))
    
    for i, element := range dataElements {
        p := &parsedTextureCoordinates{}
        binary.Read(bytes.NewReader(element), binary.LittleEndian, p)
        
        textureCoordinates := &TextureCoordinates{
            flags: p.Flags,
            triangleIndex: p.Index,
            s1: p.S[0],
            s2: p.S[1],
            s3: p.S[2],
            t1: p.T[0],
            t2: p.T[1],
            t3: p.T[2],
        }
        
        textureCoordinates.Associate(model)
        model.textureCoordinates[i] = textureCoordinates
    }
}
