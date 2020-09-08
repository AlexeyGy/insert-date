package main

import (
	"bytes"
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/getlantern/systray"
)

const (
	ModAlt = 1 << iota
	ModCtrl
	ModShift
	ModWin
)

// HOTKEYS to listen to
var HOTKEYS = map[int16]*Hotkey{
	1: &Hotkey{4, ModAlt + ModCtrl, 'D'}, // WIN + /
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

// MSG see https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-msg
type MSG struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}

func main() {
	systray.Run(onReady, func() {})
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

func setupSystray() {
	data, err := Asset("images/bitmap.ico")
	if err != nil {
		fmt.Println("Icon reading error", err)
		return
	}

	systray.SetTemplateIcon(data, data)

	systray.SetTitle("Insert Date")
	systray.SetTooltip("Insert the current date (Hotkey CTRL+ALT+D)")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()
}

func onReady() {
	setupSystray()
	run()

}

func run() {
	// from https://stackoverflow.com/questions/38646794/implement-a-global-hotkey-in-golang/
	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

	registerHotkeys(user32)
	getmsg := user32.MustFindProc("GetMessageW")

	sendInputProc := user32.MustFindProc("SendInput")

	for {
		var msg = &MSG{}
		getmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

		// Registered id is in the WPARAM field:
		if id := msg.WPARAM; id != 0 {
			fmt.Println("Hotkey pressed:", HOTKEYS[id])
			PrintCharacters(sendInputProc, strings.ToUpper(time.Now().Format("2006-01-02")))
		}
	}
}
