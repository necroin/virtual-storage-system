package utils

import (
	"unsafe"

	"github.com/lxn/win"
)

func KeyboardDown(keyCode uint16) {
	input := win.KEYBD_INPUT{}
	input.Type = win.INPUT_KEYBOARD
	input.Ki = win.KEYBDINPUT{
		WVk: keyCode,
	}

	win.SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))
}

func KeyboardUp(keyCode uint16) {
	input := win.KEYBD_INPUT{}
	input.Type = win.INPUT_KEYBOARD
	input.Ki = win.KEYBDINPUT{
		WVk:     keyCode,
		DwFlags: win.KEYEVENTF_KEYUP,
	}

	win.SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))
}
