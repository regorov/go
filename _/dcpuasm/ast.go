package main

const (
    O_REG = iota
    O_MEM
    O_MEMDISP
    O_POP
    O_PEEK
    O_PUSH
    O_SP
    O_PC
    O_O
    O_MEMIMM
    O_IMM
)

type Node interface {
    GetName() string
    GetChildNodes() []Node
}

func joinNodeLists(nodeLists ...[]Node) (result []Node) {
    size := 0
    for _, nodeList := range nodeLists {size += len(nodeList)}
    result = make([]Node, size)
    pos := 0
    for _, nodeList := range nodeLists {
        copy(result[pos:], nodeList)
        pos += len(nodeList)
    }
    return result
}

type Assembly struct {
    Objects []Node
}

func NewAssembly(objects []Node) (node *Assembly) {
    node = new(Assembly)
    node.Objects = objects
    return node
}

func (node *Assembly) GetName() (name string) {
    return "Assembly"
}

func (node *Assembly) GetChildNodes() (childNodes []Node) {
    return node.Objects
}

type Constant struct {
    Value int
}

func NewConstant(value int) (node *Constant) {
    node = new(Constant)
    node.Value = value
    return node
}

func (node *Constant) GetName() (name string) {
    return "Constant"
}

func (node *Constant) GetChildNodes() (childNodes []Node) {
    return []Node{}
}

type Instruction struct {
    Name string
    Operands []Node
}

func NewInstruction(name string, operands []Node) (node *Instruction) {
    node = new(Instruction)
    node.Name = name
    node.Operands = operands
    return node
}

func (node *Instruction) GetName() (name string) {
    return "Instruction"
}

func (node *Instruction) GetChildNodes() (childNodes []Node) {
    return node.Operands
}

type Label struct {
    Name string
}

func NewLabel(name string) (node *Label) {
    node = new(Label)
    node.Name = name
    return node
}

func (node *Label) GetName() (name string) {
    return "Label"
}

func (node *Label) GetChildNodes() (childNodes []Node) {
    return []Node{}
}

type LabelRef struct {
    Name string
}

func NewLabelRef(name string) (node *LabelRef) {
    node = new(LabelRef)
    node.Name = name
    return node
}

func (node *LabelRef) GetName() (name string) {
    return "LabelRef"
}

func (node *LabelRef) GetChildNodes() (childNodes []Node) {
    return []Node{}
}

type Operand struct {
    Format int
    X Node
    Y Node
}

func NewOperand(format int, x Node, y Node) (node *Operand) {
    node = new(Operand)
    node.Format = format
    node.X = x
    node.Y = y
    return node
}

func (node *Operand) GetName() (name string) {
    return "Operand"
}

func (node *Operand) GetChildNodes() (childNodes []Node) {
    return []Node{node.X, node.Y}
}

func (node *Operand) GetLength() (length int) {
    switch node.Format {
    case O_MEMDISP:
        return 1
    case O_MEMIMM:
        return 1
    case O_IMM:
        if node.X >= 0x20 {return 1}
    }
    
    return 0
}

type Register struct {
    Number int
}

func NewRegister(number int) (node *Register) {
    node = new(Register)
    node.Number = number
    return node
}

func (node *Register) GetName() (name string) {
    return "Register"
}

func (node *Register) GetChildNodes() (childNodes []Node) {
    return []Node{}
}

