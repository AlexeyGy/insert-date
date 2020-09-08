package main

import (
	"log"
	"syscall"
	"unsafe"
)

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

// PrintCharacter ...
func PrintCharacters(sendInputProc *syscall.Proc, toPrint string) {

	for i := 0; i < len(toPrint); i++ {
		toPass := []input{}

		// ALT key will be sent as a keyup
		var altKeyUp input
		altKeyUp.inputType = 1
		altKeyUp.ki.wVk = Keycodes["ALT"] // ALT
		altKeyUp.ki.dwFlags = 0x0002      // the key is being released.
		toPass = append(toPass, altKeyUp)

		// CTRL key will be sent as a keyup
		var ctrlKeyUp input
		ctrlKeyUp.inputType = 1
		ctrlKeyUp.ki.wVk = Keycodes["CTRL"] // CTRL
		ctrlKeyUp.ki.dwFlags = 0x0002       // the key is being released.
		toPass = append(toPass, ctrlKeyUp)

		// a key
		var key input
		key.inputType = 1 //INPUT_KEYBOARD
		key.ki.wVk = Keycodes[string(toPrint[i])]
		toPass = append(toPass, key)

		ret, _, err := sendInputProc.Call(
			uintptr(len(toPass)),
			uintptr(unsafe.Pointer(&toPass[0])),
			uintptr(unsafe.Sizeof(ctrlKeyUp)),
		)
		if err != nil {
			log.Printf("ret: %v error: %v", ret, err)
		}
	}
}
