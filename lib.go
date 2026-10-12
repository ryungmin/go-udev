//go:build linux

package udev

import (
	"sync"

	"github.com/ebitengine/purego"
)

// libudev functions, bound at runtime by loadLib.
// Pointers to libudev objects are passed as uintptr. purego copies string
// arguments to NUL-terminated C strings and copies char* results to Go
// strings (NULL becomes "").
var (
	udev_new   func() uintptr
	udev_unref func(udev uintptr) uintptr

	udev_list_entry_get_next  func(entry uintptr) uintptr
	udev_list_entry_get_name  func(entry uintptr) string
	udev_list_entry_get_value func(entry uintptr) string

	udev_device_new_from_syspath           func(udev uintptr, syspath string) uintptr
	udev_device_new_from_devnum            func(udev uintptr, typ byte, devnum uint64) uintptr
	udev_device_new_from_subsystem_sysname func(udev uintptr, subsystem, sysname string) uintptr
	udev_device_new_from_device_id         func(udev uintptr, id string) uintptr
	udev_device_ref                        func(dev uintptr) uintptr
	udev_device_unref                      func(dev uintptr) uintptr

	udev_device_get_parent                        func(dev uintptr) uintptr
	udev_device_get_parent_with_subsystem_devtype func(dev uintptr, subsystem, devtype string) uintptr
	udev_device_get_devpath                       func(dev uintptr) string
	udev_device_get_subsystem                     func(dev uintptr) string
	udev_device_get_devtype                       func(dev uintptr) string
	udev_device_get_sysname                       func(dev uintptr) string
	udev_device_get_syspath                       func(dev uintptr) string
	udev_device_get_sysnum                        func(dev uintptr) string
	udev_device_get_devnode                       func(dev uintptr) string
	udev_device_get_is_initialized                func(dev uintptr) int32
	udev_device_get_devlinks_list_entry           func(dev uintptr) uintptr
	udev_device_get_properties_list_entry         func(dev uintptr) uintptr
	udev_device_get_tags_list_entry               func(dev uintptr) uintptr
	udev_device_get_sysattr_list_entry            func(dev uintptr) uintptr
	udev_device_get_property_value                func(dev uintptr, key string) string
	udev_device_get_driver                        func(dev uintptr) string
	udev_device_get_devnum                        func(dev uintptr) uint64
	udev_device_get_action                        func(dev uintptr) string
	udev_device_get_seqnum                        func(dev uintptr) uint64
	udev_device_get_usec_since_initialized        func(dev uintptr) uint64
	udev_device_get_sysattr_value                 func(dev uintptr, sysattr string) string
	udev_device_set_sysattr_value                 func(dev uintptr, sysattr, value string) int32
	udev_device_has_tag                           func(dev uintptr, tag string) int32

	udev_enumerate_new                      func(udev uintptr) uintptr
	udev_enumerate_unref                    func(e uintptr) uintptr
	udev_enumerate_add_match_subsystem      func(e uintptr, subsystem string) int32
	udev_enumerate_add_nomatch_subsystem    func(e uintptr, subsystem string) int32
	udev_enumerate_add_match_sysattr        func(e uintptr, sysattr, value string) int32
	udev_enumerate_add_nomatch_sysattr      func(e uintptr, sysattr, value string) int32
	udev_enumerate_add_match_property       func(e uintptr, property, value string) int32
	udev_enumerate_add_match_sysname        func(e uintptr, sysname string) int32
	udev_enumerate_add_match_tag            func(e uintptr, tag string) int32
	udev_enumerate_add_match_parent         func(e uintptr, parent uintptr) int32
	udev_enumerate_add_match_is_initialized func(e uintptr) int32
	udev_enumerate_add_syspath              func(e uintptr, syspath string) int32
	udev_enumerate_scan_devices             func(e uintptr) int32
	udev_enumerate_scan_subsystems          func(e uintptr) int32
	udev_enumerate_get_list_entry           func(e uintptr) uintptr

	udev_monitor_new_from_netlink                   func(udev uintptr, name string) uintptr
	udev_monitor_unref                              func(m uintptr) uintptr
	udev_monitor_set_receive_buffer_size            func(m uintptr, size int32) int32
	udev_monitor_filter_add_match_subsystem_devtype func(m uintptr, subsystem string, devtype *byte) int32
	udev_monitor_filter_add_match_tag               func(m uintptr, tag string) int32
	udev_monitor_filter_update                      func(m uintptr) int32
	udev_monitor_filter_remove                      func(m uintptr) int32
	udev_monitor_receive_device                     func(m uintptr) uintptr
	udev_monitor_enable_receiving                   func(m uintptr) int32
	udev_monitor_get_fd                             func(m uintptr) int32
)

// loadLib opens libudev at runtime (once per process) and binds the functions above.
// A host without libudev only makes the constructors return nil; the process still starts.
var loadLib = sync.OnceValue(func() error {
	var h uintptr
	var err error
	for _, name := range []string{"libudev.so.1", "libudev.so"} {
		if h, err = purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL); err == nil {
			break
		}
	}
	if err != nil {
		return err
	}

	for name, fn := range map[string]any{
		"udev_new":   &udev_new,
		"udev_unref": &udev_unref,

		"udev_list_entry_get_next":  &udev_list_entry_get_next,
		"udev_list_entry_get_name":  &udev_list_entry_get_name,
		"udev_list_entry_get_value": &udev_list_entry_get_value,

		"udev_device_new_from_syspath":           &udev_device_new_from_syspath,
		"udev_device_new_from_devnum":            &udev_device_new_from_devnum,
		"udev_device_new_from_subsystem_sysname": &udev_device_new_from_subsystem_sysname,
		"udev_device_new_from_device_id":         &udev_device_new_from_device_id,
		"udev_device_ref":                        &udev_device_ref,
		"udev_device_unref":                      &udev_device_unref,

		"udev_device_get_parent":                        &udev_device_get_parent,
		"udev_device_get_parent_with_subsystem_devtype": &udev_device_get_parent_with_subsystem_devtype,
		"udev_device_get_devpath":                       &udev_device_get_devpath,
		"udev_device_get_subsystem":                     &udev_device_get_subsystem,
		"udev_device_get_devtype":                       &udev_device_get_devtype,
		"udev_device_get_sysname":                       &udev_device_get_sysname,
		"udev_device_get_syspath":                       &udev_device_get_syspath,
		"udev_device_get_sysnum":                        &udev_device_get_sysnum,
		"udev_device_get_devnode":                       &udev_device_get_devnode,
		"udev_device_get_is_initialized":                &udev_device_get_is_initialized,
		"udev_device_get_devlinks_list_entry":           &udev_device_get_devlinks_list_entry,
		"udev_device_get_properties_list_entry":         &udev_device_get_properties_list_entry,
		"udev_device_get_tags_list_entry":               &udev_device_get_tags_list_entry,
		"udev_device_get_sysattr_list_entry":            &udev_device_get_sysattr_list_entry,
		"udev_device_get_property_value":                &udev_device_get_property_value,
		"udev_device_get_driver":                        &udev_device_get_driver,
		"udev_device_get_devnum":                        &udev_device_get_devnum,
		"udev_device_get_action":                        &udev_device_get_action,
		"udev_device_get_seqnum":                        &udev_device_get_seqnum,
		"udev_device_get_usec_since_initialized":        &udev_device_get_usec_since_initialized,
		"udev_device_get_sysattr_value":                 &udev_device_get_sysattr_value,
		"udev_device_set_sysattr_value":                 &udev_device_set_sysattr_value,
		"udev_device_has_tag":                           &udev_device_has_tag,

		"udev_enumerate_new":                      &udev_enumerate_new,
		"udev_enumerate_unref":                    &udev_enumerate_unref,
		"udev_enumerate_add_match_subsystem":      &udev_enumerate_add_match_subsystem,
		"udev_enumerate_add_nomatch_subsystem":    &udev_enumerate_add_nomatch_subsystem,
		"udev_enumerate_add_match_sysattr":        &udev_enumerate_add_match_sysattr,
		"udev_enumerate_add_nomatch_sysattr":      &udev_enumerate_add_nomatch_sysattr,
		"udev_enumerate_add_match_property":       &udev_enumerate_add_match_property,
		"udev_enumerate_add_match_sysname":        &udev_enumerate_add_match_sysname,
		"udev_enumerate_add_match_tag":            &udev_enumerate_add_match_tag,
		"udev_enumerate_add_match_parent":         &udev_enumerate_add_match_parent,
		"udev_enumerate_add_match_is_initialized": &udev_enumerate_add_match_is_initialized,
		"udev_enumerate_add_syspath":              &udev_enumerate_add_syspath,
		"udev_enumerate_scan_devices":             &udev_enumerate_scan_devices,
		"udev_enumerate_scan_subsystems":          &udev_enumerate_scan_subsystems,
		"udev_enumerate_get_list_entry":           &udev_enumerate_get_list_entry,

		"udev_monitor_new_from_netlink":                   &udev_monitor_new_from_netlink,
		"udev_monitor_unref":                              &udev_monitor_unref,
		"udev_monitor_set_receive_buffer_size":            &udev_monitor_set_receive_buffer_size,
		"udev_monitor_filter_add_match_subsystem_devtype": &udev_monitor_filter_add_match_subsystem_devtype,
		"udev_monitor_filter_add_match_tag":               &udev_monitor_filter_add_match_tag,
		"udev_monitor_filter_update":                      &udev_monitor_filter_update,
		"udev_monitor_filter_remove":                      &udev_monitor_filter_remove,
		"udev_monitor_receive_device":                     &udev_monitor_receive_device,
		"udev_monitor_enable_receiving":                   &udev_monitor_enable_receiving,
		"udev_monitor_get_fd":                             &udev_monitor_get_fd,
	} {
		sym, err := purego.Dlsym(h, name)
		if err != nil {
			return err
		}
		purego.RegisterFunc(fn, sym)
	}
	return nil
})
