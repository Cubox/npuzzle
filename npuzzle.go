package main

import (
    "flag"
    "log"
    "runtime"
    //"github.com/davecheney/profile"
    _ "net/http/pprof"
    "net/http"
    "time"
)

var (
    file       string
    methodName string
    showSteps  bool
    algo       string
    dumb       bool
    outputF    string
    status     time.Duration
    threads    int
    genSize    int
)

type Algorithm interface {
    solve(initial *Board) (steps uint16)
}

func main() {
    go func() {
	log.Println(http.ListenAndServe("localhost:6060", nil))
}()
    //defer profile.Start(profile.CPUProfile).Stop()
    flag.StringVar(&file, "f", "", "Specify an input file. If empty, will generate randomly")
    flag.IntVar(&genSize, "s", 4, "Size of the puzzle generated if no file are given")
    flag.StringVar(&methodName, "m", "linear", "Method used to determine priority")
    flag.StringVar(&algo, "a", "a", "Algorithm to use")
    flag.StringVar(&outputF, "o", "", "File where to output calculated steps")
    flag.BoolVar(&dumb, "dumb", false, "Interactive stats?")
    flag.DurationVar(&status, "status", time.Second, "Interval to print the calculation status")
    flag.IntVar(&threads, "t", 0, "Threads to spawn. 0 for default value.")
    flag.Parse()

    var board Board
    board.parent = nil
    if file != "" {
        board.parseBoard(file)
    } else {
        board.generateBoard()
    }

    switch methodName {
    case "manhattan":
        method = manhattanCell
    case "hamming":
        method = hammingCell
    case "linear":
        method = linearConflictCell
    default:
        log.Fatalln("Unknown method:", methodName)
    }

    if threads == 0 {
        threads = runtime.NumCPU() * 2
    }

    runtime.GOMAXPROCS(threads)

    var a Algorithm

    switch algo {
    case "a":
        a = new(Astar)
    case "ida":
        a = new(Ida)
    default:
        log.Fatalln("Unknown algorithm: ", algo)
    }

    a.solve(&board)
}
