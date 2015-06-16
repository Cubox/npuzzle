package main

func abs(x int8) uint16 {
    if x < 0 {
        return uint16(-x)
    }

    return uint16(x)
}

func linearConflictCell(value uint16, x, y uint8, board *Board, goal []Cell) uint16 {
    distance := manhattanCell(value, x, y, board, goal)

    for _, n := range board.cells[x][y+1:] {
        if n != 0 && goal[n].x == goal[value].x && goal[n].y < goal[value].y {
            distance += 2
        }
    }

    for i := x + 1; i < size; i++ {
        n := board.cells[i][y]
        if n != 0 && goal[n].y == goal[value].y && goal[n].x < goal[value].x {
            distance += 2
        }
    }

    return distance
}

func manhattanCell(value uint16, x, y uint8, board *Board, goal []Cell) uint16 {
    cell := goal[value]
    var distance uint16
    distance += abs(int8(cell.x - x))
    distance += abs(int8(cell.y - y))

    return distance
}

func hammingCell(value uint16, x, y uint8, board *Board, goal []Cell) uint16 {
    cell := goal[value]
    if cell.x == x && cell.y == y {
        return 0
    } else {
        return 1
    }
}

func (b *Board) priority(goal []Cell) uint16 {
    var distance uint16

    for i, n := range b.cells {
        for j, m := range n {
            if m != 0 {
                distance += method(m, uint8(i), uint8(j), b, goal) // method is a pointer to one of the functions up here
            }
        }
    }

    return distance
}
