package mm3dmodel

import (
    "github.com/kierdavis/go/bytereader"
    "io"
)

var MAGIC = []byte("MISFIT3D")

func Read(reader io.ReadSeeker) (model *Model, err error) {
    br := bytereader.NewByteReader(reader, bytereader.LITTLE_ENDIAN)
    
    magic, err := br.ReadBytes(8); if err != nil {return nil, err}
    for i := 0; i < 8; i++ {
        if magic[i] != MAGIC[i] {return nil, NewError(E_MAGIC_MISMATCH, "Read: magic number mismatch")}
    }
    
    majorVersion, err := br.ReadUnsigned8(); if err != nil {return nil, err}
    if majorVersion != 0x01 {return nil, NewError(E_UNIMPLEMENTED_VERSION, "Read: this library can only handle major version 1 files.")}
    
    minorVersion, err := br.ReadUnsigned8(); if err != nil {return nil, err}
    modelFlags, err := br.ReadUnsigned8(); if err != nil {return nil, err}
    offsetCount, err := br.ReadUnsigned8(); if err != nil {return nil, err}
    
    model = &Model{
        majorVersion: majorVersion,
        minorVersion: minorVersion,
        modelFlags: modelFlags,
        dirtySegments: make([]*DirtySegment, 0),
    }
    
    for i := uint8(0); i < offsetCount; i++ {
        offsetType, err := br.ReadUnsigned16(); if err != nil {return nil, err}
        offsetValue, err := br.ReadUnsigned32(); if err != nil {return nil, err}
        
        if offsetType == 0x3FFF {
            continue
        }
        
        br.PushPos()
        err = br.JumpTo(int64(offsetValue)); if err != nil {return nil, err}
        
        dataFlags, err := br.ReadUnsigned16(); if err != nil {return nil, err}
        dataCount, err := br.ReadUnsigned32(); if err != nil {return nil, err}
        
        dataElements := make([][]byte, dataCount)
        
        if offsetType & 0x8000 != 0 {
            dataSize, err := br.ReadUnsigned32(); if err != nil {return nil, err}
            
            for j := uint32(0); j < dataCount; j++ {
                element, err := br.ReadBytes(int(dataSize)); if err != nil {return nil, err}
                dataElements[j] = element
            }
        
        } else {
            for j := uint32(0); j < dataCount; j++ {
                dataSize, err := br.ReadUnsigned32(); if err != nil {return nil, err}
                element, err := br.ReadBytes(int(dataSize)); if err != nil {return nil, err}
                dataElements[j] = element
            }
        }
        
        switch (offsetType & 0xBFFF) { // Exclude bit 14, the dirty bit
        case 0x1001: ParseMetadataSegment                    (model, dataFlags, dataElements)
        //case 0x0101: ParseGroupsSegment                      (model, dataFlags, dataElements)
        //case 0x0142: ParseExternalTexturesSegment            (model, dataFlags, dataElements)
        //case 0x0161: ParseMaterialsSegment                   (model, dataFlags, dataElements)
        //case 0x016c: ParseTextureProjectionsTrianglesSegment (model, dataFlags, dataElements)
        //case 0x0191: ParseCanvasBackgroundImagesSegment      (model, dataFlags, dataElements)
        //case 0x0301: ParseSkeletalAnimationsSegment          (model, dataFlags, dataElements)
        //case 0x0321: ParseFrameAnimationsSegment             (model, dataFlags, dataElements)
        //case 0x0326: ParseFrameAnimationPointsSegment        (model, dataFlags, dataElements)
        case 0x8001: ParseVerticesSegment                    (model, dataFlags, dataElements)
        case 0x8021: ParseTrianglesSegment                   (model, dataFlags, dataElements)
        case 0x8026: ParseTriangleNormalsSegment             (model, dataFlags, dataElements)
        //case 0x8041: ParseJointsSegment                      (model, dataFlags, dataElements)
        //case 0x8046: ParseJointVerticesSegment               (model, dataFlags, dataElements)
        //case 0x8061: ParsePointsSegment                      (model, dataFlags, dataElements)
        //case 0x8106: ParseSmoothnessAnglesSegment            (model, dataFlags, dataElements)
        //case 0x8146: ParseWeightedInfluencesSegment          (model, dataFlags, dataElements)
        //case 0x8168: ParseTextureProjectionsSegment          (model, dataFlags, dataElements)
        case 0x8121: ParseTextureCoordinatesSegment          (model, dataFlags, dataElements)
        
        default:
            model.dirtySegments = append(model.dirtySegments, &DirtySegment{offsetType: offsetType, flags: dataFlags, elements: dataElements})
        }
        
        err = br.PopPos(); if err != nil {return nil, err}
    }
    
    return model, nil
}
