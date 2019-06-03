package input

import (
	"errors"
	"strings"
	"unicode"

	"github.com/gonutz/w32"
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

// Type will write the given text by sequentially pressing its keys. Valid
// inputs are characters [a..z] and [A..Z], digits [0..9], spaces ' ', line
// breaks '\n', '\r' or "\r\n" and the backspace character '\b'.
func Type(s string) error {
	// Unify line breaks to '\r' which is the virtual key code for VK_RETURN.
	s = strings.Replace(s, "\r\n", "\r", -1)
	s = strings.Replace(s, "\n", "\r", -1)

	// Predefine shift keys, we might need them more than once.
	shiftDown := w32.KeyboardInput(w32.KEYBDINPUT{
		Vk: w32.VK_LSHIFT,
	})
	shiftUp := w32.KeyboardInput(w32.KEYBDINPUT{
		Vk:    w32.VK_LSHIFT,
		Flags: w32.KEYEVENTF_KEYUP,
	})

	var events []w32.INPUT
	for _, key := range s {
		shift := unicode.IsUpper(key)
		r := unicode.ToUpper(key)
		if !(r == ' ' || r == '\n' || r == '\r' || r == '\b' ||
			'0' <= r && r <= '9' ||
			'A' <= r && r <= 'Z') {
			return errors.New("input.Type: invalid key in string: '" + string(key) + "'")
		}
		down := w32.KeyboardInput(w32.KEYBDINPUT{
			Vk: uint16(r),
		})
		up := w32.KeyboardInput(w32.KEYBDINPUT{
			Vk:    uint16(r),
			Flags: w32.KEYEVENTF_KEYUP,
		})
		if shift {
			events = append(events, shiftDown, down, up, shiftUp)
		} else {
			events = append(events, down, up)
		}
	}
	n := w32.SendInput(events...)
	if n == 0 {
		return ErrBlocked
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
