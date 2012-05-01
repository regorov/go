package mm3dmodel

import (
    "bufio"
    "bytes"
    "encoding/binary"
)

//import "fmt"

// Function ParseMetadataSegment parses the supplied metadata segment data and adds it to the model.
func ParseMetadataSegment(model *Model, dataFlags uint16, dataElements [][]byte) (err error) {
    if model.metadata == nil {
        model.metadata = make(map[string]string)
    }

    for _, element := range dataElements {
        reader := bufio.NewReader(bytes.NewReader(element))
        key, err := reader.ReadString(0)
        if err != nil {
            return err
        }

        value, err := reader.ReadString(0)
        if err != nil {
            return err
        }

        model.metadata[key[:-1]] = value[:-1]
    }

    return nil
}

// Type parsedGroup16 is used to parse part of a raw group definition from the MM3D file.
type parsedGroup16 struct {
    N uint16
}

// Type parsedGroup32 is used to parse part of a raw group definition from the MM3D file.
type parsedGroup32 struct {
    N uint32
}

// Type parsedGroupFooter is used to parse part of a raw group definition from the MM3D file.
type parsedGroupFooter struct {
    Smoothness uint8
    Material   uint32
}

// Function ParseGroupsSegment parses the supplied group segment data and adds it to the model.
func ParseGroupsSegment(model *Model, dataFlags uint16, dataElements [][]byte) (err error) {
    p16 := new(parsedGroup16)
    p32 := new(parsedGroup32)
    footer := new(parsedGroupFooter)

    if model.groups == nil {
        model.groups = make([]*Group, 0, len(dataElements))
    }

    for i, element := range dataElements {
        reader := bufio.NewReader(bytes.NewReader(element))

        err = binary.Read(reader, binary.LittleEndian, p16)
        if err != nil {
            return err
        }
        flags := p16.N

        name, err := reader.ReadString(0)
        if err != nil {
            return err
        }
        name = name[:len(name)-1]

        err = binary.Read(reader, binary.LittleEndian, p32)
        if err != nil {
            return err
        }
        triangleCount := p32.N

        triangleIndices := make([]uint32, triangleCount)
        for i := uint32(0); i < triangleCount; i++ {
            err = binary.Read(reader, binary.LittleEndian, p32)
            if err != nil {
                return err
            }
            triangleIndices[i] = p32.N
        }

        err = binary.Read(reader, binary.LittleEndian, footer)
        if err != nil {
            return err
        }

        group := &Group{
            index:           i,
            flags:           flags,
            name:            name,
            triangleIndices: triangleIndices,
            smoothness:      footer.Smoothness,
            materialIndex:   footer.Material,
        }

        group.Associate(model)
        model.groups = append(model.groups, group)
    }

    return nil
}

// Type parsedVertex is used to parse a raw vertex definition from the MM3D file.
type parsedVertex struct {
    Flags uint16
    X     float32
    Y     float32
    Z     float32
}

// Function ParseVerticesSegment parses the supplied vertices segment data and adds it to the model.
func ParseVerticesSegment(model *Model, dataFlags uint16, dataElements [][]byte) (err error) {
    p := new(parsedVertex)

    if model.vertices == nil {
        model.vertices = make([]*Vertex, 0, len(dataElements))
    }

    for i, element := range dataElements {
        err = binary.Read(bytes.NewReader(element), binary.LittleEndian, p)
        if err != nil {
            return err
        }

        vertex := &Vertex{
            index: i,
            flags: p.Flags,
            x:     p.X,
            y:     p.Y,
            z:     p.Z,
        }

        model.vertices = append(model.vertices, vertex)
    }

    return nil
}

// Type parsedTriangle is used to parse a raw triangle definition from the MM3D file.
type parsedTriangle struct {
    Flags   uint16
    V1Index uint32
    V2Index uint32
    V3Index uint32
}

// Function ParseTrianglesSegment parses the supplied triangles segment data and adds it to the
// model.
func ParseTrianglesSegment(model *Model, dataFlags uint16, dataElements [][]byte) (err error) {
    p := new(parsedTriangle)

    if model.triangles == nil {
        model.triangles = make([]*Triangle, 0, len(dataElements))
    }

    for i, element := range dataElements {
        err = binary.Read(bytes.NewReader(element), binary.LittleEndian, p)
        if err != nil {
            return err
        }

        triangle := &Triangle{
            index:   i,
            flags:   p.Flags,
            v1index: p.V1Index,
            v2index: p.V2Index,
            v3index: p.V3Index,
        }

        triangle.Associate(model)
        model.triangles = append(model.triangles, triangle)
    }

    return nil
}

// Type parsedTriangleNormals is used to parse a raw triangle normals definition from the MM3D file.
type parsedTriangleNormals struct {
    Flags    uint16
    Index    uint32
    V1Normal [3]float32
    V2Normal [3]float32
    V3Normal [3]float32
}

// Function ParseTriangleNormalsSegment parses the supplied triangle normals segment data and adds
// it to the model.
func ParseTriangleNormalsSegment(model *Model, dataFlags uint16, dataElements [][]byte) (err error) {
    p := new(parsedTriangleNormals)

    if model.triangleNormals == nil {
        model.triangleNormals = make([]*TriangleNormals, 0, len(dataElements))
    }

    for i, element := range dataElements {
        err = binary.Read(bytes.NewReader(element), binary.LittleEndian, p)
        if err != nil {
            return err
        }

        triangleNormals := &TriangleNormals{
            index:         i,
            flags:         p.Flags,
            triangleIndex: p.Index,
            v1x:           p.V1Normal[0],
            v1y:           p.V1Normal[1],
            v1z:           p.V1Normal[2],
            v2x:           p.V2Normal[0],
            v2y:           p.V2Normal[1],
            v2z:           p.V2Normal[2],
            v3x:           p.V3Normal[0],
            v3y:           p.V3Normal[1],
            v3z:           p.V3Normal[2],
        }

        triangleNormals.Associate(model)
        model.triangleNormals = append(model.triangleNormals, triangleNormals)
    }

    return nil
}

// Type parsedTextureCoordinates is used to parse a raw texture coordinates definition from the MM3D
// file.
type parsedTextureCoordinates struct {
    Flags uint16
    Index uint32
    S     [3]float32
    T     [3]float32
}

// Function ParseTextureCoordinatesSegment parses the supplied texture coordinates segment data and
// adds it to the model.
func ParseTextureCoordinatesSegment(model *Model, dataFlags uint16, dataElements [][]byte) (err error) {
    p := new(parsedTextureCoordinates)

    if model.textureCoordinates == nil {
        model.textureCoordinates = make([]*TextureCoordinates, 0, len(dataElements))
    }

    for i, element := range dataElements {
        err = binary.Read(bytes.NewReader(element), binary.LittleEndian, p)
        if err != nil {
            return err
        }

        textureCoordinates := &TextureCoordinates{
            index:         i,
            flags:         p.Flags,
            triangleIndex: p.Index,
            s1:            p.S[0],
            s2:            p.S[1],
            s3:            p.S[2],
            t1:            p.T[0],
            t2:            p.T[1],
            t3:            p.T[2],
        }

        textureCoordinates.Associate(model)
        model.textureCoordinates = append(model.textureCoordinates, textureCoordinates)
    }

    return nil
}

// Type parsedSmoothnessAngles is used to parse a raw smoothness angles definition from the MM3D
// file.
type parsedSmoothnessAngles struct {
    GroupIndex uint32
    Angle      uint8
}

// Function ParseSmoothnessAnglesSegment parses the supplied smoothness angles segment data and adds
// it to the model.
func ParseSmoothnessAnglesSegment(model *Model, dataFlags uint16, dataElements [][]byte) (err error) {
    p := new(parsedSmoothnessAngles)

    if model.smoothnessAngles == nil {
        model.smoothnessAngles = make([]*SmoothnessAngle, 0, len(dataElements))
    }

    for i, element := range dataElements {
        err = binary.Read(bytes.NewReader(element), binary.LittleEndian, p)
        if err != nil {
            return err
        }

        smoothnessAngle := &SmoothnessAngle{
            index:      i,
            groupIndex: p.GroupIndex,
            angle:      p.Angle,
        }

        smoothnessAngle.Associate(model)
        model.smoothnessAngles = append(model.smoothnessAngles, smoothnessAngle)
    }

    return nil
}
