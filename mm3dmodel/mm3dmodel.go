// Package mm3dmodel provides MM3D model file loading and rendering for Amberfell.
package mm3dmodel

import (
    "encoding/binary"
    "errors"
    "io"
)

var (
    ErrMagicMismatch        = errors.New("magic number mismatch")
    ErrUnimplementedVersion = errors.New("unimplemented version; can only handle major version 1")
)

// Variable MAGIC is the expected magic number in MM3D files.
var MAGIC = []byte("MISFIT3D")

// Type parsedHeader is used to parse a the raw header from the MM3D file.
type parsedHeader struct {
    Magic        [8]uint8
    MajorVersion uint8
    MinorVersion uint8
    ModelFlags   uint8
    OffsetCount  uint8
}

// Type parsedOffset is used to parse a raw segment offset from the MM3D file.
type parsedOffset struct {
    Type  uint16
    Value uint32
}

// Type parsedSegmentHeader is used to parse a raw segment header from the MM3D file.
type parsedSegmentHeader struct {
    Flags uint16
    Count uint32
}

// Type parsedSegmentSize is used to parse a raw segment size from the MM3D file, whether it be
// attached to the header or a data element.
type parsedSegmentSize struct {
    Size uint32
}

// Function Read reads an MM3D model definition from the specified io.Reader and returns either the
// parsed model and nil, or nil and an error.
func Read(reader io.ReadSeeker) (model *Model, err error) {
    var header *parsedHeader
    var offset *parsedOffset
    var segheader *parsedSegmentHeader
    var segsize *parsedSegmentSize

    err = binary.Read(reader, binary.LittleEndian, header)
    if err != nil {
        return nil, err
    }

    for i := 0; i < 8; i++ {
        if header.Magic[i] != MAGIC[i] {
            return nil, ErrMagicMismatch
        }
    }

    if header.MajorVersion != 0x01 {
        return nil, ErrUnimplementedVersion
    }

    model = &Model{
        majorVersion:  header.MajorVersion,
        minorVersion:  header.MinorVersion,
        modelFlags:    header.ModelFlags,
        dirtySegments: make([]*DirtySegment, 0),
    }

    offsets := make([]*parsedOffset, 0, header.OffsetCount)

    for i := uint8(0); i < header.OffsetCount; i++ {
        offset = new(parsedOffset)
        err = binary.Read(reader, binary.LittleEndian, offset)
        if err != nil {
            return nil, err
        }

        if offset.Type == 0x3FFF {
            continue
        }

        offsets = append(offsets, offset)
    }

    for _, offset := range offsets {
        offsetType := offset.Type
        offsetValue := offset.Value

        _, err = reader.Seek(int64(offsetValue), 0)
        if err != nil {
            return nil, err
        }

        err = binary.Read(reader, binary.LittleEndian, segheader)
        if err != nil {
            return nil, err
        }

        dataFlags := segheader.Flags
        dataCount := segheader.Count

        dataElements := make([][]byte, dataCount)

        if (offsetType & 0x8000) != 0 {
            err = binary.Read(reader, binary.LittleEndian, segsize)
            if err != nil {
                return nil, err
            }

            for j := uint32(0); j < dataCount; j++ {
                element := make([]byte, segsize.Size)
                _, err = reader.Read(element)
                if err != nil {
                    return nil, err
                }

                dataElements[j] = element
            }

        } else {
            for j := uint32(0); j < dataCount; j++ {
                err = binary.Read(reader, binary.LittleEndian, segsize)
                if err != nil {
                    return nil, err
                }

                element := make([]byte, segsize.Size)
                _, err = reader.Read(element)
                if err != nil {
                    return nil, err
                }

                dataElements[j] = element
            }
        }

        err = nil

        switch offsetType & 0xBFFF { // Exclude bit 14, the dirty bit
        case 0x1001:
            err = ParseMetadataSegment(model, dataFlags, dataElements)
        case 0x0101:
            err = ParseGroupsSegment(model, dataFlags, dataElements)
        //case 0x0142: err = ParseExternalTexturesSegment            (model, dataFlags, dataElements)
        //case 0x0161: err = ParseMaterialsSegment                   (model, dataFlags, dataElements)
        //case 0x016c: err = ParseTextureProjectionsTrianglesSegment (model, dataFlags, dataElements)
        //case 0x0191: err = ParseCanvasBackgroundImagesSegment      (model, dataFlags, dataElements)
        //case 0x0301: err = ParseSkeletalAnimationsSegment          (model, dataFlags, dataElements)
        //case 0x0321: err = ParseFrameAnimationsSegment             (model, dataFlags, dataElements)
        //case 0x0326: err = ParseFrameAnimationPointsSegment        (model, dataFlags, dataElements)
        case 0x8001:
            err = ParseVerticesSegment(model, dataFlags, dataElements)
        case 0x8021:
            err = ParseTrianglesSegment(model, dataFlags, dataElements)
        case 0x8026:
            err = ParseTriangleNormalsSegment(model, dataFlags, dataElements)
        //case 0x8041: err = ParseJointsSegment                      (model, dataFlags, dataElements)
        //case 0x8046: err = ParseJointVerticesSegment               (model, dataFlags, dataElements)
        //case 0x8061: err = ParsePointsSegment                      (model, dataFlags, dataElements)
        //case 0x8106: err = ParseSmoothnessAnglesSegment            (model, dataFlags, dataElements)
        //case 0x8146: err = ParseWeightedInfluencesSegment          (model, dataFlags, dataElements)
        //case 0x8168: err = ParseTextureProjectionsSegment          (model, dataFlags, dataElements)
        case 0x8121:
            err = ParseTextureCoordinatesSegment(model, dataFlags, dataElements)

        default:
            model.dirtySegments = append(model.dirtySegments, &DirtySegment{
                index:      len(model.dirtySegments),
                offsetType: offsetType,
                flags:      dataFlags,
                elements:   dataElements,
            })
        }

        if err != nil {
            return nil, err
        }
    }

    return model, nil
}
