package main

import (
    "fmt"
    "os"
)

func output(board *Board) uint16 {
    var file *os.File
    var err error

    if outputF != "" {
        file, err = os.Create(outputF)
        if err != nil {
            fmt.Println(err)
        }
    }

    moves := []*Board{board}
    node := board
    for node != nil {
        moves = append(moves, node)
        node = node.parent
    }

    if file != nil {
        for i := range moves {
            fmt.Fprintf(file, "Step number: %d\n", i+1)
            file.WriteString(moves[(len(moves)-1)-i].String())
            file.Write([]byte{'\n', '\n'})
        }

        fmt.Println("Steps written to file:", file.Name())
    }

    return uint16(len(moves))
}
