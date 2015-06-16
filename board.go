package main

import (
    "bufio"
    "fmt"
    "log"
    "math/rand"
    "os"
    "sort"
    "strconv"
    "sync"
    "time"
)

var (
    size   uint8
    method func(value uint16, x, y uint8, board *Board, goal []Cell) uint16
    pool   = sync.Pool{New: allocateBoard}
)

type Board struct {
    cells      [][]uint16
    parent     *Board
    f, g uint16
    x, y uint8
}

type Cell struct {
    x, y uint8
}

func allocateBoard() interface{} {
    board := make([][]uint16, size)
    for i := uint8(0); i < size; i++ {
        board[i] = make([]uint16, size)
    }

    return board
}

func isGoal(b *Board, goal []Cell) bool {
    for x := range b.cells {
        for y, m := range b.cells[x] {
            if goal[m].x != uint8(x) || goal[m].y != uint8(y) {
                return false
            }
        }
    }

    return true
}

func initGoal(goal *[]Cell) {
    var x, y int8
    var vy int8
    value := uint16(1)
    vx := int8(1)

    *goal = make([]Cell, size*size)
    board := make([][]bool, size)
    for i := uint8(0); i < size; i++ {
        board[i] = make([]bool, size)
    }

    for value < uint16(size*size) {
        board[x][y] = true
        (*goal)[value] = Cell{x: uint8(y), y: uint8(x)}

        value++
        if uint8(x+vx) == size || x+vx < 0 || (vx != 0 && board[x+vx][y]) {
            vy, vx = vx, 0
        }
        if uint8(y+vy) == size || y+vy < 0 || (vy != 0 && board[x][y+vy]) {
            vx, vy = -vy, 0
        }
        x += vx
        y += vy
    }

    (*goal)[0] = Cell{x: uint8(y), y: uint8(x)}
}

func (b *Board) copy() *Board {
    board := *b
    board.cells = pool.Get().([][]uint16)
    for i := uint8(0); i < size; i++ {
        copy(board.cells[i], b.cells[i])
    }

    return &board
}

func (b *Board) neighbours() []*Board {
    neighbours := make([]*Board, 0)

    if b.y > 0 {
        n := b.copy()
        n.cells[b.x][b.y], n.cells[b.x][b.y-1] = n.cells[b.x][b.y-1], n.cells[b.x][b.y]
        n.y--
        neighbours = append(neighbours, n)
    }

    if b.x < size-1 {
        n := b.copy()
        n.cells[b.x][b.y], n.cells[b.x+1][b.y] = n.cells[b.x+1][b.y], n.cells[b.x][b.y]
        n.x++
        neighbours = append(neighbours, n)
    }

    if b.y < size-1 {
        n := b.copy()
        n.cells[b.x][b.y], n.cells[b.x][b.y+1] = n.cells[b.x][b.y+1], n.cells[b.x][b.y]
        n.y++
        neighbours = append(neighbours, n)
    }

    if b.x > 0 {
        n := b.copy()
        n.cells[b.x][b.y], n.cells[b.x-1][b.y] = n.cells[b.x-1][b.y], n.cells[b.x][b.y]
        n.x--
        neighbours = append(neighbours, n)
    }

    return neighbours
}

func (b *Board) check() {
    cells := make([]int, size*size)
    for i := range b.cells {
        for j := range b.cells[i] {
            cells[uint8(i)*size+uint8(j)] = int(b.cells[i][j])
        }
    }

    sort.Ints(cells)

    for i := range cells {
        if cells[i] != i {
            log.Fatalln("Board invalid!")
        }
    }
}

func parseToken(data []byte, atEOF bool) (advance int, token []byte, err error) {
    if atEOF && len(data) == 0 {
        return 0, nil, nil
    }

    i := 0

    for i < len(data) && (data[i] == ' ' || data[i] == '\n') {
        i++
    }

    if i >= len(data) {
        return i, make([]byte, 0), nil
    }

    if data[i] == '#' { // We ignore until next \n
        advance, _, err = bufio.ScanLines(data[i:], atEOF)
        return advance + i, make([]byte, 0), err
    } else {
        advance, buf, err := bufio.ScanWords(data[i:], atEOF)
        return advance + i, buf, err
    }
}

func (b *Board) parseBoard(fileName string) {
    file, err := os.Open(fileName)
    if err != nil {
        log.Fatalln(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Split(parseToken)

    for scanner.Scan() {
        pSize, err := strconv.Atoi(scanner.Text())
        size = uint8(pSize)

        if err == nil {
            break
        }
    }

    if size < 3 || size > 20 {
        log.Fatalln("Invalid size")
    }

    b.cells = pool.Get().([][]uint16)

    for i := uint8(0); scanner.Scan(); {
        text := scanner.Text()

        if len(text) > 0 {
            value, err := strconv.Atoi(text)
            if err != nil {
                log.Fatalln(err)
            }

            b.cells[i/size][i%size] = uint16(value)

            if b.cells[i/size][i%size] == 0 {
                b.x = i / size
                b.y = i % size
            }

            i++
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatalln(err)
    }

    b.check()
}

func (b *Board) generateBoard() {
    if genSize < 3 || genSize > 20 {
        log.Fatalln("Invalid size")
    }

    size = uint8(genSize)

    b.cells = pool.Get().([][]uint16)

    rand.Seed(time.Now().UnixNano())
    numbers := rand.Perm(int(size * size))

    for n, i := range numbers {
        b.cells[uint8(i)/size][uint8(i)%size] = uint16(n)
        if n == 0 {
            b.x = uint8(i) / size
            b.y = uint8(i) % size
        }
    }

    b.check()
}

func (b *Board) String() string {
    var str string
    str += fmt.Sprintln("Board size:", size)
    for _, n := range b.cells {
        for _, m := range n {
            str += fmt.Sprintf("%6d", m)
        }
        str += fmt.Sprintln()
    }

    return str[:len(str)-1] // Remove extra \n
}
