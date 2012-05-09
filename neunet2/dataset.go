package neunet2

type Dataset struct {
    InputCount    int
    OutputCount   int
    TrainingSet   *Set
    ValidationSet *Set
    TestSet       *Set
}

type Set struct {
    Input  Matrix
    Output Matrix
    Class  Matrix
    Bias   Matrix
    Count  int
}
