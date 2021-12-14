package input

import (
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/gonutz/w32/v2"
)

type errString string

func (e errString) Error() string {
	return string(e)
}

const (
	// ErrBlocked happens when the winapi function SendInput is blocked by User
	// Interface Privilege Isolation.
	ErrBlocked errString = "w32.SendInput returned 0, meaning input was blocked"

	// ErrGetCursorFailed happens when the winapi function GetCursorPos fails.
	ErrGetCursorFailed errString = "w32.GetCursorPos failed"

	// ErrSetCursorFailed happens when the winapi function SetCursorPos fails.
	ErrSetCursorFailed errString = "w32.SetCursorPos failed"
)

func clickAt(x, y int, down, up uint32) error {
	if !w32.SetCursorPos(x, y) {
		return ErrSetCursorFailed
	}
	return click(down, up)
}

func click(down, up uint32) error {
	n := w32.SendInput(
		w32.MouseInput(w32.MOUSEINPUT{Flags: down}),
		w32.MouseInput(w32.MOUSEINPUT{Flags: up}),
	)
	if n == 0 {
		return ErrBlocked
	}
	return nil
}

// LeftClickAt moves the mouse to screen coordinates x,y and clicks the left
// mouse button.
func LeftClickAt(x, y int) error {
	return clickAt(x, y, w32.MOUSEEVENTF_LEFTDOWN, w32.MOUSEEVENTF_LEFTUP)
}

// RightClickAt moves the mouse to screen coordinates x,y and clicks the right
// mouse button.
func RightClickAt(x, y int) error {
	return clickAt(x, y, w32.MOUSEEVENTF_RIGHTDOWN, w32.MOUSEEVENTF_RIGHTUP)
}

// MiddleClickAt moves the mouse to screen coordinates x,y and clicks the middle
// mouse button.
func MiddleClickAt(x, y int) error {
	return clickAt(x, y, w32.MOUSEEVENTF_MIDDLEDOWN, w32.MOUSEEVENTF_MIDDLEUP)
}

// LeftClick clicks the left mouse button, i.e. presses and releases it.
func LeftClick() error {
	return click(w32.MOUSEEVENTF_LEFTDOWN, w32.MOUSEEVENTF_LEFTUP)
}

// RightClick clicks the right mouse button, i.e. presses and releases it.
func RightClick() error {
	return click(w32.MOUSEEVENTF_RIGHTDOWN, w32.MOUSEEVENTF_RIGHTUP)
}

// MiddleClick clicks the middle mouse button, i.e. presses and releases it.
func MiddleClick() error {
	return click(w32.MOUSEEVENTF_MIDDLEDOWN, w32.MOUSEEVENTF_MIDDLEUP)
}

func buttonDown(x, y int, down uint32) error {
	if !w32.SetCursorPos(x, y) {
		return ErrSetCursorFailed
	}
	n := w32.SendInput(
		w32.MouseInput(w32.MOUSEINPUT{Flags: down}),
	)
	if n == 0 {
		return ErrBlocked
	}
	return nil
}

func buttonUp(up uint32) error {
	n := w32.SendInput(
		w32.MouseInput(w32.MOUSEINPUT{Flags: up}),
	)
	if n == 0 {
		return ErrBlocked
	}
	return nil
}

// LeftDown presses the left mouse button down.
func LeftDown(x, y int) error {
	return buttonDown(x, y, w32.MOUSEEVENTF_LEFTDOWN)
}

// RightDown presses the right mouse button down.
func RightDown(x, y int) error {
	return buttonDown(x, y, w32.MOUSEEVENTF_RIGHTDOWN)
}

// MiddleDown presses the middle mouse button down.
func MiddleDown(x, y int) error {
	return buttonDown(x, y, w32.MOUSEEVENTF_MIDDLEDOWN)
}

// LeftUp releases the left mouse button.
func LeftUp() error {
	return buttonUp(w32.MOUSEEVENTF_LEFTUP)
}

// RightUp releases the right mouse button.
func RightUp() error {
	return buttonUp(w32.MOUSEEVENTF_RIGHTUP)
}

// MiddleUp releases the middle mouse button.
func MiddleUp() error {
	return buttonUp(w32.MOUSEEVENTF_MIDDLEUP)
}

// MoveMouseTo move the mouse cursor to the given screen coordinates.
func MoveMouseTo(x, y int) error {
	if !w32.SetCursorPos(x, y) {
		return ErrSetCursorFailed
	}
	return nil
}

// MoveMouseBy moves the mouse cursor by the given amount of pixels in x and y.
// Note that positive dy means the cursor is move down on the screen.
func MoveMouseBy(dx, dy int) error {
	x, y, ok := w32.GetCursorPos()
	if !ok {
		return ErrGetCursorFailed
	}
	if !w32.SetCursorPos(x+dx, y+dy) {
		return ErrSetCursorFailed
	}
	return nil
}

// LeftDoubleClickAt moves the mouse to screen coordinates x,y and clicks the
// left mouse button twice.
func LeftDoubleClickAt(x, y int) error {
	if !w32.SetCursorPos(x, y) {
		return ErrSetCursorFailed
	}
	n := w32.SendInput(
		w32.MouseInput(w32.MOUSEINPUT{Flags: w32.MOUSEEVENTF_LEFTDOWN}),
		w32.MouseInput(w32.MOUSEINPUT{Flags: w32.MOUSEEVENTF_LEFTUP}),
		w32.MouseInput(w32.MOUSEINPUT{Flags: w32.MOUSEEVENTF_LEFTDOWN}),
		w32.MouseInput(w32.MOUSEINPUT{Flags: w32.MOUSEEVENTF_LEFTUP}),
	)
	if n == 0 {
		return ErrBlocked
	}
	return nil
}

// LeftDoubleClick clicks the left mouse button twice.
func LeftDoubleClick() error {
	n := w32.SendInput(
		w32.MouseInput(w32.MOUSEINPUT{Flags: w32.MOUSEEVENTF_LEFTDOWN}),
		w32.MouseInput(w32.MOUSEINPUT{Flags: w32.MOUSEEVENTF_LEFTUP}),
		w32.MouseInput(w32.MOUSEINPUT{Flags: w32.MOUSEEVENTF_LEFTDOWN}),
		w32.MouseInput(w32.MOUSEINPUT{Flags: w32.MOUSEEVENTF_LEFTUP}),
	)
	if n == 0 {
		return ErrBlocked
	}
	return nil
}

// Type will write the given text using Alt+Numpad numbers. It will sleep the
// smallest, non-0 delay between two letters.
func Type(s string) error {
	return TypeWithDelay(s, 1)
}

// TypeWithDelay will write the given text using Alt+Numpad numbers. It will
// sleep the given delay between two letters.
func TypeWithDelay(s string, delay time.Duration) error {
	toScanCode := func(vk uint) uint16 {
		return uint16(w32.MapVirtualKey(vk, w32.MAPVK_VK_TO_VSC))
	}

	const (
		down = 0
		up   = 1
	)

	upDown := func(vk uint) [2]w32.INPUT {
		return [2]w32.INPUT{
			down: w32.KeyboardInput(w32.KEYBDINPUT{
				Scan:  toScanCode(vk),
				Flags: w32.KEYEVENTF_SCANCODE,
			}),

			up: w32.KeyboardInput(w32.KEYBDINPUT{
				Scan:  toScanCode(vk),
				Flags: w32.KEYEVENTF_SCANCODE | w32.KEYEVENTF_KEYUP,
			}),
		}
	}

	alt := upDown(w32.VK_LMENU)
	nums := [][2]w32.INPUT{
		upDown(w32.VK_NUMPAD0),
		upDown(w32.VK_NUMPAD1),
		upDown(w32.VK_NUMPAD2),
		upDown(w32.VK_NUMPAD3),
		upDown(w32.VK_NUMPAD4),
		upDown(w32.VK_NUMPAD5),
		upDown(w32.VK_NUMPAD6),
		upDown(w32.VK_NUMPAD7),
		upDown(w32.VK_NUMPAD8),
		upDown(w32.VK_NUMPAD9),
	}

	keys := []w32.INPUT{alt[down], nums[0][down], nums[0][up]}

	// Unify line breaks to '\r' which is the virtual key code for VK_RETURN.
	s = strings.Replace(s, "\r\n", "\r", -1)
	s = strings.Replace(s, "\n", "\r", -1)

	for _, r := range s {
		if r == '\r' {
			if err := PressKey(w32.VK_RETURN); err != nil {
				return err
			}
		} else if r == '\b' {
			if err := PressKey(w32.VK_BACK); err != nil {
				return err
			}
		} else {
			keys = keys[:3] // Keep Alt down and type 0.
			for _, digit := range fmt.Sprint(int(r)) {
				d := digit - '0'
				keys = append(keys, nums[d][down], nums[d][up])
			}
			keys = append(keys, alt[up])

			if w32.SendInput(keys...) == 0 {
				return ErrBlocked
			}
		}
		time.Sleep(delay)
	}
	return nil
}

// KeyDown presses the given key on the keyboard. The value must be a virtual
// keycode like 'A', '1' or VK_RETURN (you can use the constants in
// github.com/gonutz/w32 VK_...).
func KeyDown(key uint16) error {
	n := w32.SendInput(w32.KeyboardInput(w32.KEYBDINPUT{Vk: key}))
	if n == 0 {
		return ErrBlocked
	}
	return nil
}

// KeyUp releases the given key on the keyboard. The value must be a virtual
// keycode like 'A', '1' or VK_RETURN (you can use the constants in
// github.com/gonutz/w32 VK_...).
func KeyUp(key uint16) error {
	n := w32.SendInput(w32.KeyboardInput(w32.KEYBDINPUT{
		Vk:    key,
		Flags: w32.KEYEVENTF_KEYUP,
	}))
	if n == 0 {
		return ErrBlocked
	}
	return nil
}

// PressKey presses and releases the given key on the keyboard. The value must
// be a virtual keycode like 'A', '1' or VK_RETURN (you can use the constants in
// github.com/gonutz/w32 VK_...).
func PressKey(key uint16) error {
	n := w32.SendInput(
		w32.KeyboardInput(w32.KEYBDINPUT{
			Vk: key,
		}),
		w32.KeyboardInput(w32.KEYBDINPUT{
			Vk:    key,
			Flags: w32.KEYEVENTF_KEYUP,
		}),
	)
	if n == 0 {
		return ErrBlocked
	}
	return nil
}

// ForegroundWindowTitle returns the title of the window that currently has the
// focus. The desktop window usually has title "".
func ForegroundWindowTitle() string {
	return w32.GetWindowText(w32.GetForegroundWindow())
}

// ForegroundWindowClassName returns the class name of the window that currently
// has the focus.
func ForegroundWindowClassName() string {
	name, ok := w32.GetClassName(w32.GetForegroundWindow())
	if ok {
		return name
	}
	return ""
}

// ClipboardText returns the contents of the clipboard as text. If the clipboard
// is empty or does not contain text it returns "".
func ClipboardText() string {
	var text string
	if w32.OpenClipboard(0) {
		defer w32.CloseClipboard()
		data := (*uint16)(unsafe.Pointer(w32.GetClipboardData(w32.CF_UNICODETEXT)))
		if data != nil {
			var characters []uint16
			for *data != 0 {
				characters = append(characters, *data)
				data = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(data)) + 2))
			}
			text = syscall.UTF16ToString(characters)
		}
	}
	return text
}

// SetClipboardText sets the contents of the clipboard to the given string.
func SetClipboardText(text string) {
	if w32.OpenClipboard(0) {
		w32.EmptyClipboard()
		data := syscall.StringToUTF16(text)
		clipBuffer := w32.GlobalAlloc(w32.GMEM_DDESHARE, uint32(len(data)*2))
		w32.MoveMemory(
			w32.GlobalLock(clipBuffer),
			unsafe.Pointer(&data[0]),
			uint32(len(data)*2),
		)
		w32.GlobalUnlock(clipBuffer)
		w32.SetClipboardData(
			w32.CF_UNICODETEXT,
			w32.HANDLE(unsafe.Pointer(clipBuffer)),
		)
		w32.CloseClipboard()
	}
}
