package main

import (
	"fmt"
	"strconv"
	"time"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func isASCII(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func main() {
	// Create the application and window
	myApp := app.New()
	myWindow := myApp.NewWindow("Go Paster Helper")
	myWindow.Resize(fyne.NewSize(500, 400))

	// UI Components

	// Instructions
	instructionLabel := widget.NewLabel("1. Enter ASCII text below.\n2. Set Delay & Interval.\n3. Click Start and focus the target window.")
	instructionLabel.Wrapping = fyne.TextWrapWord

	// Text Input Area
	inputEntry := widget.NewMultiLineEntry()
	inputEntry.SetPlaceHolder("Type text here (ASCII only)...")
	inputEntry.Wrapping = fyne.TextWrapWord

	// Settings
	intervalLabel := widget.NewLabel("Interval (sec):")
	intervalEntry := widget.NewEntry()
	intervalEntry.SetText("0.1")

	delayLabel := widget.NewLabel("Start Delay (sec):")
	delayEntry := widget.NewEntry()
	delayEntry.SetText("5")

	// Status Label
	statusLabel := widget.NewLabelWithStyle("Ready", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Progress Bar
	progressBar := widget.NewProgressBar()
	progressBar.SetValue(0)
	progressBar.Hide()

	// State
	isTyping := false

	// Start Button
	var startBtn *widget.Button
	startBtn = widget.NewButtonWithIcon("Start Typing", theme.MediaPlayIcon(), func() {
		if isTyping {
			return
		}

		// 1. Validate Inputs
		textToType := inputEntry.Text
		if textToType == "" {
			statusLabel.SetText("Error: Text is empty")
			return
		}

		if !isASCII(textToType) {
			statusLabel.SetText("Error: Text contains non-ASCII characters")
			return
		}

		interval, err := strconv.ParseFloat(intervalEntry.Text, 64)
		if err != nil || interval < 0 {
			statusLabel.SetText("Error: Invalid Interval")
			return
		}

		startDelay, err := strconv.ParseFloat(delayEntry.Text, 64)
		if err != nil || startDelay < 0 {
			statusLabel.SetText("Error: Invalid Delay")
			return
		}

		// 2. Set State
		isTyping = true
		startBtn.SetText("Running...")
		progressBar.Show()

		// 3. Start Logic in Goroutine
		go func() {
			defer func() {
				isTyping = false
				statusLabel.SetText("Done!")
				startBtn.SetText("Start Typing")
				progressBar.Hide()
			}()

			// Countdown
			for i := int(startDelay); i > 0; i-- {
				statusLabel.SetText(fmt.Sprintf("Starting in %d seconds...", i))
				progressBar.SetValue(1.0 - (float64(i) / startDelay))
				time.Sleep(1 * time.Second)
			}

			statusLabel.SetText("Typing...")
			progressBar.SetValue(0)

			// Typing Loop
			// Use our custom TypeString function which handles low-level Windows API
			TypeString(textToType, interval)
		}()
	})

	// Layout Construction
	settingsGrid := container.NewGridWithColumns(4,
		delayLabel, delayEntry,
		intervalLabel, intervalEntry,
	)

	bottomContainer := container.NewVBox(
		settingsGrid,
		statusLabel,
		progressBar,
		startBtn,
	)

	content := container.NewBorder(
		instructionLabel,  // Top
		bottomContainer,   // Bottom
		nil,               // Left
		nil,               // Right
		inputEntry,        // Center (expands)
	)

	myWindow.SetContent(content)
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}