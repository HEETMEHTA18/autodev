//go:build windows

package cmd

import "os"

func notifyResize(ch chan<- os.Signal) {}

func stopResize(ch chan<- os.Signal) {}

func resizeSignal() os.Signal { return nil }
