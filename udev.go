//go:build linux

// Package udev provides a wrapper around the libudev C library.
// libudev is loaded at runtime with purego, so cgo is not required.
package udev

import (
	"runtime"
	"sync"
)

// Udev is an opaque struct wraping a udev library context
type Udev struct {
	// A pointer to the C struct udev context (0 if libudev could not be loaded)
	ptr uintptr
	// Mutex for thread sync as libudev is not thread safe when called with the same struct udev
	m    sync.Mutex
	once sync.Once
}

func udevUnref(u *Udev) {
	udev_unref(u.ptr)
}

// Lock locks a udev context
func (u *Udev) lock() {
	u.once.Do(func() {
		if loadLib() != nil {
			return
		}
		u.ptr = udev_new()
		if u.ptr != 0 {
			runtime.SetFinalizer(u, udevUnref)
		}
	})
	u.m.Lock()
}

// Unlock unlocks a udev context
func (u *Udev) unlock() {
	u.m.Unlock()
}

// newDevice is a private helper function and returns a pointer to a new device.
// The device is also added t the devices map in the udev context.
// The agrument ptr is a pointer to the underlying C udev_device structure.
// The function returns nil if the pointer passed is NULL.
func (u *Udev) newDevice(ptr uintptr) (d *Device) {
	// If passed a NULL pointer, return nil
	if ptr == 0 {
		return nil
	}
	// Create a new device object
	d = &Device{
		ptr: ptr,
		u:   u,
	}
	runtime.SetFinalizer(d, deviceUnref)
	// Return the device object
	return
}

// newMonitor is a private helper function and returns a pointer to a new monitor.
// The monitor is also added t the monitors map in the udev context.
// The agrument ptr is a pointer to the underlying C udev_monitor structure.
// The function returns nil if the pointer passed is NULL.
func (u *Udev) newMonitor(ptr uintptr) (m *Monitor) {
	// If passed a NULL pointer, return nil
	if ptr == 0 {
		return nil
	}
	// Create a new device object
	m = &Monitor{
		ptr: ptr,
		u:   u,
	}
	runtime.SetFinalizer(m, monitorUnref)
	// Return the device object
	return
}

func (u *Udev) newEnumerate(ptr uintptr) (e *Enumerate) {
	// If passed a NULL pointer, return nil
	if ptr == 0 {
		return nil
	}
	// Create a new device object
	e = &Enumerate{
		ptr: ptr,
		u:   u,
	}
	runtime.SetFinalizer(e, enumerateUnref)
	// Return the device object
	return
}

// NewDeviceFromSyspath returns a pointer to a new device identified by its syspath, and nil on error
// The device is identified by the syspath argument
func (u *Udev) NewDeviceFromSyspath(syspath string) *Device {
	// Lock the udev context
	u.lock()
	defer u.unlock()
	if u.ptr == 0 {
		return nil
	}
	// Return a new device
	return u.newDevice(udev_device_new_from_syspath(u.ptr, syspath))
}

// NewDeviceFromDevnum returns a pointer to a new device identified by its Devnum, and nil on error
// deviceType is 'c' for a character device and 'b' for a block device
func (u *Udev) NewDeviceFromDevnum(deviceType uint8, n Devnum) *Device {
	u.lock()
	defer u.unlock()
	if u.ptr == 0 {
		return nil
	}
	return u.newDevice(udev_device_new_from_devnum(u.ptr, deviceType, n.d))
}

// NewDeviceFromSubsystemSysname returns a pointer to a new device identified by its subystem and sysname, and nil on error
func (u *Udev) NewDeviceFromSubsystemSysname(subsystem, sysname string) *Device {
	u.lock()
	defer u.unlock()
	if u.ptr == 0 {
		return nil
	}
	return u.newDevice(udev_device_new_from_subsystem_sysname(u.ptr, subsystem, sysname))
}

// NewDeviceFromDeviceID returns a pointer to a new device identified by its device id, and nil on error
func (u *Udev) NewDeviceFromDeviceID(id string) *Device {
	u.lock()
	defer u.unlock()
	if u.ptr == 0 {
		return nil
	}
	return u.newDevice(udev_device_new_from_device_id(u.ptr, id))
}

// NewEnumerate returns a pointer to a new enumerate, and nil on error
func (u *Udev) NewEnumerate() *Enumerate {
	u.lock()
	defer u.unlock()
	if u.ptr == 0 {
		return nil
	}
	return u.newEnumerate(udev_enumerate_new(u.ptr))
}

// NewMonitorFromNetlink returns a pointer to a new monitor listening to a NetLink socket, and nil on error
// The name argument is either "kernel" or "udev".
// When passing "kernel" the events are received before they are processed by udev.
// When passing "udev" the events are received after udev has processed the events and created device nodes.
// In most cases you will want to use "udev".
func (u *Udev) NewMonitorFromNetlink(name string) *Monitor {
	u.lock()
	defer u.unlock()
	if u.ptr == 0 {
		return nil
	}
	return u.newMonitor(udev_monitor_new_from_netlink(u.ptr, name))
}
