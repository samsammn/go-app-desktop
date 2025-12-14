package main

import (
	"crypto/rand"
	"fmt"

	"math/big"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func generateOne(charset string, length int) (string, error) {
	if len(charset) == 0 || length <= 0 {
		return "", nil
	}

	var sb strings.Builder
	max := big.NewInt(int64(len(charset)))
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		sb.WriteByte(charset[n.Int64()])
	}

	return sb.String(), nil
}

func buildCharset(lower, upper, digits, symbols bool) string {
	var b strings.Builder
	if lower {
		b.WriteString("abcdefghijklmnopqrstuvwxyz")
	}
	if upper {
		b.WriteString("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
	if digits {
		b.WriteString("0123456789")
	}
	if symbols {
		// common safe symbols
		b.WriteString("!@#$%&*()-_+=[]{}<>?,.:;")
	}

	return b.String()
}

func windowStringGenerator(w fyne.Window) func() {
	return func() {
		notePreviewToolbar.ToolbarObject().Hide()

		lowerCheck := widget.NewCheck("Lowercase (a-z)", func(bool) {})
		lowerCheck.SetChecked(true)
		upperCheck := widget.NewCheck("Uppercase (A-Z)", func(bool) {})
		upperCheck.SetChecked(true)
		digitCheck := widget.NewCheck("Digits (0-9)", func(bool) {})
		digitCheck.SetChecked(true)
		symbolCheck := widget.NewCheck("Symbols (!@#...)", func(bool) {})

		lengthEntry := widget.NewEntry()
		lengthEntry.SetPlaceHolder("Length (e.g. 16)")
		lengthEntry.SetText("16")

		countEntry := widget.NewEntry()
		countEntry.SetPlaceHolder("Count (e.g. 5)")
		countEntry.SetText("5")

		seedLabel := widget.NewLabel("")

		// right: result area
		resultArea := widget.NewMultiLineEntry()
		resultArea.SetPlaceHolder("Generated strings will appear here (one per line).")
		resultArea.Wrapping = fyne.TextWrapWord

		generateBtn := widget.NewButtonWithIcon("Generate", theme.DocumentPrintIcon(), nil)
		clearBtn := widget.NewButtonWithIcon("Clear", theme.ContentClearIcon(), func() {
			resultArea.SetText("")
		})

		generate := func() {
			// parse inputs
			length := 0
			count := 0

			_, err := fmt.Sscan(lengthEntry.Text, &length)
			if err != nil || length <= 0 {
				dialog.ShowError(fmt.Errorf("invalid length: %v", lengthEntry.Text), w)
				return
			}
			_, err = fmt.Sscan(countEntry.Text, &count)
			if err != nil || count <= 0 {
				dialog.ShowError(fmt.Errorf("invalid count: %v", countEntry.Text), w)
				return
			}

			charset := buildCharset(lowerCheck.Checked, upperCheck.Checked, digitCheck.Checked, symbolCheck.Checked)
			if charset == "" {
				dialog.ShowInformation("No charset", "Please choose at least one charset option.", w)
				return
			}

			var out []string
			for i := 0; i < count; i++ {
				s, err := generateOne(charset, length)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				out = append(out, s)
			}
			resultArea.SetText(strings.Join(out, "\n"))
			seedLabel.SetText(fmt.Sprintf("Generated %d strings, %d chars each", count, length))
			dialog.ShowInformation("Generate", "String has been Generated ✔", w)
		}

		generateBtn.OnTapped = generate

		// shortcuts: Ctrl+G for generate
		if c := w.Canvas(); c != nil {
			c.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyG, Modifier: fyne.KeyModifierControl}, func(fyne.Shortcut) {
				generate()
			})
		}

		left := container.NewVBox(
			widget.NewLabelWithStyle("Charset", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			lowerCheck,
			upperCheck,
			digitCheck,
			symbolCheck,
			vPadding(10),
			widget.NewSeparator(),
			widget.NewLabelWithStyle("Options", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			container.NewGridWithColumns(3,
				widget.NewLabel("Length:"), lengthEntry, widget.NewLabel(""),
				widget.NewLabel("Count:"), countEntry, widget.NewLabel(""),
			),
			vPadding(10),
			widget.NewSeparator(),
			vPadding(10),
			container.NewHBox(generateBtn, clearBtn),
			seedLabel,
		)

		resultLabel := widget.NewLabelWithStyle("Result", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		right := container.NewBorder(resultLabel, nil, nil, nil, container.NewScroll(resultArea))

		rightPanel := padding(10, right)
		leftPanel := padding(15, left)

		split := container.NewHSplit(leftPanel, rightPanel)
		split.SetOffset(0.2)

		w.SetContent(
			container.NewBorder(toolbars(w), nil, nil, nil, split),
		)
	}
}
