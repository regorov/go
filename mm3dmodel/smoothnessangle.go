package mm3dmodel

// Type SmoothnessAngle represents the normals for the three vertices of a triangle.
type SmoothnessAngle struct {
    index      int
    groupIndex uint32
    group      *Group
    angle      uint8
}

// Function Index returns the index of this definition into its parent model's list.
func (smoothnessAngle *SmoothnessAngle) Index() (index int) {
    return smoothnessAngle.index
}

// Function TriangleIndex returns the index into the associated model's triangles of the triangle
// that this structure describes.
func (smoothnessAngle *SmoothnessAngle) GroupIndex() (groupIndex uint32) {
    return smoothnessAngle.groupIndex
}

// Function Triangle returns the group that this structure describes.
func (smoothnessAngle *SmoothnessAngle) Group() (group *Group) {
    return smoothnessAngle.group
}

// Function Angle returns the angle of the structure.
func (smoothnessAngle *SmoothnessAngle) Angle() (angle uint8) {
    return smoothnessAngle.angle
}

// Function Associate updates the links between this structure and others, using the model it was
// parsed from as a reference point.
func (smoothnessAngle *SmoothnessAngle) Associate(model *Model) {
    if model.groups != nil {
        smoothnessAngle.group = model.groups[smoothnessAngle.groupIndex]
        smoothnessAngle.group.smoothnessAngle = smoothnessAngle
    }
}
