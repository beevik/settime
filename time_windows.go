// +build windows,cgo
package main

/*
#include <windows.h>
#include <stdint.h>

int set_time(uint32_t ticksUpper, uint32_t ticksLower) {
	FILETIME fileTime;
	fileTime.dwLowDateTime  = (DWORD)ticksLower;
	fileTime.dwHighDateTime = (DWORD)ticksUpper;

	SYSTEMTIME sysTime;
	BOOL result = FileTimeToSystemTime(&fileTime, &sysTime);
	if (!result)
		return GetLastError();

	result = SetSystemTime(&sysTime);
	if (!result)
		return GetLastError();

	return 0;
}
*/
import "C"

import (
	"fmt"
	"time"
)

const (
	nsecPerTick        = 100
	windowsEpochOffset = 0x019db1ded53e8000 // ticks[UnixEpoch] - ticks[WindowsEpoch]
)

var (
	unixEpoch = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
)

func setPlatformTime(t time.Time) error {
	// Calculate the number of 100ns ticks since the Windows epoch (Jan 1,
	// 1601). We have to start with the time since the unix epoch, since go's
	// Duration field cannot hold a duration longer than ~290 years.
	nsec := t.Sub(unixEpoch).Nanoseconds()
	ticks := uint64(nsec)/nsecPerTick + windowsEpochOffset

	// Split ticks into upper and lower dword values.
	ticksUpper := uint32(ticks >> 32)
	ticksLower := uint32(ticks & 0xffffffff)

	// Call the C function to set the system time.
	errcode := C.set_time(C.uint32_t(ticksUpper), C.uint32_t(ticksLower))
	switch errcode {
	case 0:
		return nil
	case 19:
		fallthrough
	case 1314:
		return fmt.Errorf("insufficient permissions")
	default:
		return fmt.Errorf("unable to set system time (errcode=%v)", errcode)
	}
}
