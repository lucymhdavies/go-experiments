package main

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/gopherjs/gopherjs/js"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

var (
	// User entered text
	text string
	// Notification text
	notification string

	// an error, to catch the user quitting
	regularTermination = errors.New("regular termination")
)

func update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return regularTermination
	}

	// TODO: delete one char every ~0.1s.
	// e.g. with a waitUntilCanDeleteAgain counter
	// Using ebiten.IsKeyPressed is too fast, and using IsKeyJustPressed will
	// only match the first char
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		tLen := len(text)
		if tLen > 0 {
			text = text[:tLen-1]
		}
	}

	ic := ebiten.InputChars()
	if len(ic) > 0 {
		text = text + string(ic)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Text: %v\n%v", text, notification))

	return nil
}

func jsPrompt() {
	wait := 5
	for wait > 0 {
		notification = fmt.Sprintf("Prompting in %v...", wait)
		time.Sleep(time.Second)
		wait--
	}
	notification = ""

	// Use Javascript's prompt() method
	window := js.Global.Get("window")
	v := window.Call("prompt", "Test", text+" from prompt()")
	if v != nil {
		text = v.String()
	}

	// Wait 5 seconds before trying again
	wait = 5
	for wait > 0 {
		notification = fmt.Sprintf("Focusing textbox in %v...", wait)
		time.Sleep(time.Second)
		wait--
	}
	notification = ""

	/*
		// JS version...

		// Append a textbox
		inputHack = document.createElement("input")
		inputHack.id = "inputHack"
		inputHack.style.cssText = "background: red"
		document.body.appendChild(inputHack)

		// Focus the text box, which should launch vKeyboard
		inputHack.focus
	*/

	// TODO: check if one already exists?

	document := js.Global.Get("document")
	inputHack := document.Call("createElement", "input")
	inputHack.Set("id", "inputHack")
	inputHack.Set("value", text+" from input")
	// Add more CSS to make this actually transparent, but for now make it visible for debugging...
	inputHack.Get("style").Set("cssText", "background: white; color:black")
	document.Get("body").Call("appendChild", inputHack)

	time.Sleep(10 * time.Millisecond)
	inputHack.Call("focus")

	for {
		text = inputHack.Get("value").String()
		time.Sleep(10 * time.Millisecond)

		// Check if still focused?
		if document.Get("activeElement") != inputHack {
			break
		}
	}

	notification = "Text Area Unfocused"
	inputHack.Call("remove")

}

func main() {
	screenWidth := 320
	screenHeight := 240
	scaleFactor := 2.0

	if runtime.GOARCH == "js" {
		// TODO: need to figure out if this is a mobile device
		// i.e. something which only has a virtual keyboard
		text = "JS"
		go jsPrompt()

		scaleFactor = ebiten.DeviceScaleFactor()
		w, h := ebiten.ScreenSizeInFullscreen()
		ebiten.SetFullscreen(true)
		screenWidth = int(float64(w) / scaleFactor)
		screenHeight = int(float64(h) / scaleFactor)

	}

	if err := ebiten.Run(update, screenWidth, screenHeight, scaleFactor, "Keyboard Input Hack"); err != nil {
		if err != regularTermination {
			panic(err)
		}
	}
}
