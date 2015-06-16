package main

import (
    "runtime"
    "syscall"
    "time"
)

func cpu(c chan int) {
    prevTime := time.Now().UnixNano()
    var prevUsage int64
    var rusage syscall.Rusage
    var memstats runtime.MemStats

    t := time.Tick(status)

    for {
        select {
        case <-t:
            syscall.Getrusage(syscall.RUSAGE_SELF, &rusage)

            curTime := time.Now().UnixNano()
            timeDiff := curTime - prevTime
            curUsage := rusage.Utime.Nano() + rusage.Stime.Nano()
            usageDiff := curUsage - prevUsage

            cpuUsagePercent := 100 * float64(usageDiff) / float64(timeDiff)
            prevTime = curTime
            prevUsage = curUsage

            runtime.ReadMemStats(&memstats)

            c <- int(cpuUsagePercent)
            c <- int(memstats.Alloc / 1024 / 1024)
        }
    }
}
