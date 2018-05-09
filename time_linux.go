// +build linux,cgo

package main

/*
#include <sys/time.h>
#include <time.h>
#include <errno.h>

int set_time_of_day(size_t sec, size_t usec) {
	struct timeval tv;
	tv.tv_sec = (time_t)sec;
	tv.tv_usec = (time_t)usec;

	int err = settimeofday(&tv, NULL);
	if (err) {
		switch (errno) {
			case EPERM:
				return 1;
			case EFAULT:
				return 2;
			case EINVAL:
				return 3;
			default:
				return 4;
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

const (
	nsecPerUsec = 1000
	usecPerSec  = 1000000
)

var (
	unixEpoch = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
)

func setPlatformTime(t time.Time) error {
	usec := int64(t.Sub(unixEpoch).Nanoseconds() / nsecPerUsec)
	sec := usec / usecPerSec
	usec -= sec * usecPerSec
	errno := C.set_time_of_day(C.size_t(sec), C.size_t(usec))
	switch errno {
	case 0:
		return nil
	case 1:
		return errors.New("insufficient permissions")
	case 2:
		return errors.New("error getting information from user space")
	case 3:
		return errors.New("invalid data")
	default:
		return errors.New("unknown error")
	}
}
