package main

import (
	"bytes"
	"fmt"
	"log"
	"syscall"
	"time"
	"unsafe"
)

const (
	ModAlt = 1 << iota
	ModCtrl
	ModShift
	ModWin
)

type Hotkey struct {
	Id        int // Unique id
	Modifiers int // Mask of modifiers
	KeyCode   int // Key code, e.g. 'A'
}

// String returns a human-friendly display name of the hotkey
// such as "Hotkey[Id: 1, Alt+Ctrl+O]"
func (h *Hotkey) String() string {
	mod := &bytes.Buffer{}
	if h.Modifiers&ModAlt != 0 {
		mod.WriteString("Alt+")
	}
	if h.Modifiers&ModCtrl != 0 {
		mod.WriteString("Ctrl+")
	}
	if h.Modifiers&ModShift != 0 {
		mod.WriteString("Shift+")
	}
	if h.Modifiers&ModWin != 0 {
		mod.WriteString("Win+")
	}
	return fmt.Sprintf("Hotkey[Id: %d, %s%c]", h.Id, mod, h.KeyCode)
}

type MSG struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}

type keyboardInput struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uint64
}

type input struct {
	inputType uint32
	ki        keyboardInput
	padding   uint64
}

func main() {

	// from https://stackoverflow.com/questions/38646794/implement-a-global-hotkey-in-golang/
	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()
	sendInputProc := user32.MustFindProc("SendInput")
	reghotkey := user32.MustFindProc("RegisterHotKey")

	// Hotkeys to listen to:
	keys := map[int16]*Hotkey{
		1: &Hotkey{1, ModAlt + ModCtrl, 'O'},  // ALT+CTRL+O
		2: &Hotkey{2, ModAlt + ModShift, 'M'}, // ALT+SHIFT+M
		3: &Hotkey{3, ModAlt + ModCtrl, 'X'},  // ALT+CTRL+X
		4: &Hotkey{4, ModAlt + ModCtrl, 'D'},  // WIN + /
	}

	// Register hotkeys:
	for _, v := range keys {
		r1, _, err := reghotkey.Call(
			0, uintptr(v.Id), uintptr(v.Modifiers), uintptr(v.KeyCode))
		if r1 == 1 {
			fmt.Println("Registered", v)
		} else {
			fmt.Println("Failed to register", v, ", error:", err)
		}
	}
	getmsg := user32.MustFindProc("GetMessageW")

	for {
		var msg = &MSG{}
		getmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

		// Registered id is in the WPARAM field:
		if id := msg.WPARAM; id != 0 {
			fmt.Println("Hotkey pressed:", keys[id])
			if id == 3 { // CTRL+ALT+X = Exit
				fmt.Println("CTRL+ALT+X pressed, goodbye...")
				return
			}
			if id == 4 {
				fmt.Println(time.Now().Format(time.RFC3339))

				// send input
				var i input
				i.inputType = 1 //INPUT_KEYBOARD https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-input
				i.ki.wVk = 0x41 // virtual key code for a
				i.ki.time = 0

				// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendinput
				ret, _, err := sendInputProc.Call(
					uintptr(1),
					uintptr(unsafe.Pointer(&i)),
					uintptr(unsafe.Sizeof(i)),
				)
				log.Printf("ret: %v error: %v", ret, err)

				// Release the "..." key https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-keybd_event
				i.ki.dwFlags = 0x0002 // KEYEVENTF_KEYUP for key release
				ret, _, err = sendInputProc.Call(
					uintptr(1),
					uintptr(unsafe.Pointer(&i)),
					uintptr(unsafe.Sizeof(i)))
				log.Printf("ret: %v error: %v", ret, err)
			}
		}

	}
}
