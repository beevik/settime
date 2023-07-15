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

	const format = "Mon Jan _2 2006 15:04:05 (MST)"
	fmt.Printf("Time sucessfully set to: %s\n", tm.Format(format))
}

func setTime(host string) (time.Time, error) {
	r, err := ntp.Query(host)
	if err != nil {
		return emptyTime, err
	}

	err = r.Validate()
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
