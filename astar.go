package main

import (
    "bytes"
    "fmt"
    "hash/crc32"
    "sync/atomic"
    "time"
    "runtime"
)

type Astar struct {
    open           Queue
    closed         Map
    goal   []Cell
    opened, states uint32 // Shared variables, need atomic access
}

func sum(b *Board) uint32 { // Return hash of the board
    buf := new(bytes.Buffer)
    buf.Grow(int(size*size)*2)

    for i := range b.cells {
        for j := range b.cells[i] {
            buf.WriteByte(byte(b.cells[i][j]))
            buf.WriteByte(byte(b.cells[i][j] >> 8))
        }
    }

    return crc32.ChecksumIEEE(buf.Bytes())
}

func (a *Astar) show(dur time.Duration, cpuU, memStats int) {
    if !dumb {
        fmt.Printf("\033[1A") // Move cursor one line up
    }

    print("Time spent: %s Open nodes: %d Closed nodes: %d Total opened nodes: %d Maximum states ever in memory: %d\n",
                dur, a.open.len(), a.closed.len(), a.opened, a.states)

    print("CPU usage: %d%% Memory usage: %dMo Routines: %d",
                cpuU, memStats, runtime.NumGoroutine())
    if dumb {
        fmt.Println()
    }
}

func (a *Astar) solve(initial *Board) (steps uint16) {
    fmt.Println("Starting path search.\n")

    quit := make(chan *Board) // Used to send search result
    c := make(chan int)       // Used to send cpu/mem stats

    a.open = NewQueue(initial)
    a.closed = NewMap()

    initGoal(&a.goal)

    if isGoal(initial, a.goal) {
        fmt.Println("Success!")
        fmt.Printf("Numbers of steps: %d\n", output(initial))
        fmt.Println("Just for your information, you gave me the goal board.")
        return
    }

    for i := 0; i < threads; i++ {
        go a.astar(quit)
    }

    go cpu(c)
    startTime := time.Now()

    cpuU := 0
    memStats := 0

    for {
        select {
        case cpuU = <-c: // Will be received according to ticker in cpu()
            memStats = <-c

            a.show(time.Now().Sub(startTime), cpuU, memStats)

        case n := <-quit: // Will only be used once
            a.show(time.Now().Sub(startTime), cpuU, memStats)

            if n != nil {
                fmt.Println()
                fmt.Printf("Success!\n")
                steps = output(n)
                fmt.Printf("Numbers of steps: %d\n", steps)
            } else {
                fmt.Println()
                fmt.Println("Failed to find a valid path!")
            }
            return
        }
    }
}

func (a *Astar) astar(quit chan *Board) {
    for {
        node := a.open.pop()
        if node == nil { // Queue empty and all threads sleeping as well
            quit <- nil
            return
        }

        atomic.AddUint32(&a.opened, 1)

        for _, n := range node.neighbours() {
            n.parent = node

            if isGoal(n, a.goal) {
                quit <- n
                return
            }

            nSum := sum(n)

            n.g = node.g + 1
            n.f = n.g + n.priority(a.goal)

            m, ok := a.closed.get(nSum)

            if ok && m < n.f {
                pool.Put(n.cells)
                continue
            }

            a.open.push(n)

            a.closed.set(nSum, n.f)
        }

        if l := uint32(a.closed.len()); l > atomic.LoadUint32(&a.states) {
            atomic.StoreUint32(&a.states, l)
        }
    }
}
