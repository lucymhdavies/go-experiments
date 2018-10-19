package selector

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/sahilm/fuzzy"
)

var g *gocui.Gui

// Select string from array
// Based heavily on github.com/sahilm/fuzzy/_example
//
// Displays textbox, followed by list.
// User can type a number to select an element from the list.
// User can press up/down to pick from the list.
// User can type to narrow down items from the list.

type selector struct {
	Name string
	// TODO: other config, e.g.
	// fuzzy matching off/on
	// up/down selection
	// mouse enabled
	// numbers off/on
	// multi-select? - probably want this as a different function, with different signature
	// replace spaces with underscore - definitely shouldn't be the default
	// should ctrl-c terminate the running app?
	// show/hide help (keybinding, config) - bottom of screen?
	Items        []string
	SelectedItem string
}

func NewSelector(name string) *selector {
	return &selector{
		Name: name,
	}
}

//
// TODO: almost all the below is copypasta, and is all crap that needs refactoring
//

// Filenames, because I copied from gocui/_example and I'm too lazy to change it
var filenames []string
var filenamesFiltered []string
var filenamesBytes []byte

// err....
var err error

// Store which element is selected
var selectedIndex int

// and its value
var selectedValue string

// if we've passed in a filter, use it
var initialFilter string

// How big is the results pane?
var resultsHeight int

// Select takes a slice of strings and displays an fzf inspired selector,
// allowing the user to pick an item from the list
func (selector selector) SelectFromSlice(list []string) (string, error) {
	return selector.SelectFromSliceWithFilter(list, "")
}

// Select takes a slice of strings and displays an fzf inspired selector,
// allowing the user to pick an item from the list
//
// This version allows you to pass in a predefined filter
func (selector selector) SelectFromSliceWithFilter(list []string, filter string) (string, error) {
	selector.Items = list
	// TODO: temporary. in future, use selector.Items (and make all the below functions Selector methods)
	filenames = selector.Items
	filenamesFiltered = filenames
	filenamesBytes = []byte(strings.Join(filenames, "\n"))
	err = nil
	selectedValue = ""
	selectedIndex = -1
	initialFilter = filter

	g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return "", err
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = false

	g.SetManagerFunc(selector.Layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return "", err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return "", err
	}
	if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, cursorPageUp); err != nil {
		return "", err
	}
	if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, cursorPageDown); err != nil {
		return "", err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return "", err
	}
	// TODO: left/right/home/end

	// TODO: pgup / pgdn
	// FZF behaviour:
	// pgdn - move distance from top to bottom. if this puts cursor off screen, scroll
	// pgup - same as above, in the other direction

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, selector.enter); err != nil {
		return "", err
	}

	//
	// Main Loop. Run, until error
	//

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return "", err
	}

	//return selector.SelectedItem, err
	return selectedValue, err
}
func cursorDown(g *gocui.Gui, v *gocui.View) error {
	results, err := g.View("results")
	if err != nil {
		// handle error
	}
	_, rsY := results.Size()

	if selectedIndex < len(filenamesFiltered)-1 {
		selectedIndex++

		// If we get to the bottom of the screen, scroll

		// TODO: handle scroll down, then back up, then down again
		// i.e. waht I want is "is my selected item at the bottom of the screen?"

		if selectedIndex > rsY-1 {
			scrollView(results, 1)
		}
	} else {
		// Select first element
		selectedIndex = 0
		// Reset scrolling
		results.SetOrigin(0, 0)
	}

	return updateResults()
}
func cursorPageDown(g *gocui.Gui, v *gocui.View) error {
	results, err := g.View("results")
	if err != nil {
		// handle error
	}
	_, rsY := results.Size()

	logs, err := g.View("logs")
	if err != nil {
		// handle error
	}
	logs.Clear()

	if selectedIndex < len(filenamesFiltered)-1 {
		scrollDistance := 10

		// Check if adding that would put us past the last element
		if selectedIndex+scrollDistance > len(filenamesFiltered)-1 {
			scrollDistance = len(filenamesFiltered) - 1 - selectedIndex
		}

		selectedIndex += scrollDistance

		// If we get to the bottom of the screen, scroll

		// TODO: handle scroll down, then back up, then down again
		// i.e. waht I want is "is my selected item at the bottom of the screen?"

		if selectedIndex > rsY-1 {
			scrollView(results, scrollDistance)
		}

		fmt.Fprintf(logs, "resultsSize:%d selected:%d scrollDistance:%d",
			rsY, selectedIndex, scrollDistance)
	} else {
		// Select first element
		selectedIndex = 0
		// Reset scrolling
		results.SetOrigin(0, 0)
	}
	return updateResults()
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	results, err := g.View("results")
	if err != nil {
		// handle error
	}
	_, rsY := results.Size()

	if selectedIndex > 0 {
		selectedIndex--

		// If we get to the top of the screen, scroll
		// TODO: figure out what to do here
		// What we want is "is my selected item at the top of the screen?)
		// For now, always scroll up
		scrollView(results, -1)
	} else {
		// select the last element
		selectedIndex = len(filenamesFiltered) - 1

		// TODO: you are here
		// Scroll to bottom

		scrollView(results, len(filenamesFiltered)-rsY)
	}

	return updateResults()
}
func cursorPageUp(g *gocui.Gui, v *gocui.View) error {
	results, err := g.View("results")
	if err != nil {
		// handle error
	}
	_, rsY := results.Size()
	_, oy := results.Origin()

	logs, err := g.View("logs")
	if err != nil {
		// handle error
	}
	logs.Clear()

	if selectedIndex > 0 {
		scrollDistance := 10
		// Check if removing that would put us past the first element
		if selectedIndex-scrollDistance < 0 {
			// Select first element
			selectedIndex = 0
			// Reset scrolling
			results.SetOrigin(0, 0)
		} else {
			selectedIndex -= scrollDistance

			// If we're very nearly at the top...
			if oy < scrollDistance {
				// Reset scrolling
				results.SetOrigin(0, 0)
			} else {
				scrollView(results, -scrollDistance)
			}
		}
		fmt.Fprintf(logs, "resultsSize:%d selected:%d scrollDistance:%d oy:%d",
			rsY, selectedIndex, scrollDistance, oy)
	} else {
		// select the last element
		selectedIndex = len(filenamesFiltered) - 1

		// Scroll to bottom
		scrollView(results, len(filenamesFiltered)-rsY)
	}

	return updateResults()
}

func (selector selector) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("help", 0, 0, maxX-1, 4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//
		fmt.Fprintf(v, `
- Type digits, or press up/down/pgup/pgdn to select from list, or
- Type letters to filter list
- CTRL-C to cancel
`)
		v.Editable = false
		v.Wrap = true
		v.Frame = true
		v.Title = selector.Name
	}

	if v, err := g.SetView("finder", 0, 5, maxX-1, 7); err != nil {

		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Editable = true
		v.Frame = true
		v.Title = "Type pattern here"
		if _, err := g.SetCurrentView("finder"); err != nil {
			return err
		}
		v.Editor = gocui.EditorFunc(finder)

		// TODO: iterate through filter
		for _, char := range initialFilter {
			v.EditWrite(char)
		}
	}
	if v, err := g.SetView("results", 0, 8, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Wrap = true
		v.Frame = true
		v.Title = "Search Results"
	}

	if v, err := g.SetView("logs", maxX-50, maxY-10, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Wrap = true
		v.Frame = true
		v.Title = "Debug"
	}

	updateResults()
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (selector selector) enter(g *gocui.Gui, v *gocui.View) error {
	// TODO: set this to whichever item is highlighed
	selectedValue = filenamesFiltered[selectedIndex]

	return gocui.ErrQuit
}

func finder(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {

	switch {
	case ch != 0 && mod == 0:
		// Reset index to 0 when typing
		selectedIndex = 0

		// Reset scrolling
		results, err := g.View("results")
		if err != nil {
			// handle error
		}
		results.SetOrigin(0, 0)

		// Add typed character to view
		v.EditWrite(ch)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	case key == gocui.KeySpace:
		// TODO: this should not be the default
		v.EditWrite('_')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	}

	updateResults()
}

func updateResults() error {

	g.Update(func(gui *gocui.Gui) error {
		finder, err := g.View("finder")
		if err != nil {
			// handle error
		}

		results, err := g.View("results")
		if err != nil {
			// handle error
		}
		results.Clear()

		// 		logs, err := g.View("logs")
		// 		if err != nil {
		// 			// handle error
		// 		}
		// 		logs.Clear()

		viewBufferText := strings.TrimSpace(finder.ViewBuffer())

		// Regex returns sequential digits
		re := regexp.MustCompile("[0-9]+")
		// Typed numbers is an array of integers
		typedNumbers := re.FindAllString(viewBufferText, -1)
		// If there are any numbers, we only care about the latest one
		// (kinda arbitrary, but we have to pick
		lastTypedNumber := ""
		if len(typedNumbers) > 0 {
			lastTypedNumber = typedNumbers[len(typedNumbers)-1]
			// TODO: handle the err
			n, _ := strconv.Atoi(lastTypedNumber)

			// If this is within the bounds of the filtered list, select it
			if n < len(filenamesFiltered) {
				selectedIndex = n
			}
		}

		// Now, remove all numbers from the input

		// Regex returns everything but digits
		re = regexp.MustCompile("[^0-9]+")
		typedText := strings.TrimSpace(strings.Join(re.FindAllString(viewBufferText, -1), ""))

		//		rsX, rsY := results.Size()
		// 		fmt.Fprintf(logs, "finder:%s, text:%s, number:%s, lastNum:%s, resultsSize:%d,%d",
		// 			viewBufferText, typedText, typedNumbers, lastTypedNumber, rsX, rsY)

		// TODO: scroll through of there are more than fit on screen

		if typedText == "" {
			// TODO: pre-append numbers
			//fmt.Fprintf(results, "%s", filenamesBytes)

			filenamesFiltered = filenames

			for m, item := range filenames {

				selectChar := ""
				highlightChar := ""
				unHighlightChar := ""

				if m == selectedIndex {
					selectChar = ">"
					highlightChar = "\033[7m"
					unHighlightChar = "\033[0m"
				}

				fmt.Fprintf(results, fmt.Sprintf("%s%1s %3d - %s", highlightChar, selectChar, m, item))
				fmt.Fprintln(results, unHighlightChar)

			}

		} else {
			matches := fuzzy.Find(typedText, filenames)

			filenamesFiltered = []string{}

			for m, match := range matches {

				filenamesFiltered = append(filenamesFiltered, match.Str)

				selectChar := ""
				highlightChar := ""
				unHighlightChar := ""
				matchChar := "\033[31;1m"
				resetMatchChar := "\033[39;0m"

				if m == selectedIndex {
					selectChar = ">"
					highlightChar = "\033[7m"
					unHighlightChar = "\033[0m"
					matchChar = "\033[31;1;7m"
				}

				fmt.Fprintf(results, fmt.Sprintf("%s%1s %3d - ", highlightChar, selectChar, m))
				for i := 0; i < len(match.Str); i++ {
					if contains(i, match.MatchedIndexes) {
						fmt.Fprintf(results, fmt.Sprintf("%s%s%s%s%s", unHighlightChar, matchChar, string(match.Str[i]), resetMatchChar, highlightChar))
					} else {
						fmt.Fprintf(results, string(match.Str[i]))
					}

				}
				fmt.Fprintln(results, unHighlightChar)
			}
		}
		return nil
	})

	return nil

}

func contains(needle int, haystack []int) bool {
	for _, i := range haystack {
		if needle == i {
			return true
		}
	}
	return false
}

func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}
