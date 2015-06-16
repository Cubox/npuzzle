package main

import (
    "sync"
)

type Queue struct {
    m map[uint16]*Head
    sync.Cond
    *sync.RWMutex
    size uint32
    smallest uint16
    count    uint8
}

type Head struct {
    *Node
    count uint32
}

type Node struct {
    *Board
    next *Node
}

type Map struct {
    m   map[uint32]uint16
    *sync.RWMutex
}

func NewMap() Map {
    m := Map{}
    m.m = make(map[uint32]uint16)
    m.RWMutex = new(sync.RWMutex)

    return m
}

func (m *Map) set(key uint32, f uint16) {
    m.Lock()
    defer m.Unlock()

    m.m[key] = f
}

func (m *Map) get(key uint32) (uint16, bool) {
    m.RLock()
    defer m.RUnlock()

    val, ok := m.m[key]
    return val, ok
}

func (m *Map) del(key uint32) {
    m.Lock()
    defer m.Unlock()

    delete(m.m, key)
}

func (m *Map) len() int {
    m.RLock()
    defer m.RUnlock()

    return len(m.m)
}

func NewQueue(initial *Board) Queue {
    queue := Queue{size: 1}
    queue.m = make(map[uint16]*Head)
    queue.m[initial.f] = &Head{Node: &Node{Board: initial}, count: 1}

    queue.RWMutex = new(sync.RWMutex)
    queue.Cond.L = queue.RWMutex

    return queue
}

func (q *Queue) push(board *Board) {
    q.Lock()
    defer q.Unlock()
    defer q.Signal()

    q.size++
    head := q.m[board.f]

    if head == nil {
        head = &Head{}
        q.m[board.f] = head
    }

    node := Node{Board: board, next: head.Node}

    head.Node = &node
    head.count++

    if q.smallest == 0 || board.f < q.smallest {
        q.smallest = board.f
    }
}

func (q *Queue) pop() *Board {
    q.Lock()
    defer q.Unlock()

    for len(q.m) == 0 {
        if q.count >= uint8(threads)-1 { // Deadlock if all threads are waiting here
            return nil
        }

        q.count++
        q.Wait()
        q.count--
    }

    q.size--
    head := q.m[q.smallest]
    board := head.Board

    if head.next != nil {
        head.Node = head.next
    } else { // We need to find the next best priority
        delete(q.m, q.smallest)

        if len(q.m) == 0 {
            q.smallest = 0
            return board
        }

        for { // We will never hit an endless loop
            q.smallest++
            _, ok := q.m[q.smallest]
            if ok {
                return board
            }
        }
    }

    return board
}

func (q *Queue) len() uint32 {
    q.RLock()
    defer q.RUnlock()

    return q.size
}
