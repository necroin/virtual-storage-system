package utils

import (
	"unsafe"

	"github.com/lxn/win"
	"github.com/necroin/golibs/utils/winapi"
	"github.com/necroin/golibs/utils/winutils"
	"golang.org/x/sys/windows"
)

func MouseMove(x int32, y int32) {
	cx_screen := win.GetSystemMetrics(win.SM_CXSCREEN)
	cy_screen := win.GetSystemMetrics(win.SM_CYSCREEN)
	real_x := 65535 * x / cx_screen
	real_y := 65535 * y / cy_screen

	mouseInput := win.MOUSE_INPUT{}
	mouseInput.Type = win.INPUT_MOUSE
	mouseInput.Mi = win.MOUSEINPUT{
		Dx: int32(real_x),
		Dy: int32(real_y),
	}
	mouseInput.Mi.MouseData = 0
	mouseInput.Mi.DwExtraInfo = 0
	mouseInput.Mi.Time = 0

	mouseInput.Mi.DwFlags = win.MOUSEEVENTF_ABSOLUTE | win.MOUSEEVENTF_MOVE
	win.SendInput(1, unsafe.Pointer(&mouseInput), int32(unsafe.Sizeof(mouseInput)))
}

func MouseLeftClick(x int32, y int32) {
	mouseInput := win.MOUSE_INPUT{}
	mouseInput.Type = win.INPUT_MOUSE
	mouseInput.Mi = win.MOUSEINPUT{
		Dx:          0,
		Dy:          0,
		MouseData:   0,
		DwExtraInfo: 0,
		Time:        0,
	}

	mouseInput.Mi.DwFlags = win.MOUSEEVENTF_LEFTDOWN
	win.SendInput(1, unsafe.Pointer(&mouseInput), int32(unsafe.Sizeof(mouseInput)))

	mouseInput.Mi.DwFlags = win.MOUSEEVENTF_LEFTUP
	win.SendInput(1, unsafe.Pointer(&mouseInput), int32(unsafe.Sizeof(mouseInput)))
}

func MouseRightClick() {
	mouseInput := win.MOUSE_INPUT{}
	mouseInput.Type = win.INPUT_MOUSE
	mouseInput.Mi = win.MOUSEINPUT{
		Dx:          0,
		Dy:          0,
		MouseData:   0,
		DwExtraInfo: 0,
		Time:        0,
	}

	mouseInput.Mi.DwFlags = win.MOUSEEVENTF_RIGHTDOWN
	win.SendInput(1, unsafe.Pointer(&mouseInput), int32(unsafe.Sizeof(mouseInput)))

	mouseInput.Mi.DwFlags = win.MOUSEEVENTF_RIGHTUP
	win.SendInput(1, unsafe.Pointer(&mouseInput), int32(unsafe.Sizeof(mouseInput)))
}

func FindValidRect(pid winapi.ProcessId) (windows.Rect, bool) {
	windowHandles := winutils.GetWindowHandlesByProcessId(pid)
	clientRects := winutils.GetWindowHandlesClientRects(windowHandles)
	for _, clientRect := range clientRects {
		if winutils.IsValidRect(clientRect) {
			continue
		}
		return clientRect, true
	}
	return windows.Rect{}, false
}
