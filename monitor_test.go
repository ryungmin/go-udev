//go:build linux

package udev

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func ExampleMonitor() {

	// Create Udev and Monitor
	u := Udev{}
	m := u.NewMonitorFromNetlink("udev")

	// Add filters to monitor
	_ = m.FilterAddMatchSubsystemDevtype("block", "disk")
	_ = m.FilterAddMatchTag("systemd")

	// Create a context
	ctx, cancel := context.WithCancel(context.Background())

	// Start monitor goroutine and get receive channel
	ch, errCh, _ := m.DeviceChan(ctx)

	// WaitGroup for timers
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		fmt.Println("Started listening on channel")
	Loop:
		for {
			select {
			case <-ctx.Done():
				break Loop
			case d := <-ch:
				fmt.Println("Event:", d.Syspath(), d.Action())
			case e := <-errCh:
				fmt.Println("Error:", e)
			}
		}
		fmt.Println("Channel closed")
		wg.Done()
	}()
	go func() {
		fmt.Println("Starting timer to update filter")
		<-time.After(2 * time.Second)
		fmt.Println("Removing filter")
		_ = m.FilterRemove()
		fmt.Println("Updating filter")
		_ = m.FilterUpdate()
		wg.Done()
	}()
	go func() {
		fmt.Println("Starting timer to signal done")
		<-time.After(4 * time.Second)
		fmt.Println("Signalling done")
		cancel()
		wg.Done()
	}()
	wg.Wait()
}

func TestMonitorDeviceChan(t *testing.T) {
	u := Udev{}
	m := u.NewMonitorFromNetlink("udev")
	_ = m.FilterAddMatchSubsystemDevtype("block", "disk")
	_ = m.FilterAddMatchTag("systemd")
	ctx, cancel := context.WithCancel(context.Background())
	ch, errCh, e := m.DeviceChan(ctx)
	if e != nil {
		t.Fail()
	}
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		fmt.Println("Started listening on channel")
	Loop:
		for {
			select {
			case <-ctx.Done():
				break Loop
			case d := <-ch:
				fmt.Println(d.Syspath(), d.Action())
			case e := <-errCh:
				fmt.Println("Error:", e)
			}
		}
		fmt.Println("Channel closed")
		wg.Done()
	}()
	go func() {
		fmt.Println("Starting timer to update filter")
		<-time.After(2 * time.Second)
		fmt.Println("Removing filter")
		_ = m.FilterRemove()
		fmt.Println("Updating filter")
		_ = m.FilterUpdate()
		wg.Done()
	}()
	go func() {
		fmt.Println("Starting timer to signal done")
		<-time.After(4 * time.Second)
		fmt.Println("Signalling done")
		cancel()
		wg.Done()
	}()
	wg.Wait()

}

func TestMonitorGC(t *testing.T) {
	runtime.GC()
}
