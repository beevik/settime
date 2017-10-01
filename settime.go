package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

var emptyTime time.Time

var usage = `Usage: settime [HOST]
Set the system clock using the time reported by the NTP server at the IP
address 'HOST'.`

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println(usage)
		os.Exit(0)
	}

	tm, err := setTime(args[0])
	if err != nil {
		fmt.Printf("Time could not be set: %v\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Time sucessfully updated to %v\n", tm)
}

func setTime(host string) (time.Time, error) {
	r, err := ntp.Query(host)
	if err != nil {
		return emptyTime, err
	}

	t := time.Now().Add(r.ClockOffset)
	err = setPlatformTime(t)
	if err != nil {
		return emptyTime, err
	}
	return t, err
}
