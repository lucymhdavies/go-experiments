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
		// works on desktop
		inputHack.focus

		// on iOS, we need to do a bit more...
		// (something to do with not being able to automatically focus without
		//  it being prompted by a user input...)
		document.querySelector('canvas').ontouchstart  = function() { document.getElementById("inputHack").focus() }

	*/

	// TODO: check if one already exists?

	document := js.Global.Get("document")
	inputHack := document.Call("createElement", "input")
	inputHack.Set("id", "inputHack")
	inputHack.Set("value", text+" from input")
	// Put it approximately where the an in-game text area would be (vertically)
	// (if this was real, take into account scaleFactor)
	// Also make it really small, and off-screen to the left, so it's unseen
	// Could also add some CSS to make it actually invisible
	inputHack.Get("style").Set("cssText", `
		background: black; color:black;
		border: 0;
		position: absolute;
		top: 100px; left: -20px;
		width: 0px; height: 0px;
`)
	document.Get("body").Call("appendChild", inputHack)

	// Sufficient for desktop, but not for mobile...
	inputHack.Call("focus")
	notification = "Tap anywhere on screen to focus"

	canvas := document.Call("querySelector", "canvas")

	// on iOS (and maybe android?) looks like you need to actually tap
	// somewhere on screen, otherwise JavaScript isn't allowed to focus an input
	canvas.Set("ontouchstart", js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
		inputHack.Call("focus")
		return nil
	},
	))

	for {
		text = inputHack.Get("value").String()
		time.Sleep(10 * time.Millisecond)

		// Check if still focused?
		if document.Get("activeElement") != inputHack {
			notification = "Text Area Unfocused"
		} else {
			notification = "Text Area Focused"
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
