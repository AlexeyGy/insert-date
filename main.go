package main

import (
	"bytes"
	"fmt"
	"strings"
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

// HOTKEYS to listen to
var HOTKEYS = map[int16]*Hotkey{
	// 1: &Hotkey{1, ModAlt + ModCtrl, 'O'},  // ALT+CTRL+O
	// 2: &Hotkey{2, ModAlt + ModShift, 'M'}, // ALT+SHIFT+M
	// 3: &Hotkey{3, ModAlt + ModCtrl, 'X'},  // ALT+CTRL+X Example of other keys
	4: &Hotkey{4, ModAlt + ModCtrl, 'D'}, // WIN + /
}

//Hotkey ..
type Hotkey struct {
	ID        int // Unique id
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
	return fmt.Sprintf("Hotkey[Id: %d, %s%c]", h.ID, mod, h.KeyCode)
}

// MSG...
type MSG struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}

func main() {

	// from https://stackoverflow.com/questions/38646794/implement-a-global-hotkey-in-golang/
	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

	registerHotkeys(user32)
	getmsg := user32.MustFindProc("GetMessageW")

	for {
		var msg = &MSG{}
		getmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

		// Registered id is in the WPARAM field:
		if id := msg.WPARAM; id != 0 {
			fmt.Println("Hotkey pressed:", HOTKEYS[id])
			PrintCharacters(strings.ToUpper(time.Now().Format("January 02, Mon")))
		}
	}
}

func registerHotkeys(user32 *syscall.DLL) {

	reghotkey := user32.MustFindProc("RegisterHotKey")

	// Register hotkeys:
	for _, v := range HOTKEYS {
		r1, _, err := reghotkey.Call(
			0, uintptr(v.ID), uintptr(v.Modifiers), uintptr(v.KeyCode))
		if r1 == 1 {
			fmt.Println("Registered", v)
		} else {
			fmt.Println("Failed to register", v, ", error:", err)
		}
	}
}
