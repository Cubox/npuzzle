package main

import (
    "fmt"
)

func print(format string, args... interface{}) {
    if !dumb {
        fmt.Printf("\r\033[2K") // Clean current line
    }

    fmt.Printf(format, args...)

    if dumb && format[len(format)-1] != '\n' {
        fmt.Println()
    }
}
