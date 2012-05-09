package neunet2

import (
    "math"
    "math/rand"
)

const PLOT_GRAPHS = true

func OutputToClass(output Matrix) (class Matrix) {
    class = NewMatrix(output.Rows, 1)

    for i := 0; i < output.Rows; i++ {
        var maxOutput float32 = 0.0
        var maxIndex int = 0

        for j := 0; j < output.Cols; j++ {
            o := output.Data[i][j]
            if o > maxOutput {
                maxOutput = o
                maxIndex = j
            }
        }

        class.Data[i][0] = float32(maxIndex + 1)
    }

    return class
}

func ClassToOutput(class Matrix, numClasses int) (output Matrix) {
    output = NewMatrix(class.Rows, numClasses)

    for i := 0; i < class.Rows; i++ {
        c := class.Data[i][0]
        output.Data[i][int(c)] = 1.0
    }

    return output
}

func Activate(x Matrix) (y Matrix) {
    return x.Apply(func(a float32) (b float32) {
        t := float32(math.Tanh(float64(a)))
        return (t + 1.0) / 2.0
    })
}

func DActivate(x Matrix) (y Matrix) {
    return x.Apply(func(a float32) (b float32) {
        t := float32(math.Tanh(float64(a)))
        return (1.0 - (t * t)) / 2.0
    })
}

func FeedForward(input Matrix, weight Matrix, bias Matrix) (net Matrix, output Matrix) {
    net = weight.Mul(input.HConcat(bias))
    output = Activate(net)
    return net, output
}

func InitWeights(maxWeight float32, rows int, cols int) (matrix Matrix) {
    matrix = NewMatrix(rows, cols)

    for i := 0; i < rows; i++ {
        for j := 0; j < cols; j++ {
            matrix.Data[i][j] = (rand.Float32() * 2.0 * maxWeight) - maxWeight
        }
    }

    return matrix
}

func square(x float32) (y float32) {
    return x * x
}

func CalcError(set *Set, weight Matrix) (e float32, c float32) {
    _, output := FeedForward(set.Input, weight, set.Bias)
    sum := set.Output.Sub(output).Apply(square).Sum()
    class := OutputToClass(output)
    e = sum / float32(set.Count*weight.Cols)

    c = 0.0
    for i := 0; i < class.Rows; i++ {
        if class.Data[i][0] == set.Class.Data[i][0] {
            c += 1.0
        }
    }

    c /= float32(set.Count)
    return e, c
}

func BackPropagate(input Matrix, weight Matrix, bias Matrix, eta float32) (weight2 Matrix) {
    return weight
}

func Train(dataset *Dataset) (weight Matrix, e_train, e_validate, e_test, c_train, c_validate, c_test float32) {
    weight = InitWeights(0.5, dataset.InputCount, dataset.OutputCount)

    bias_train := FillMatrix(dataset.TrainingSet.Count, 1, 1.0)
    bias_validate := FillMatrix(dataset.ValidationSet.Count, 1, 1.0)
    bias_test := FillMatrix(dataset.TestSet.Count, 1, 1.0)

    var e_train_l []float32
    var e_validate_l []float32
    var e_test_l []float32
    var c_train_l []float32
    var c_validate_l []float32
    var c_test_l []float32

    if PLOT_GRAPHS {
        e_train_l := make([]float32, 0)
        e_validate_l := make([]float32, 0)
        e_test_l := make([]float32, 0)
        c_train_l := make([]float32, 0)
        c_validate_l := make([]float32, 0)
        c_test_l := make([]float32, 0)
    }

    epoch := 0

    for epoch < 500 {
        weight = BackPropagate(data.TrainingSet.Input, weight, 0.1, bias_train)

        if PLOT_GRAPHS {
            e_train, c_train = CalcError(dataset.TrainingSet, weight)
            e_validate, c_validate = CalcError(dataset.ValidationSet, weight)
            e_test, c_test = CalcError(dataset.TestSet, weight)

            e_train_l = append(e_train_l, e_train)
            e_validate_l = append(e_validate_l, e_validate)
            e_test_l = append(e_test_l, e_test)
            c_train_l = append(c_train_l, c_train)
            c_validate_l = append(c_validate_l, c_validate)
            c_test_l = append(c_test_l, c_test)
        }

        epoch++
    }

    if PLOT_GRAPHS {
        PlotGoogleChart()
    }

    e_train, c_train = CalcError(dataset.TrainingSet, weight)
    e_validate, c_validate = CalcError(dataset.ValidationSet, weight)
    e_test, c_test = CalcError(dataset.TestSet, weight)

    return
}
