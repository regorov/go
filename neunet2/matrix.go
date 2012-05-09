package neunet2

import (
    "fmt"
)

type Applicator func(float32) float32

type Matrix struct {
    Data [][]float32
    Rows int
    Cols int
}

func NewMatrix(rows int, cols int) (matrix Matrix) {
    matrix.Data = make([][]float32, rows)
    matrix.Rows = rows
    matrix.Cols = cols

    for i := 0; i < rows; i++ {
        matrix.Data[i] = make([]float32, cols)
    }

    return matrix
}

func FillMatrix(rows int, cols int, value float32) (matrix Matrix) {
    matrix = NewMatrix(rows, cols)
    matrix.IFill(value)
    return matrix
}

func (matrix Matrix) IFill(value float32) {
    for i := 0; i < matrix.Rows; i++ {
        for j := 0; j < matrix.Cols; j++ {
            matrix.Data[i][j] = value
        }
    }
}

func (matrix Matrix) Fill(value float32) (result Matrix) {
    result = matrix.Copy()
    result.IFill(value)
    return result
}

func (matrix Matrix) IApply(f Applicator) {
    for i := 0; i < matrix.Rows; i++ {
        for j := 0; j < matrix.Cols; j++ {
            matrix.Data[i][j] = f(matrix.Data[i][j])
        }
    }
}

func (matrix Matrix) Apply(f Applicator) (result Matrix) {
    result = matrix.Copy()
    result.IApply(f)
    return result
}

func (matrix Matrix) Copy() (result Matrix) {
    result = NewMatrix(matrix.Rows, matrix.Cols)

    for i := 0; i < matrix.Rows; i++ {
        for j := 0; j < matrix.Cols; j++ {
            result.Data[i][j] = matrix.Data[i][j]
        }
    }

    return result
}

func (matrix Matrix) Transpose() (result Matrix) {
    result = NewMatrix(matrix.Cols, matrix.Rows)

    for i := 0; i < matrix.Rows; i++ {
        for j := 0; j < matrix.Cols; j++ {
            result.Data[j][i] = matrix.Data[i][j]
        }
    }

    return result
}

func (matrix Matrix) IAdd(other Matrix) {
    for i := 0; i < matrix.Rows; i++ {
        for j := 0; j < matrix.Cols; j++ {
            matrix.Data[i][j] += other.Data[i][j]
        }
    }
}

func (matrix Matrix) Add(other Matrix) (result Matrix) {
    result = matrix.Copy()
    result.IAdd(other)
    return result
}

func (matrix Matrix) ISub(other Matrix) {
    for i := 0; i < matrix.Rows; i++ {
        for j := 0; j < matrix.Cols; j++ {
            matrix.Data[i][j] -= other.Data[i][j]
        }
    }
}

func (matrix Matrix) Sub(other Matrix) (result Matrix) {
    result = matrix.Copy()
    result.ISub(other)
    return result
}

func (matrix Matrix) ISMul(x float32) {
    for i := 0; i < matrix.Rows; i++ {
        for j := 0; j < matrix.Cols; j++ {
            matrix.Data[i][j] *= x
        }
    }
}

func (matrix Matrix) SMul(x float32) (result Matrix) {
    result = matrix.Copy()
    result.ISMul(x)
    return result
}

func (matrix Matrix) ISDiv(x float32) {
    for i := 0; i < matrix.Rows; i++ {
        for j := 0; j < matrix.Cols; j++ {
            matrix.Data[i][j] /= x
        }
    }
}

func (matrix Matrix) SDiv(x float32) (result Matrix) {
    result = matrix.Copy()
    result.ISDiv(x)
    return result
}

func (matrix Matrix) Mul(other Matrix) (result Matrix) {
    if matrix.Cols != other.Rows {
        matrix, other = other, matrix
    }

    if matrix.Cols != other.Rows {
        panic("The two specified matrices cannot be multiplied")
    }

    rows := matrix.Rows
    cols := other.Cols
    cross := matrix.Cols // = other.Rows

    result = NewMatrix(rows, cols)

    for i := 0; i < rows; i++ {
        for j := 0; j < cols; j++ {
            var sum float32 = 0.0

            for k := 0; k < cross; k++ {
                sum += matrix.Data[i][k] * other.Data[k][j]
            }

            result.Data[i][j] = sum
        }
    }

    return result
}

func (matrix Matrix) IHMul(other Matrix) {
    for i := 0; i < matrix.Rows; i++ {
        for j := 0; j < matrix.Cols; j++ {
            matrix.Data[i][j] *= other.Data[i][j]
        }
    }
}

func (matrix Matrix) HMul(other Matrix) (result Matrix) {
    result = matrix.Copy()
    result.IHMul(other)
    return result
}

func (matrix Matrix) KMul(other Matrix) (result Matrix) {
    result = NewMatrix(matrix.Rows*other.Rows, matrix.Cols*other.Cols)

    r1 := matrix.Rows
    r2 := other.Rows
    c1 := matrix.Cols
    c2 := other.Cols

    r := 0
    for i := 0; i < r1; i++ {
        for j := 0; j < r2; j++ {
            c := 0

            for k := 0; k < c1; k++ {
                for l := 0; l < c2; l++ {
                    //c := (k * c2) + l
                    //r := (i * r2) + j

                    result.Data[r][c] = matrix.Data[i][k] * other.Data[j][l]

                    c++
                }
            }

            r++
        }
    }

    return result
}

func (matrix Matrix) HConcat(other Matrix) (result Matrix) {
    result = NewMatrix(matrix.Rows, matrix.Cols+other.Cols)

    for i := 0; i < matrix.Rows; i++ {
        for j := 0; j < matrix.Cols; j++ {
            result.Data[i][j] = matrix.Data[i][j]
        }

        for j := 0; j < other.Cols; j++ {
            result.Data[i][matrix.Cols+j] = other.Data[i][j]
        }
    }

    return result
}

func (matrix Matrix) Sum() (sum float32) {
    sum = 0

    for i := 0; i < matrix.Rows; i++ {
        for j := 0; j < matrix.Cols; j++ {
            sum += matrix.Data[i][j]
        }
    }

    return sum
}

func (matrix Matrix) Show() {
    for i := 0; i < matrix.Rows; i++ {
        for j := 0; j < matrix.Cols; j++ {
            fmt.Printf("%.3f\t", matrix.Data[i][j])
        }
        fmt.Printf("\n")
    }

    fmt.Printf("\n")
}
