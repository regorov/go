// Package hydraulics provides a backend for Amberfell's hydraulics system.
package hydraulics

import (
    
)

// Function balance returns the two arguments averaged. With a normal integral mean of say 7 and 10,
// the result is 17 / 2 = 8. However 8 + 8 (16) != 7 + 10 (17) - one unit has been lost. This
// function ensures that even if the total is not divisible by two, the two results still add to the
// same sum as the arguments: balance(7, 10) = (8, 9), and 8 + 9 == 7 + 10.
func balance(a int, b int) (x int, y int) {
    tot := a + b
    avg := tot / 2
    return tot - avg, avg
}

// Function balanceYLimited works exactly like balance, except that it places a maximum limit on the
// second return value.
func balanceYLimited(a int, b int, yLimit int) (x int, y int) {
    tot := a + b
    avg := tot / 2
    if tot >= yLimit {
        avg = yLimit
    }
    
    return tot - avg, avg
}

// Function balance3ZLimited works exactly like balanceYLimited, except that it processes 3
// arguments instead of two.
func balance3ZLimited(a int, b int, c int, zLimit int) (x int, y int, z int) {
    tot := a + b + c
    avg := tot / 3
    rem := tot % 3
    
    if avg <= zLimit {
        if rem == 0 {
            return avg, avg, avg
        } else if rem == 1 {
            return avg + 1, avg, avg
        } else {
            return avg, avg + 1, avg + 1
        }
    }
    
    x = zLimit
    y = zLimit
    z = zLimit
    tot -= zLimit
    
    x_extra, y_extra := balance(tot, 0)
    return x + x_extra, y + y_extra, z
}

// Function balane3YZLimited works exactly like balance3ZLimited, except two of the results have
// limits rather than just one.
func balance3YZLimited(a int, b int, c int, yLimit int, zLimit int) (x int, y int, z int) {
    tot := a + b + c
    avg := tot / 3
    rem := tot % 3
    
    if avg <= yLimit && avg <= zLimit {
        if rem == 0 {
            return avg, avg, avg
        } else if rem == 1 {
            return avg + 1, avg, avg
        } else {
            return avg, avg + 1, avg + 1
        }
    }
    
    if yLimit < zLimit {
        x = yLimit
        y = yLimit
        z = yLimit
        tot -= yLimit
        
        x_extra, z_extra := balanceYLimited(tot, 0, zLimit)
        return x + x_extra, y, z + z_extra
    }
    
    x = zLimit
    y = zLimit
    z = zLimit
    tot -= zLimit
    
    x_extra, y_extra := balanceYLimited(tot, 0, yLimit)
    return x + x_extra, y + y_extra, z
}

// Function imin is just like math.Min, except that it works with ints.
func imin(a int, b int) (x int) {
    if a < b {
        return a
    } else {
        return b
    }
    
    return 0
}

// Interface Component represents any simulatable hydraulic component.
type Component interface {
    Flow()
}

// Interface Producer represents any simulatable hydraulic component that can provide fluid to
// others.
type Producer interface {
    Component
    
    GetOutput() Receiver
    SetOutput(Receiver)
}

// Interface Receiver represents any simulatable hydraulic component that can accept fluid from
// others.
type Receiver interface {
    Component
    
    GetQuantity() (int)
    SetQuantity(int)
    AddQuantity(int)
    
    GetCapacity() (int)
    SetCapacity(int)
}
