//go:build linux

package udev

import "golang.org/x/sys/unix"

// Devnum is a kernel device number
type Devnum struct {
	d uint64
}

// Major returns the major part of a Devnum
func (d Devnum) Major() int {
	return int(unix.Major(d.d))
}

// Minor returns the minor part of a Devnum
func (d Devnum) Minor() int {
	return int(unix.Minor(d.d))
}

// MkDev creates a Devnum from a major and minor number
func MkDev(major, minor int) Devnum {
	return Devnum{unix.Mkdev(uint32(major), uint32(minor))}
}
