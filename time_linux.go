// +build linux,cgo
package main

/*
#include <time.h>
#include <errno.h>

int set_time(size_t t) {
	int err;
	const time_t tm = (time_t)t;
	err = stime(&tm);
	if (err != 0) {
		switch (errno) {
			case EPERM:
				return 1;
			case EFAULT:
				return 2;
			default:
				return 3;
		}
	}
	return 0;
}
*/
import "C"

import (
	"errors"
	"time"
)

const nanoPerSec = 1000000000

var (
	unixEpoch = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
)

func setPlatformTime(t time.Time) error {
	elapsed := int64(t.Sub(unixEpoch).Nanoseconds() / nanoPerSec)
	errno := C.set_time(C.size_t(elapsed))
	switch errno {
	case 0:
		return nil
	case 1:
		return errors.New("insufficient permissions")
	case 2:
		return errors.New("error getting information from user space")
	default:
		return errors.New("unknown error")
	}
}
