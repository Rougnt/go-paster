//go:build windows

package main

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procSendInput    = user32.NewProc("SendInput")
	procVkKeyScanW   = user32.NewProc("VkKeyScanW")
	procMapVirtualKeyW = user32.NewProc("MapVirtualKeyW")
)

// Constants for INPUT structure
const (
	INPUT_KEYBOARD = 1
	KEYEVENTF_KEYUP     = 0x0002
	KEYEVENTF_UNICODE   = 0x0004
	KEYEVENTF_SCANCODE  = 0x0008
)

type KEYBDINPUT struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
}

type INPUT struct {
	type_ uint32
	ki    KEYBDINPUT
	padding [4]byte // Padding to match C struct alignment on 64-bit
}

// TypeString mimics the typing of a string using low-level SendInput with ScanCodes
func TypeString(str string, interval float64) {
	runes := []rune(str)
	for _, r := range runes {
		sendChar(r)
		time.Sleep(time.Duration(interval * float64(time.Second)))
	}
}

func sendChar(r rune) {
	// 1. Convert Rune to Virtual Key (VK) using current keyboard layout
	// VkKeyScanW returns: low byte = VK, high byte = shift state
	res, _, _ := procVkKeyScanW.Call(uintptr(r))
	vk := uint16(res & 0xff)
	shiftState := uint16((res >> 8) & 0xff)

	needShift := (shiftState & 1) != 0
	needCtrl  := (shiftState & 2) != 0
	needAlt   := (shiftState & 4) != 0

	// Helper to press/release keys
	press := func(vkCode uint16, scanCode uint16, flags uint32) {
		sendInput(vkCode, scanCode, flags)
	}
	release := func(vkCode uint16, scanCode uint16, flags uint32) {
		sendInput(vkCode, scanCode, flags|KEYEVENTF_KEYUP)
	}

	// 2. Resolve Scan Codes for modifiers
	// VK_SHIFT = 0x10, VK_CONTROL = 0x11, VK_MENU (Alt) = 0x12
	scanShift := mapVirtualKey(0x10)
	scanCtrl  := mapVirtualKey(0x11)
	scanAlt   := mapVirtualKey(0x12)

	// 3. Press Modifiers if needed
	if needShift { press(0x10, scanShift, KEYEVENTF_SCANCODE) }
	if needCtrl  { press(0x11, scanCtrl, KEYEVENTF_SCANCODE) }
	if needAlt   { press(0x12, scanAlt, KEYEVENTF_SCANCODE) }

	// 4. Get Scan Code for the actual character
	scanCode := mapVirtualKey(vk)

	// 5. Press and Release the character (Using SCANCODE flag is crucial for KVMs!)
	press(vk, scanCode, KEYEVENTF_SCANCODE)
	// Small hardware debounce delay
	time.Sleep(5 * time.Millisecond) 
	release(vk, scanCode, KEYEVENTF_SCANCODE)

	// 6. Release Modifiers
	if needAlt   { release(0x12, scanAlt, KEYEVENTF_SCANCODE) }
	if needCtrl  { release(0x11, scanCtrl, KEYEVENTF_SCANCODE) }
	if needShift { release(0x10, scanShift, KEYEVENTF_SCANCODE) }
}

func mapVirtualKey(vk uint16) uint16 {
	// MapVirtualKeyW(vk, 0) -> returns scan code
	ret, _, _ := procMapVirtualKeyW.Call(uintptr(vk), 0)
	return uint16(ret)
}

func sendInput(vk uint16, scan uint16, flags uint32) {
	var i INPUT
	i.type_ = INPUT_KEYBOARD
	i.ki.wVk = vk
	i.ki.wScan = scan
	i.ki.dwFlags = flags

	// SendInput expects an array of INPUT structs. We send 1.
	// Size of INPUT struct is 40 bytes on 64-bit windows (usually).
	procSendInput.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&i)),
		uintptr(unsafe.Sizeof(i)),
	)
}
