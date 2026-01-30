//go:build !windows

package main

import "fmt"

// TypeString stub for non-windows systems
func TypeString(str string, interval float64) {
	fmt.Println("Typing simulation is only available on Windows for this version.")
	fmt.Printf("Mock typing: %s\n", str)
}
