package main

import (
    // "fmt"
    // "sync/atomic"
)

type Ida struct {
    // running uint32
    // found   chan *Board
}

func (a *Ida) search(node *Board, depth, bound uint16, c chan uint16, routine bool) {
    // end := func() {
    //     if routine {
    //         atomic.AddUint32(&a.running, ^uint32(0))
    //     }
    // }
    //
    // f := depth + node.priority()
    // if f > bound {
    //     end()
    //     c <- f
    //     return
    // }
    //
    // if isGoal(node) {
    //     a.found <- node
    //     return
    // }
    //
    // var min uint16
    // childC := make(chan uint16, 4)
    //
    // neighbours := node.neighbours()
    //
    // for _, n := range neighbours {
    //     n.parent = node
    //     if atomic.LoadUint32(&a.running) < uint32(threads) {
    //         atomic.AddUint32(&a.running, 1)
    //         go a.search(n, depth+node.priority(), bound, childC, true)
    //     } else {
    //         a.search(n, depth+node.priority(), bound, childC, false)
    //     }
    // }
    //
    // end()
    //
    // for i := 0; i < len(neighbours); i++ {
    //     result := <-childC
    //     if min == 0 || result < min {
    //         min = result
    //     }
    // }
    //
    // c <- min
}

func (a *Ida) solve(initial *Board) (steps uint16) {
    // initGoal()
    //
    // fmt.Println("Starting path search.")
    // bound := initial.priority()
    //
    // c := make(chan uint16, 1)
    // a.found = make(chan *Board)
    //
    // for {
    //     atomic.AddUint32(&a.running, 1)
    //     go a.search(initial, 0, bound, c, true)
    //
    //     select {
    //     case t := <-c:
    //         bound = t
    //     case result := <-a.found:
    //         fmt.Println("found!!!!")
    //         output(result)
    //         return
    //     }
    // }
    return
}
