package mm3dmodel

// Type Group represents a group of triangles with a specified smoothness and material.
type Group struct {
    index           int
    flags           uint16
    name            string
    triangleIndices []uint32
    triangles       []*Triangle
    smoothness      uint8
    materialIndex   uint32
    smoothnessAngle *SmoothnessAngle
    //material        *Material
}

// Function Index returns the index of this definition into its parent model's list.
func (group *Group) Index() (index int) {
    return group.index
}

// Function Flags returns the flags for this group definition.
func (group *Group) Flags() (flags uint16) {
    return group.flags
}

// Function Name returns the name of this group definition.
func (group *Group) Name() (name string) {
    return group.name
}

// Function NTriangles returns the number of triangles in this group definition.
func (group *Group) NTriangles() (num int) {
    return len(group.triangleIndices)
}

func (group *Group) Associate(model *Model) {

}
