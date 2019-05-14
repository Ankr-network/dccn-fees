package main

import (
    "log"
    "time"
    "fmt"
)

func main() {
    now := time.Now()
    currentYear, currentMonth, _ := now.Date()
    currentLocation := now.Location()

    firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
    startLastMonth := firstOfMonth.AddDate(0, -1, 0)
    endlastMonth := firstOfMonth.AddDate(0, 0, -1)

    start := startLastMonth.Unix()
    end := endlastMonth.Unix()

    fmt.Println(startLastMonth)
    fmt.Println(endlastMonth)

    log.Printf("start %d %d \n", start, end)
}
